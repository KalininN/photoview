package exif_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/photoview/photoview/api/scanner/exif"
	"github.com/photoview/photoview/api/test_utils"
)

// Fork-specific: not part of upstream Photoview.
func TestExiftoolParserRating(t *testing.T) {
	test_utils.FilesystemTest(t)

	parser, err := exif.NewExiftoolParser()
	if err != nil {
		t.Skip("exiftool is not available")
	}

	rated, err := parser.ParseExif("./test_data/rated.jpg")
	require.NoError(t, err)
	require.NotNil(t, rated)
	require.NotNil(t, rated.Rating)
	assert.EqualValues(t, 4, *rated.Rating)

	unrated, err := parser.ParseExif("./test_data/bird.jpg")
	require.NoError(t, err)
	require.NotNil(t, unrated)
	assert.Nil(t, unrated.Rating)
}
