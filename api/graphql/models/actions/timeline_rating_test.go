package actions_test

import (
	"testing"
	"time"

	"github.com/photoview/photoview/api/graphql/models"
	"github.com/photoview/photoview/api/graphql/models/actions"
	"github.com/photoview/photoview/api/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Fork-specific: not part of upstream Photoview.
func TestMyTimelineMinRating(t *testing.T) {
	db := test_utils.DatabaseTest(t)

	password := "1234"
	user, err := models.RegisterUser(db, "user", &password, false)
	require.NoError(t, err)

	album := models.Album{Title: "root", Path: "/photos"}
	require.NoError(t, db.Save(&album).Error)
	require.NoError(t, db.Model(&user).Association("Albums").Append(&album))

	media := []models.Media{
		{Title: "unrated", Path: "/photos/unrated", AlbumID: album.ID, DateShot: time.Unix(1632758400, 0)},
		{Title: "two", Path: "/photos/two", AlbumID: album.ID, DateShot: time.Unix(1632758450, 0), Rating: 2},
		{Title: "four", Path: "/photos/four", AlbumID: album.ID, DateShot: time.Unix(1632758500, 0), Rating: 4},
		{Title: "five", Path: "/photos/five", AlbumID: album.ID, DateShot: time.Unix(1632758600, 0), Rating: 5},
	}
	require.NoError(t, db.Save(&media).Error)

	titles := func(minRating *int) []string {
		result, err := actions.MyTimeline(db, user, nil, nil, nil, minRating)
		require.NoError(t, err)

		titles := make([]string, 0, len(result))
		for _, m := range result {
			titles = append(titles, m.Title)
		}
		return titles
	}

	intPtr := func(v int) *int { return &v }

	assert.ElementsMatch(t, []string{"unrated", "two", "four", "five"}, titles(nil))
	assert.ElementsMatch(t, []string{"unrated", "four", "five"}, titles(intPtr(3)))
	assert.ElementsMatch(t, []string{"four", "five"}, titles(intPtr(4)))
	assert.ElementsMatch(t, []string{"five"}, titles(intPtr(5)))

	// Filters are clamped to [3, 5]
	assert.ElementsMatch(t, []string{"unrated", "four", "five"}, titles(intPtr(1)))
	assert.ElementsMatch(t, []string{"five"}, titles(intPtr(9)))
}
