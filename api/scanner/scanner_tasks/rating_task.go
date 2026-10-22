package scanner_tasks

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/gorm"

	"github.com/photoview/photoview/api/graphql/models"
	"github.com/photoview/photoview/api/scanner/exif"
	"github.com/photoview/photoview/api/scanner/face_detection"
	"github.com/photoview/photoview/api/scanner/scanner_task"
	"github.com/photoview/photoview/api/scanner/scanner_utils"
)

// RatingTask keeps EXIF data, `media.rating` and, with the faces_exif detector,
// faces in sync with the files on disk.
// It must run after ExifTask, which parses EXIF only for newly found media.
//
// Fork-specific: not part of upstream Photoview.
type RatingTask struct {
	scanner_task.ScannerTaskBase
}

func (t RatingTask) AfterMediaFound(ctx scanner_task.TaskContext, media *models.Media, newMedia bool) error {
	tx := ctx.GetDB()

	if !newMedia {
		changed, err := modifiedSince(media.Path, media.UpdatedAt)
		if err != nil {
			log.Printf("WARN: Could not stat media file %s: %s\n", media.Path, err)
			return nil
		}
		if !changed {
			return nil
		}

		if err := refreshEXIF(tx, media); err != nil {
			log.Printf("WARN: Refreshing EXIF for %s failed: %s\n", media.Title, err)
		}
	}

	if err := syncRating(tx, media); err != nil {
		log.Printf("WARN: Updating rating for %s failed: %s\n", media.Title, err)
	}

	if err := syncFaces(tx, media); err != nil {
		scanner_utils.ScannerError(ctx, "Error reading faces from metadata (%s): %s", media.Path, err)
	}

	return nil
}

// syncFaces reads the faces of a new or changed photo from its metadata.
func syncFaces(tx *gorm.DB, media *models.Media) error {
	if !face_detection.FromMetadata || face_detection.GlobalFaceDetector == nil {
		return nil
	}

	if media.Type != models.MediaTypePhoto {
		return nil
	}

	return face_detection.GlobalFaceDetector.DetectFaces(tx, media)
}

func modifiedSince(path string, since time.Time) (bool, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return stat.ModTime().After(since), nil
}

func saveEXIF(tx *gorm.DB, media *models.Media) error {
	exif.ForgetCachedMetadata(media.Path)
	_, err := exif.SaveEXIF(tx, media)
	return err
}

// refreshEXIF re-parses EXIF for a media file that changed on disk, replacing the stored row.
func refreshEXIF(tx *gorm.DB, media *models.Media) error {
	oldExifID := media.ExifID
	media.ExifID = nil

	if err := saveEXIF(tx, media); err != nil {
		media.ExifID = oldExifID
		return err
	}

	// SaveEXIF links the new EXIF row in the database; read the link back
	var stored models.Media
	if err := tx.Select("exif_id").First(&stored, media.ID).Error; err != nil {
		return fmt.Errorf("failed to reload %q: %w", media.Path, err)
	}
	media.ExifID = stored.ExifID

	if oldExifID != nil && media.ExifID != nil && *oldExifID != *media.ExifID {
		if err := tx.Delete(&models.MediaEXIF{}, *oldExifID).Error; err != nil {
			return fmt.Errorf("failed to delete outdated EXIF for %q: %w", media.Path, err)
		}
	}

	// Mark the media as up to date, so it is not refreshed again on the next scan.
	now := time.Now()
	if err := tx.Model(media).UpdateColumn("updated_at", now).Error; err != nil {
		return fmt.Errorf("failed to mark %q as refreshed: %w", media.Path, err)
	}
	media.UpdatedAt = now

	return nil
}

// syncRating copies the rating from the stored EXIF data to the media.
func syncRating(tx *gorm.DB, media *models.Media) error {
	if media.ExifID == nil {
		return nil
	}

	var exif models.MediaEXIF
	if err := tx.First(&exif, *media.ExifID).Error; err != nil {
		return fmt.Errorf("failed to get EXIF for %q: %w", media.Path, err)
	}

	if exif.Rating == nil || *exif.Rating == media.Rating {
		return nil
	}

	if err := tx.Model(media).UpdateColumn("rating", *exif.Rating).Error; err != nil {
		return fmt.Errorf("failed to update rating for %q: %w", media.Path, err)
	}
	media.Rating = *exif.Rating

	return nil
}
