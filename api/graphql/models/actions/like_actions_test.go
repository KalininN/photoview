package actions_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/photoview/photoview/api/graphql/models"
	"github.com/photoview/photoview/api/graphql/models/actions"
	"github.com/photoview/photoview/api/test_utils"
)

func TestLikeMedia(t *testing.T) {
	db := test_utils.DatabaseTest(t)

	password := "1234"
	owner, err := models.RegisterUser(db, "owner", &password, false)
	require.NoError(t, err)
	other, err := models.RegisterUser(db, "other", &password, false)
	require.NoError(t, err)

	album := models.Album{Title: "root", Path: "/photos"}
	require.NoError(t, db.Save(&album).Error)
	require.NoError(t, db.Model(&owner).Association("Albums").Append(&album))

	updatedAt := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	media := models.Media{Title: "pic", Path: "/photos/pic", AlbumID: album.ID}
	require.NoError(t, db.Save(&media).Error)
	require.NoError(t, db.Model(&media).UpdateColumn("updated_at", updatedAt).Error)

	t.Run("owner likes twice", func(t *testing.T) {
		_, err := actions.LikeMedia(db, owner, media.ID)
		require.NoError(t, err)
		liked, err := actions.LikeMedia(db, owner, media.ID)
		require.NoError(t, err)

		assert.EqualValues(t, 2, liked.Likes)
		// The scanner compares updated_at with the file's modification time
		assert.True(t, liked.UpdatedAt.Equal(updatedAt), "updated_at changed to %s", liked.UpdatedAt)
	})

	t.Run("other user is forbidden", func(t *testing.T) {
		_, err := actions.LikeMedia(db, other, media.ID)
		assert.Error(t, err)

		var stored models.Media
		require.NoError(t, db.First(&stored, media.ID).Error)
		assert.EqualValues(t, 2, stored.Likes)
	})

	t.Run("unknown media", func(t *testing.T) {
		_, err := actions.LikeMedia(db, owner, media.ID+100)
		assert.Error(t, err)
	})
}
