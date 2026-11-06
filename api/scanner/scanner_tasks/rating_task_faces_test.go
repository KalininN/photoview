//go:build !no_face_detection && faces_exif

package scanner_tasks

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/otiai10/copy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/photoview/photoview/api/graphql/models"
	"github.com/photoview/photoview/api/scanner/face_detection"
	"github.com/photoview/photoview/api/scanner/scanner_task"
	"github.com/photoview/photoview/api/test_utils"
)

// Fork-specific: not part of upstream Photoview. Run with `go test -tags faces_exif`.
func TestRatingTaskReadsFaces(t *testing.T) {
	test_utils.FilesystemTest(t)
	db := test_utils.DatabaseTest(t)

	require.NoError(t, face_detection.InitializeFaceDetector(db))
	defer func() { face_detection.GlobalFaceDetector = nil }()

	dir := t.TempDir()
	mediaPath := filepath.Join(dir, "photo.jpg")
	require.NoError(t, copy.Copy(test_utils.PathFromAPIRoot("scanner", "externaltools", "exif", "test_data", "faces.jpg"), mediaPath))

	album := models.Album{Title: "album", Path: dir}
	require.NoError(t, db.Create(&album).Error)

	media := models.Media{Title: "photo.jpg", Path: mediaPath, AlbumID: album.ID, Type: models.MediaTypePhoto}
	require.NoError(t, db.Create(&media).Error)

	ctx := scanner_task.NewTaskContext(context.Background(), db, &album, nil)
	require.NoError(t, ExifTask{}.AfterMediaFound(ctx, &media, true))
	require.NoError(t, RatingTask{}.AfterMediaFound(ctx, &media, true))

	var groups []models.FaceGroup
	require.NoError(t, db.Preload("ImageFaces").Order("label").Find(&groups).Error)
	require.Len(t, groups, 2)

	assert.Equal(t, "2024", *groups[0].Label)
	assert.Equal(t, "Anna Smith", *groups[1].Label)
	assert.Equal(t, "Anna", *groups[1].DisplayLabel)

	require.Len(t, groups[1].ImageFaces, 1)
	assert.InDelta(t, 0.2, groups[1].ImageFaces[0].Rectangle.MinX, 1e-6)
	assert.InDelta(t, 0.6, groups[1].ImageFaces[0].Rectangle.MaxY, 1e-6)
}
