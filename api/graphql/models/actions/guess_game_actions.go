package actions

// Fork-specific: not part of upstream Photoview.

import (
	"errors"
	"fmt"
	"math"
	mathrand "math/rand/v2"
	"strings"
	"time"

	"github.com/photoview/photoview/api/graphql/models"
	"github.com/photoview/photoview/api/guessgame"
	"github.com/photoview/photoview/api/utils"
	"gorm.io/gorm"
)

const (
	// GuessDefaultMinRating is the minimum rating of the photos a game picks from, unless the filter sets one.
	GuessDefaultMinRating = 4
	// GuessMaxRounds is the most rounds a game can have.
	GuessMaxRounds = 10
	// GuessMaxOpenGames is the most unfinished games an account can have at once.
	GuessMaxOpenGames = 50
	// GuessIdleTimeout is how long an unfinished game may go without a guess before it is abandoned.
	GuessIdleTimeout = 7 * 24 * time.Hour
	// GuessRepeatWindow is how long a photo played on a device is avoided in that device's next games.
	GuessRepeatWindow = 30 * 24 * time.Hour

	guessMaxDeviceKeyLength = 64
)

var (
	ErrGuesserDisabled   = errors.New("4D-guesser is disabled")
	ErrGuessGameNotFound = errors.New("game not found")
	ErrGuessNotFound     = errors.New("round not found")
	ErrGuessGameClosed   = errors.New("this game was replaced by a newer game or expired")
	ErrGuessOutOfOrder   = errors.New("earlier rounds of this game have not been guessed yet")
	ErrGuessTooManyGames = errors.New("too many unfinished games on this account, finish or replace one first")
)

// guessPoolItem is an eligible photo with its answer.
type guessPoolItem struct {
	MediaID   int
	Latitude  float64
	Longitude float64
	DateShot  time.Time
}

func (item guessPoolItem) point() guessgame.Point {
	// DateShot holds the wall-clock time of the photo, stored as UTC.
	date := item.DateShot.UTC()
	return guessgame.Point{
		Latitude:  item.Latitude,
		Longitude: item.Longitude,
		Year:      date.Year(),
		Month:     int(date.Month()),
	}
}

func guessMinRating(filter *models.GuessGameFilter) int {
	if filter != nil && filter.MinRating != nil {
		return models.ClampMinRating(*filter.MinRating)
	}
	return GuessDefaultMinRating
}

// guessPool loads the photos of the account that a game can use:
// photos with GPS coordinates, an EXIF date and a displayable image, rated at least minRating.
func guessPool(db *gorm.DB, user *models.User, minRating int, filter *models.GuessGameFilter) ([]guessPoolItem, error) {
	query := db.Table("media").
		Select("media.id AS media_id, media_exif.gps_latitude AS latitude, "+
			"media_exif.gps_longitude AS longitude, media_exif.date_shot AS date_shot").
		Joins("JOIN media_exif ON media.exif_id = media_exif.id").
		Where("media.album_id IN (?)", db.Table("user_albums").Select("user_albums.album_id").Where("user_albums.user_id = ?", user.ID)).
		Where("media.type = ?", models.MediaTypePhoto).
		Where("media_exif.gps_latitude IS NOT NULL AND media_exif.gps_longitude IS NOT NULL").
		Where("NOT (media_exif.gps_latitude = 0 AND media_exif.gps_longitude = 0)").
		Where("media_exif.date_shot IS NOT NULL").
		Where("media_exif.date_shot >= ?", time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)).
		Where("media_exif.date_shot <= ?", time.Now().UTC().Add(24*time.Hour)).
		Where("EXISTS (SELECT 1 FROM media_urls WHERE media_urls.media_id = media.id AND media_urls.purpose IN (?))",
			[]models.MediaPurpose{models.PhotoHighRes, models.PhotoThumbnail})

	query = models.FilterMinRating(query, &minRating)

	if filter != nil && len(filter.AlbumIds) > 0 {
		albums, err := models.GetChildrenFromAlbums(db, nil, filter.AlbumIds)
		if err != nil {
			return nil, fmt.Errorf("get sub-albums: %w", err)
		}

		albumIDs := make([]int, 0, len(albums))
		for _, album := range albums {
			albumIDs = append(albumIDs, album.ID)
		}
		query = query.Where("media.album_id IN (?)", albumIDs)
	}

	var pool []guessPoolItem
	if err := query.Order("media.id").Scan(&pool).Error; err != nil {
		return nil, fmt.Errorf("load guess game photos: %w", err)
	}

	return pool, nil
}

// GuessGameInfo describes the photos a game with this filter would pick from.
func GuessGameInfo(db *gorm.DB, user *models.User, filter *models.GuessGameFilter) (*models.GuessGameInfo, error) {
	if !utils.GuesserEnabled() {
		return nil, ErrGuesserDisabled
	}

	pool, err := guessPool(db, user, guessMinRating(filter), filter)
	if err != nil {
		return nil, err
	}

	info := &models.GuessGameInfo{EligibleCount: len(pool)}
	if len(pool) == 0 {
		return info, nil
	}

	minYear, maxYear := math.MaxInt, math.MinInt
	south, west, north, east := 90.0, 180.0, -90.0, -180.0
	for _, item := range pool {
		point := item.point()
		minYear = min(minYear, point.Year)
		maxYear = max(maxYear, point.Year)
		south = min(south, point.Latitude)
		north = max(north, point.Latitude)
		west = min(west, point.Longitude)
		east = max(east, point.Longitude)
	}

	info.MinYear = &minYear
	info.MaxYear = &maxYear
	info.SouthWest = &models.Coordinates{Latitude: south, Longitude: west}
	info.NorthEast = &models.Coordinates{Latitude: north, Longitude: east}
	return info, nil
}

// StartGuessGameArgs are the settings of a new game.
type StartGuessGameArgs struct {
	DeviceKey     string
	Rounds        int
	Difficulty    models.GuessDifficulty
	Filter        *models.GuessGameFilter
	ReplaceGameID *string
}

// StartGuessGame picks the photos of a new game and stores it with its answers.
func StartGuessGame(db *gorm.DB, user *models.User, args StartGuessGameArgs) (*models.GuessGame, error) {
	if !utils.GuesserEnabled() {
		return nil, ErrGuesserDisabled
	}

	deviceKey := strings.TrimSpace(args.DeviceKey)
	if deviceKey == "" || len(deviceKey) > guessMaxDeviceKeyLength {
		return nil, fmt.Errorf("deviceKey must be 1-%d characters", guessMaxDeviceKeyLength)
	}
	if args.Rounds < 1 || args.Rounds > GuessMaxRounds {
		return nil, fmt.Errorf("rounds must be 1-%d", GuessMaxRounds)
	}
	difficulty := args.Difficulty
	if !difficulty.IsValid() {
		difficulty = models.GuessDifficultyAuto
	}

	minRating := guessMinRating(args.Filter)
	pool, err := guessPool(db, user, minRating, args.Filter)
	if err != nil {
		return nil, err
	}
	if len(pool) < args.Rounds {
		return nil, fmt.Errorf("only %d photos with GPS coordinates and an EXIF date match, a game of %d rounds needs more",
			len(pool), args.Rounds)
	}

	var game *models.GuessGame
	err = db.Transaction(func(tx *gorm.DB) error {
		now := time.Now()

		if args.ReplaceGameID != nil && *args.ReplaceGameID != "" {
			if err := tx.Model(&models.GuessGame{}).
				Where("public_id = ? AND user_id = ? AND finished_at IS NULL", *args.ReplaceGameID, user.ID).
				Update("abandoned", true).Error; err != nil {
				return fmt.Errorf("abandon replaced game: %w", err)
			}
		}

		if err := guessCheckOpenGames(tx, user, now); err != nil {
			return err
		}

		lastPlayed, err := guessLastPlayed(tx, user, deviceKey)
		if err != nil {
			return err
		}

		poolIDs := make([]int, len(pool))
		poolByID := make(map[int]guessPoolItem, len(pool))
		points := make([]guessgame.Point, len(pool))
		for i, item := range pool {
			poolIDs[i] = item.MediaID
			poolByID[item.MediaID] = item
			points[i] = item.point()
		}

		rng := mathrand.New(mathrand.NewPCG(mathrand.Uint64(), mathrand.Uint64()))
		picked := guessgame.PickMedia(poolIDs, lastPlayed, now.Add(-GuessRepeatWindow), args.Rounds, rng)
		scales := guessgame.ScalesFor(guessgame.Difficulty(difficulty), points)

		game = &models.GuessGame{
			PublicID:        guessgame.NewPublicID(),
			UserID:          user.ID,
			DeviceKey:       deviceKey,
			RoundCount:      len(picked),
			Difficulty:      string(difficulty),
			DistanceScaleKm: scales.DistanceKm,
			TimeScaleMonths: scales.TimeMonths,
			MinRating:       minRating,
		}
		if err := tx.Create(game).Error; err != nil {
			return fmt.Errorf("create guess game: %w", err)
		}

		rounds := make([]*models.GuessRound, len(picked))
		for i, mediaID := range picked {
			answer := poolByID[mediaID].point()
			rounds[i] = &models.GuessRound{
				PublicID:        guessgame.NewPublicID(),
				GameID:          game.ID,
				RoundIndex:      i,
				MediaID:         &mediaID,
				AnswerLatitude:  answer.Latitude,
				AnswerLongitude: answer.Longitude,
				AnswerYear:      answer.Year,
				AnswerMonth:     answer.Month,
			}
		}
		if err := tx.Create(&rounds).Error; err != nil {
			return fmt.Errorf("create guess rounds: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return game, nil
}

// guessCheckOpenGames makes room for a new game: it abandons idle games when the account has
// too many unfinished ones, and fails if that is not enough.
func guessCheckOpenGames(tx *gorm.DB, user *models.User, now time.Time) error {
	openGames := func() (int64, error) {
		var count int64
		err := tx.Model(&models.GuessGame{}).
			Where("user_id = ? AND finished_at IS NULL AND abandoned = ?", user.ID, false).
			Count(&count).Error
		return count, err
	}

	count, err := openGames()
	if err != nil {
		return fmt.Errorf("count open guess games: %w", err)
	}
	if count < GuessMaxOpenGames {
		return nil
	}

	if err := tx.Model(&models.GuessGame{}).
		Where("user_id = ? AND finished_at IS NULL AND abandoned = ? AND updated_at < ?", user.ID, false, now.Add(-GuessIdleTimeout)).
		Update("abandoned", true).Error; err != nil {
		return fmt.Errorf("abandon idle guess games: %w", err)
	}

	if count, err = openGames(); err != nil {
		return fmt.Errorf("count open guess games: %w", err)
	}
	if count >= GuessMaxOpenGames {
		return ErrGuessTooManyGames
	}
	return nil
}

// guessLastPlayed returns when each photo was last played on the device.
func guessLastPlayed(tx *gorm.DB, user *models.User, deviceKey string) (map[int]time.Time, error) {
	var rows []struct {
		MediaID  int
		PlayedAt time.Time
	}

	err := tx.Table("guess_rounds").
		Select("guess_rounds.media_id AS media_id, guess_rounds.created_at AS played_at").
		Joins("JOIN guess_games ON guess_rounds.game_id = guess_games.id").
		Where("guess_games.user_id = ? AND guess_games.device_key = ?", user.ID, deviceKey).
		Where("guess_rounds.media_id IS NOT NULL").
		Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("load played guess photos: %w", err)
	}

	lastPlayed := make(map[int]time.Time, len(rows))
	for _, row := range rows {
		if row.PlayedAt.After(lastPlayed[row.MediaID]) {
			lastPlayed[row.MediaID] = row.PlayedAt
		}
	}
	return lastPlayed, nil
}

// GetGuessGame returns a game of the account by its public id.
func GetGuessGame(db *gorm.DB, user *models.User, publicID string) (*models.GuessGame, error) {
	if !utils.GuesserEnabled() {
		return nil, ErrGuesserDisabled
	}

	var game models.GuessGame
	err := db.Where("public_id = ? AND user_id = ?", publicID, user.ID).First(&game).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrGuessGameNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get guess game: %w", err)
	}

	return &game, nil
}

// GuessGameRounds returns the rounds of a game in order.
func GuessGameRounds(db *gorm.DB, game *models.GuessGame) ([]*models.GuessRound, error) {
	var rounds []*models.GuessRound
	if err := db.Where("game_id = ?", game.ID).Order("round_index").Find(&rounds).Error; err != nil {
		return nil, fmt.Errorf("get guess rounds: %w", err)
	}
	return rounds, nil
}

// SubmitGuess stores and scores the guess of a round. A round is guessed once: guessing it
// again, for example from a second browser tab, returns the stored result unchanged.
func SubmitGuess(db *gorm.DB, user *models.User, roundPublicID string, guess models.GuessInput) (*models.GuessRoundResult, error) {
	if !utils.GuesserEnabled() {
		return nil, ErrGuesserDisabled
	}

	if math.IsNaN(guess.Latitude) || guess.Latitude < -90 || guess.Latitude > 90 {
		return nil, errors.New("latitude must be between -90 and 90")
	}
	if math.IsNaN(guess.Longitude) || guess.Longitude < -180 || guess.Longitude > 180 {
		return nil, errors.New("longitude must be between -180 and 180")
	}
	if guess.Month < 1 || guess.Month > 12 {
		return nil, errors.New("month must be between 1 and 12")
	}
	if guess.Year < 1 || guess.Year > 9999 {
		return nil, errors.New("year must be between 1 and 9999")
	}

	var round models.GuessRound
	err := db.Transaction(func(tx *gorm.DB) error {
		err := tx.Select("guess_rounds.*").
			Joins("JOIN guess_games ON guess_games.id = guess_rounds.game_id").
			Where("guess_rounds.public_id = ? AND guess_games.user_id = ?", roundPublicID, user.ID).
			Preload("Game").
			First(&round).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrGuessNotFound
		}
		if err != nil {
			return fmt.Errorf("get guess round: %w", err)
		}

		if round.Guessed() {
			return nil
		}
		if round.Game.Abandoned || round.Game.FinishedAt != nil {
			return ErrGuessGameClosed
		}

		var unguessedBefore int64
		if err := tx.Model(&models.GuessRound{}).
			Where("game_id = ? AND round_index < ? AND guessed_at IS NULL", round.GameID, round.RoundIndex).
			Count(&unguessedBefore).Error; err != nil {
			return fmt.Errorf("check earlier guess rounds: %w", err)
		}
		if unguessedBefore > 0 {
			return ErrGuessOutOfOrder
		}

		scales := round.Game.Scales()
		distance := guessgame.DistanceKm(guess.Latitude, guess.Longitude, round.AnswerLatitude, round.AnswerLongitude)
		monthsOff := guessgame.MonthsOff(guess.Year, guess.Month, round.AnswerYear, round.AnswerMonth)
		now := time.Now()

		update := tx.Model(&models.GuessRound{}).
			Where("id = ? AND guessed_at IS NULL", round.ID).
			Updates(map[string]any{
				"guess_latitude":  guess.Latitude,
				"guess_longitude": guess.Longitude,
				"guess_year":      guess.Year,
				"guess_month":     guess.Month,
				"distance_km":     distance,
				"months_off":      monthsOff,
				"location_score":  guessgame.LocationScore(distance, scales),
				"time_score":      guessgame.TimeScore(monthsOff, scales),
				"guessed_at":      now,
			})
		if err := update.Error; err != nil {
			return fmt.Errorf("store guess: %w", err)
		}
		// When another request guessed the round first, its result is the one that counts.
		if update.RowsAffected > 0 {
			if err := guessUpdateGameTotal(tx, round.GameID, now); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		// A concurrent guess on the same round may have won the race (SQLite refuses the
		// second write as "database is locked"): then the stored result is the answer.
		if round.ID != 0 && db.First(&round, round.ID).Error == nil && round.Guessed() {
			if result := round.Result(); result != nil {
				return result, nil
			}
		}
		return nil, err
	}

	// Reload after the commit: inside the transaction, MySQL may still show the round
	// as it was before a concurrent guess that won the race.
	if err := db.First(&round, round.ID).Error; err != nil {
		return nil, fmt.Errorf("reload guess round: %w", err)
	}

	result := round.Result()
	if result == nil {
		return nil, errors.New("guess round has no result")
	}
	return result, nil
}

// guessUpdateGameTotal recomputes the total of a game from its rounds, and finishes it when all are guessed.
func guessUpdateGameTotal(tx *gorm.DB, gameID int, now time.Time) error {
	var totals struct {
		Total   int
		Guessed int
		Rounds  int
	}
	if err := tx.Model(&models.GuessRound{}).
		Select("COALESCE(SUM(location_score + time_score), 0) AS total, COUNT(guessed_at) AS guessed, COUNT(*) AS rounds").
		Where("game_id = ?", gameID).
		Scan(&totals).Error; err != nil {
		return fmt.Errorf("sum guess scores: %w", err)
	}

	updates := map[string]any{
		"total_score": totals.Total,
		"updated_at":  now,
	}
	if totals.Guessed == totals.Rounds {
		updates["finished_at"] = now
	}

	if err := tx.Model(&models.GuessGame{}).Where("id = ?", gameID).Updates(updates).Error; err != nil {
		return fmt.Errorf("update guess game total: %w", err)
	}
	return nil
}
