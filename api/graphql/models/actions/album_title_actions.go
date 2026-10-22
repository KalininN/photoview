package actions

import (
	"path/filepath"
	"strings"
	"unicode"

	"github.com/photoview/photoview/api/graphql/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// Custom album titles.
//
// Fork-specific: not part of upstream Photoview.

// SanitizeTitle keeps letters, numbers, spaces and common punctuation of a title.
func SanitizeTitle(title string) string {
	allowedSpecials := "-_:.,;()[]{}!@#%&+=*?| "
	notAllowed := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c) && !strings.ContainsRune(allowedSpecials, c)
	}
	allowed := strings.FieldsFunc(title, notAllowed)
	return strings.Join(allowed, "")
}

func SetAlbumTitle(db *gorm.DB, user *models.User, albumID int, title string) (*models.Album, error) {
	var album models.Album
	if err := db.Find(&album, albumID).Error; err != nil {
		return nil, err
	}

	ownsAlbum, err := user.OwnsAlbum(db, &album)
	if err != nil {
		return nil, err
	}

	if !ownsAlbum {
		return nil, errors.New("forbidden")
	}

	titleSanitized := SanitizeTitle(title)
	if titleSanitized == "" {
		return nil, errors.New("invalid title")
	}
	if err := db.Model(&album).Update("title", titleSanitized).Error; err != nil {
		return nil, err
	}

	return &album, nil
}

func ResetAlbumTitle(db *gorm.DB, user *models.User, albumID int) (*models.Album, error) {
	var album models.Album
	if err := db.Find(&album, albumID).Error; err != nil {
		return nil, err
	}

	ownsAlbum, err := user.OwnsAlbum(db, &album)
	if err != nil {
		return nil, err
	}

	if !ownsAlbum {
		return nil, errors.New("forbidden")
	}

	if err := db.Model(&album).Update("title", filepath.Base(album.Path)).Error; err != nil {
		return nil, err
	}

	return &album, nil
}
