package scanner_tasks

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/otiai10/copy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/photoview/photoview/api/graphql/models"
	"github.com/photoview/photoview/api/scanner/exif"
	"github.com/photoview/photoview/api/scanner/scanner_task"
	"github.com/photoview/photoview/api/test_utils"
)

// Fork-specific: not part of upstream Photoview.

func TestMain(m *testing.M) {
	exif.InitializeEXIFParser()
	test_utils.IntegrationTestRun(m)
}

func TestRatingTask(t *testing.T) {
	test_utils.FilesystemTest(t)
	db := test_utils.DatabaseTest(t)

	dir := t.TempDir()
	mediaPath := filepath.Join(dir, "photo.jpg")
	require.NoError(t, copy.Copy(test_utils.PathFromAPIRoot("scanner", "exif", "test_data", "rated.jpg"), mediaPath))

	album := models.Album{Title: "album", Path: dir}
	require.NoError(t, db.Create(&album).Error)

	media := models.Media{Title: "photo.jpg", Path: mediaPath, AlbumID: album.ID}
	require.NoError(t, db.Create(&media).Error)

	ctx := scanner_task.NewTaskContext(context.Background(), db, &album, nil)
	runTasks := func(newMedia bool) {
		require.NoError(t, ExifTask{}.AfterMediaFound(ctx, &media, newMedia))
		require.NoError(t, RatingTask{}.AfterMediaFound(ctx, &media, newMedia))
	}
	storedRating := func() int64 {
		var stored models.Media
		require.NoError(t, db.First(&stored, media.ID).Error)
		return stored.Rating
	}

	t.Run("new media takes the rating from EXIF", func(t *testing.T) {
		assert.EqualValues(t, 3, storedRating(), "default rating")
		runTasks(true)
		assert.EqualValues(t, 4, storedRating())
		require.NotNil(t, media.ExifID)
	})

	t.Run("unchanged media is not re-read", func(t *testing.T) {
		exifID := *media.ExifID
		runTasks(false)
		assert.Equal(t, exifID, *media.ExifID)
	})

	t.Run("media changed on disk is re-read", func(t *testing.T) {
		oldExifID := *media.ExifID

		require.NoError(t, copy.Copy(test_utils.PathFromAPIRoot("scanner", "exif", "test_data", "rated5.jpg"), mediaPath))
		future := time.Now().Add(time.Minute)
		require.NoError(t, os.Chtimes(mediaPath, future, future))

		runTasks(false)
		assert.EqualValues(t, 5, storedRating())
		assert.NotEqual(t, oldExifID, *media.ExifID)

		var count int64
		require.NoError(t, db.Model(&models.MediaEXIF{}).Where("id = ?", oldExifID).Count(&count).Error)
		assert.Zero(t, count, "outdated EXIF row is deleted")
	})
}
