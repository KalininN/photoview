package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/photoview/photoview/api/database/drivers"
	"github.com/photoview/photoview/api/database/migrations"
	"github.com/photoview/photoview/api/graphql/models"
	"github.com/photoview/photoview/api/utils"

	"github.com/go-sql-driver/mysql"
	gorm_mysql "gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func GetMysqlAddress(addressString string) (string, error) {
	if addressString == "" {
		return "", fmt.Errorf("Environment variable %s missing, exiting", utils.EnvMysqlURL.GetName())
	}

	config, err := mysql.ParseDSN(addressString)
	if err != nil {
		return "", fmt.Errorf("could not parse mysql url: %w", err)
	}

	config.MultiStatements = true
	config.ParseTime = true

	return config.FormatDSN(), nil
}

func GetPostgresAddress(addressString string) (*url.URL, error) {
	if addressString == "" {
		return nil, fmt.Errorf("Environment variable %s missing, exiting", utils.EnvPostgresURL.GetName())
	}

	address, err := url.Parse(addressString)
	if err != nil {
		return nil, fmt.Errorf("could not parse postgres url: %w", err)
	}

	return address, nil
}

func GetSqliteAddress(path string) (*url.URL, error) {
	if path == "" {
		path = "photoview.db"
	}

	address, err := url.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("could not parse sqlite url (%s): %w", path, err)
	}

	queryValues := address.Query()
	queryValues.Add("cache", "shared")
	queryValues.Add("mode", "rwc")
	// queryValues.Add("_busy_timeout", "60000") // 1 minute
	queryValues.Add("_journal_mode", "WAL")    // Write-Ahead Logging (WAL) mode
	queryValues.Add("_locking_mode", "NORMAL") // allows concurrent reads and writes
	queryValues.Add("_foreign_keys", "ON")     // Enforc foreign key constraints.
	address.RawQuery = queryValues.Encode()

	// log.Panicf("%s", address.String())

	return address, nil
}

func ConfigureDatabase(config *gorm.Config) (*gorm.DB, error) {
	var databaseDialect gorm.Dialector
	driver := drivers.DatabaseDriverFromEnv()
	log.Printf("Utilizing %s database driver based on environment variables", driver)

	switch driver {
	case drivers.MYSQL:
		mysqlAddress, err := GetMysqlAddress(utils.EnvMysqlURL.GetValue())
		if err != nil {
			return nil, err
		}
		databaseDialect = gorm_mysql.Open(mysqlAddress)
	case drivers.SQLITE:
		sqliteAddress, err := GetSqliteAddress(utils.EnvSqlitePath.GetValue())
		if err != nil {
			return nil, err
		}
		databaseDialect = sqlite.Open(sqliteAddress.String())

	case drivers.POSTGRES:
		postgresAddress, err := GetPostgresAddress(utils.EnvPostgresURL.GetValue())
		if err != nil {
			return nil, err
		}
		databaseDialect = postgres.Open(postgresAddress.String())
	}

	db, err := gorm.Open(databaseDialect, config)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// SetupDatabase connects to the database using environment variables
func SetupDatabase() (*gorm.DB, error) {

	config := gorm.Config{}

	// Configure database logging
	if utils.DevelopmentMode() {
		config.Logger = logger.Default.LogMode(logger.Info)
	} else {
		config.Logger = logger.Default.LogMode(logger.Warn)
	}

	var db *gorm.DB

	for retryCount := 1; retryCount <= 5; retryCount++ {

		var err error
		db, err = ConfigureDatabase(&config)
		if err == nil {
			sqlDB, dbErr := db.DB()
			if dbErr != nil {
				return nil, dbErr
			}

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			err = sqlDB.PingContext(ctx)
			cancel()

			sqlDB.SetMaxOpenConns(80)

			if err == nil {
				return db, nil
			}
		}

		log.Printf("WARN: Could not ping database: %s. Will retry after 5 seconds\n", err)
		time.Sleep(time.Duration(5) * time.Second)
	}

	return db, nil
}

var database_models []interface{} = []interface{}{
	&models.User{},
	&models.AccessToken{},
	&models.SiteInfo{},
	&models.Media{},
	&models.MediaURL{},
	&models.Album{},
	&models.MediaEXIF{},
	&models.VideoMetadata{},
	&models.ShareToken{},
	&models.UserMediaData{},
	&models.UserAlbums{},
	&models.UserPreferences{},

	// Face detection
	&models.FaceGroup{},
	&models.ImageFace{},
}

func MigrateDatabase(db *gorm.DB) error {

	if err := db.SetupJoinTable(&models.User{}, "Albums", &models.UserAlbums{}); err != nil {
		log.Printf("Setup UserAlbums join table failed: %v\n", err)
	}

	if err := db.AutoMigrate(database_models...); err != nil {
		log.Printf("Auto migration failed: %v\n", err)
	}

	// v2.1.0 - Replaced by Media.CreatedAt
	if db.Migrator().HasColumn(&models.Media{}, "date_imported") {
		db.Migrator().DropColumn(&models.Media{}, "date_imported")
	}

	// v2.3.0 - Changed type of MediaEXIF.Exposure and MediaEXIF.Flash
	// from string values to decimal and int respectively
	if err := migrateExifFields(db); err != nil {
		log.Printf("Failed to run exif fields migration: %v\n", err)
	}

	// Remove invalid GPS data from DB
	if err := migrations.MigrateForExifGPSCorrection(db); err != nil {
		log.Printf("Failed to run exif GPS correction migration: %v\n", err)
	}

	// v2.5.0 - Remove Thumbnail Method for Downsampliing filters
	if db.Migrator().HasColumn(&models.SiteInfo{}, "thumbnail_method") {
		db.Migrator().DropColumn(&models.SiteInfo{}, "thumbnail_method")
	}

	// Add rating column to media and media_exif tables
	if !db.Migrator().HasColumn(&models.Media{}, "rating") {
		if err := db.Migrator().AddColumn(&models.Media{}, "rating"); err != nil {
			log.Printf("Failed to add rating column to media table: %v\n", err)
		} else {
			// Set default rating to 3 for all existing media
			if err := db.Model(&models.Media{}).Where("rating IS NULL").Update("rating", 3).Error; err != nil {
				log.Printf("Failed to set default rating for existing media: %v\n", err)
			}
		}
	}

	if !db.Migrator().HasColumn(&models.MediaEXIF{}, "rating") {
		if err := db.Migrator().AddColumn(&models.MediaEXIF{}, "rating"); err != nil {
			log.Printf("Failed to add rating column to media_exif table: %v\n", err)
		} else {
			// Set default rating to 3 for all existing media exif
			if err := db.Model(&models.MediaEXIF{}).Where("rating IS NULL").Update("rating", 3).Error; err != nil {
				log.Printf("Failed to set default rating for existing media exif: %v\n", err)
			}
		}
	}

	// Add display_label column to face_groups table
	if !db.Migrator().HasColumn(&models.FaceGroup{}, "display_label") {
		if err := db.Migrator().AddColumn(&models.FaceGroup{}, "display_label"); err != nil {
			log.Printf("Failed to add display_label column to face_groups table: %v\n", err)
		} else {
			// Set default display_label to label where display_label is null
			if err := db.Model(&models.FaceGroup{}).Where("display_label IS NULL").UpdateColumn("display_label", gorm.Expr("label")).Error; err != nil {
				log.Printf("Failed to set default display_label for existing face groups: %v\n", err)
			}
		}
	}

	// Add cover_id column to face_groups table
	if !db.Migrator().HasColumn(&models.FaceGroup{}, "cover_id") {
		if err := db.Migrator().AddColumn(&models.FaceGroup{}, "cover_id"); err != nil {
			log.Printf("Failed to add cover_id column to face_groups table: %v\n", err)
		}
	}

	return nil
}

func ClearDatabase(db *gorm.DB) error {
	var errs []error
	for _, model := range database_models {
		if err := db.Migrator().DropTable(model); err != nil {
			errs = append(errs, err)
		}
	}

	if err := errors.Join(errs...); err != nil {
		return fmt.Errorf("drop tables error: %w", err)
	}

	return nil
}
