package actions

import (
	"github.com/photoview/photoview/api/graphql/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// LikeMedia adds one like to a media the user can access. Likes from all users
// add up, and a user can like a media several times.
//
// Fork-specific: not part of upstream Photoview.
func LikeMedia(db *gorm.DB, user *models.User, mediaID int) (*models.Media, error) {
	var media models.Media
	if err := db.First(&media, mediaID).Error; err != nil {
		return nil, errors.Wrapf(err, "could not find media with id %d", mediaID)
	}

	ownsAlbum, err := user.OwnsAlbum(db, &models.Album{Model: models.Model{ID: media.AlbumID}})
	if err != nil {
		return nil, err
	}
	if !ownsAlbum {
		return nil, errors.New("forbidden")
	}

	// UpdateColumn adds atomically and leaves updated_at alone, which the scanner
	// compares with the file's modification time to find changed files.
	if err := db.Model(&media).UpdateColumn("likes", gorm.Expr("likes + 1")).Error; err != nil {
		return nil, errors.Wrapf(err, "could not update likes for media with id %d", mediaID)
	}

	if err := db.First(&media, mediaID).Error; err != nil {
		return nil, err
	}

	return &media, nil
}
