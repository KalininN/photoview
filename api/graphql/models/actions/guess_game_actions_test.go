package actions_test

// Fork-specific: not part of upstream Photoview.

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/photoview/photoview/api/database/drivers"
	"github.com/photoview/photoview/api/graphql/models"
	"github.com/photoview/photoview/api/graphql/models/actions"
	"github.com/photoview/photoview/api/guessgame"
	"github.com/photoview/photoview/api/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type guessPhoto struct {
	title     string
	lat, lon  *float64
	dateShot  *time.Time
	rating    int64
	mediaType models.MediaType
	noURL     bool
}

func floatPtr(v float64) *float64 { return &v }

func datePtr(year int, month time.Month) *time.Time {
	date := time.Date(year, month, 15, 12, 0, 0, 0, time.UTC)
	return &date
}

// eligiblePhoto is a photo every game can use: GPS, an EXIF date, rating 4 and a high-res image.
func eligiblePhoto(title string, lat, lon float64, year int, month time.Month) guessPhoto {
	return guessPhoto{title: title, lat: &lat, lon: &lon, dateShot: datePtr(year, month), rating: 4}
}

func createGuessUser(t *testing.T, db *gorm.DB, name string) (*models.User, *models.Album) {
	t.Helper()

	user, err := models.RegisterUser(db, name, nil, false)
	require.NoError(t, err)

	album := models.Album{Title: name, Path: "/photos/" + name}
	require.NoError(t, db.Save(&album).Error)
	require.NoError(t, db.Model(user).Association("Albums").Append(&album))

	return user, &album
}

func createGuessPhotos(t *testing.T, db *gorm.DB, album *models.Album, photos ...guessPhoto) map[string]*models.Media {
	t.Helper()

	result := make(map[string]*models.Media, len(photos))
	for _, photo := range photos {
		mediaType := photo.mediaType
		if mediaType == "" {
			mediaType = models.MediaTypePhoto
		}

		media := models.Media{
			Title:    photo.title,
			Path:     fmt.Sprintf("%s/%s.jpg", album.Path, photo.title),
			AlbumID:  album.ID,
			DateShot: time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC), // file date, never the answer
			Type:     mediaType,
			Rating:   photo.rating,
			Exif: &models.MediaEXIF{
				GPSLatitude:  photo.lat,
				GPSLongitude: photo.lon,
				DateShot:     photo.dateShot,
			},
		}
		require.NoError(t, db.Create(&media).Error)

		if !photo.noURL {
			require.NoError(t, db.Create(&models.MediaURL{
				MediaID:     media.ID,
				MediaName:   fmt.Sprintf("highres_%s_%d.jpg", photo.title, media.ID),
				Width:       1024,
				Height:      768,
				Purpose:     models.PhotoHighRes,
				ContentType: "image/jpeg",
			}).Error)
		}

		result[photo.title] = &media
	}

	return result
}

func startGuessGame(t *testing.T, db *gorm.DB, user *models.User, deviceKey string, rounds int) *models.GuessGame {
	t.Helper()

	game, err := actions.StartGuessGame(db, user, actions.StartGuessGameArgs{
		DeviceKey:  deviceKey,
		Rounds:     rounds,
		Difficulty: models.GuessDifficultyRegion,
	})
	require.NoError(t, err)
	return game
}

func guessRounds(t *testing.T, db *gorm.DB, game *models.GuessGame) []*models.GuessRound {
	t.Helper()

	rounds, err := actions.GuessGameRounds(db, game)
	require.NoError(t, err)
	return rounds
}

func exactGuess(round *models.GuessRound) models.GuessInput {
	return models.GuessInput{
		Latitude:  round.AnswerLatitude,
		Longitude: round.AnswerLongitude,
		Year:      round.AnswerYear,
		Month:     round.AnswerMonth,
	}
}

func TestGuessGameInfo(t *testing.T) {
	db := test_utils.DatabaseTest(t)

	user, album := createGuessUser(t, db, "family")
	_, otherAlbum := createGuessUser(t, db, "neighbours")

	createGuessPhotos(t, db, album,
		eligiblePhoto("berlin", 52.52, 13.40, 2015, time.March),
		eligiblePhoto("tokyo", 35.68, 139.65, 2019, time.July),
		eligiblePhoto("sydney", -33.87, 151.21, 2021, time.December),
		guessPhoto{title: "rated-three", lat: floatPtr(48.85), lon: floatPtr(2.35), dateShot: datePtr(2010, time.May), rating: 3},
		guessPhoto{title: "rated-five", lat: floatPtr(40.71), lon: floatPtr(-74.00), dateShot: datePtr(2012, time.May), rating: 5},
		guessPhoto{title: "no-gps", dateShot: datePtr(2016, time.May), rating: 5},
		guessPhoto{title: "no-date", lat: floatPtr(1), lon: floatPtr(1), rating: 5},
		guessPhoto{title: "null-island", lat: floatPtr(0), lon: floatPtr(0), dateShot: datePtr(2016, time.May), rating: 5},
		guessPhoto{title: "too-old", lat: floatPtr(1), lon: floatPtr(1), dateShot: datePtr(1850, time.May), rating: 5},
		guessPhoto{title: "video", lat: floatPtr(1), lon: floatPtr(1), dateShot: datePtr(2016, time.May), rating: 5,
			mediaType: models.MediaTypeVideo},
		guessPhoto{title: "no-image", lat: floatPtr(1), lon: floatPtr(1), dateShot: datePtr(2016, time.May), rating: 5,
			noURL: true},
	)
	createGuessPhotos(t, db, otherAlbum, eligiblePhoto("not-ours", 10, 10, 1990, time.January))

	t.Run("default filter", func(t *testing.T) {
		info, err := actions.GuessGameInfo(db, user, nil)
		require.NoError(t, err)

		assert.Equal(t, 4, info.EligibleCount, "berlin, tokyo, sydney and rated-five")
		assert.Equal(t, 2012, *info.MinYear)
		assert.Equal(t, 2021, *info.MaxYear)
		assert.Equal(t, models.Coordinates{Latitude: -33.87, Longitude: -74.00}, *info.SouthWest)
		assert.Equal(t, models.Coordinates{Latitude: 52.52, Longitude: 151.21}, *info.NorthEast)
	})

	t.Run("min rating", func(t *testing.T) {
		minRating := 3
		info, err := actions.GuessGameInfo(db, user, &models.GuessGameFilter{MinRating: &minRating})
		require.NoError(t, err)
		assert.Equal(t, 5, info.EligibleCount)
		assert.Equal(t, 2010, *info.MinYear)

		minRating = 9 // clamped to 5
		info, err = actions.GuessGameInfo(db, user, &models.GuessGameFilter{MinRating: &minRating})
		require.NoError(t, err)
		assert.Equal(t, 1, info.EligibleCount)
	})

	t.Run("album filter includes sub-albums", func(t *testing.T) {
		child := models.Album{Title: "trip", Path: album.Path + "/trip", ParentAlbumID: &album.ID}
		require.NoError(t, db.Save(&child).Error)
		require.NoError(t, db.Model(user).Association("Albums").Append(&child))
		createGuessPhotos(t, db, &child, eligiblePhoto("trip-photo", 60, 25, 2023, time.June))

		info, err := actions.GuessGameInfo(db, user, &models.GuessGameFilter{AlbumIds: []int{child.ID}})
		require.NoError(t, err)
		assert.Equal(t, 1, info.EligibleCount)

		info, err = actions.GuessGameInfo(db, user, &models.GuessGameFilter{AlbumIds: []int{album.ID}})
		require.NoError(t, err)
		assert.Equal(t, 5, info.EligibleCount)

		info, err = actions.GuessGameInfo(db, user, &models.GuessGameFilter{AlbumIds: []int{otherAlbum.ID}})
		require.NoError(t, err)
		assert.Equal(t, 0, info.EligibleCount, "another account's album gives nothing")
		assert.Nil(t, info.MinYear)
	})

	t.Run("disabled", func(t *testing.T) {
		t.Setenv("PHOTOVIEW_DISABLE_4D_GUESSER", "1")
		_, err := actions.GuessGameInfo(db, user, nil)
		assert.ErrorIs(t, err, actions.ErrGuesserDisabled)
	})
}

func TestStartGuessGame(t *testing.T) {
	db := test_utils.DatabaseTest(t)

	user, album := createGuessUser(t, db, "family")
	photos := createGuessPhotos(t, db, album,
		eligiblePhoto("berlin", 52.52, 13.40, 2015, time.March),
		eligiblePhoto("tokyo", 35.68, 139.65, 2019, time.July),
		eligiblePhoto("sydney", -33.87, 151.21, 2021, time.December),
		eligiblePhoto("paris", 48.85, 2.35, 2010, time.May),
		guessPhoto{title: "rated-three", lat: floatPtr(40.71), lon: floatPtr(-74.00), dateShot: datePtr(2012, time.May), rating: 3},
	)

	t.Run("picks distinct eligible photos and stores the answers", func(t *testing.T) {
		game := startGuessGame(t, db, user, "device-a", 4)

		assert.Len(t, game.PublicID, 22)
		assert.Equal(t, 4, game.RoundCount)
		assert.Equal(t, 40000, game.MaxScore())
		assert.Equal(t, string(models.GuessDifficultyRegion), game.Difficulty)
		assert.Equal(t, guessgame.Scales{DistanceKm: 150, TimeMonths: 37}, game.Scales())
		assert.Equal(t, 4, game.MinRating)

		rounds := guessRounds(t, db, game)
		require.Len(t, rounds, 4)

		seen := map[int]bool{}
		for i, round := range rounds {
			assert.Equal(t, i, round.RoundIndex)
			assert.Len(t, round.PublicID, 22)
			assert.False(t, round.Guessed())
			assert.Nil(t, round.Result())

			require.NotNil(t, round.MediaID)
			assert.False(t, seen[*round.MediaID], "photo picked twice")
			seen[*round.MediaID] = true
			assert.NotEqual(t, photos["rated-three"].ID, *round.MediaID)

			var exif models.MediaEXIF
			require.NoError(t, db.Joins("JOIN media ON media.exif_id = media_exif.id").
				Where("media.id = ?", *round.MediaID).First(&exif).Error)
			assert.Equal(t, *exif.GPSLatitude, round.AnswerLatitude)
			assert.Equal(t, *exif.GPSLongitude, round.AnswerLongitude)
			assert.Equal(t, exif.DateShot.Year(), round.AnswerYear)
			assert.Equal(t, int(exif.DateShot.Month()), round.AnswerMonth)
		}
	})

	t.Run("invalid arguments", func(t *testing.T) {
		_, err := actions.StartGuessGame(db, user, actions.StartGuessGameArgs{DeviceKey: "", Rounds: 3})
		assert.Error(t, err)

		_, err = actions.StartGuessGame(db, user, actions.StartGuessGameArgs{DeviceKey: "device-a", Rounds: 0})
		assert.Error(t, err)

		_, err = actions.StartGuessGame(db, user, actions.StartGuessGameArgs{DeviceKey: "device-a", Rounds: 11})
		assert.Error(t, err)

		_, err = actions.StartGuessGame(db, user, actions.StartGuessGameArgs{DeviceKey: "device-a", Rounds: 5})
		assert.ErrorContains(t, err, "only 4 photos")
	})

	t.Run("unknown difficulty falls back to auto", func(t *testing.T) {
		game, err := actions.StartGuessGame(db, user, actions.StartGuessGameArgs{
			DeviceKey:  "device-a",
			Rounds:     1,
			Difficulty: "HARDEST",
		})
		require.NoError(t, err)
		assert.Equal(t, string(models.GuessDifficultyAuto), game.Difficulty)
	})

	t.Run("disabled", func(t *testing.T) {
		t.Setenv("PHOTOVIEW_DISABLE_4D_GUESSER", "true")
		_, err := actions.StartGuessGame(db, user, actions.StartGuessGameArgs{DeviceKey: "device-a", Rounds: 1})
		assert.ErrorIs(t, err, actions.ErrGuesserDisabled)
	})
}

func TestGuessGameRepeatsPerDevice(t *testing.T) {
	db := test_utils.DatabaseTest(t)

	user, album := createGuessUser(t, db, "family")
	createGuessPhotos(t, db, album,
		eligiblePhoto("one", 1, 1, 2001, time.January),
		eligiblePhoto("two", 2, 2, 2002, time.February),
		eligiblePhoto("three", 3, 3, 2003, time.March),
		eligiblePhoto("four", 4, 4, 2004, time.April),
		eligiblePhoto("five", 5, 5, 2005, time.May),
		eligiblePhoto("six", 6, 6, 2006, time.June),
	)

	mediaIDs := func(game *models.GuessGame) []int {
		var ids []int
		for _, round := range guessRounds(t, db, game) {
			ids = append(ids, *round.MediaID)
		}
		return ids
	}

	first := mediaIDs(startGuessGame(t, db, user, "device-a", 3))
	second := mediaIDs(startGuessGame(t, db, user, "device-a", 3))
	assert.Empty(t, intersect(first, second), "the device gets the photos it has not played")

	third := mediaIDs(startGuessGame(t, db, user, "device-a", 2))
	assert.Len(t, third, 2, "all photos played, falls back to the longest ago played")
	assert.Subset(t, first, third)

	// Another device's history does not limit this one: over many games it sees every photo.
	seen := map[int]bool{}
	for range 20 {
		for _, id := range mediaIDs(startGuessGame(t, db, user, fmt.Sprintf("device-b-%d", len(seen)), 3)) {
			seen[id] = true
		}
	}
	assert.Len(t, seen, 6)
}

func intersect(a, b []int) []int {
	inA := map[int]bool{}
	for _, id := range a {
		inA[id] = true
	}

	var both []int
	for _, id := range b {
		if inA[id] {
			both = append(both, id)
		}
	}
	return both
}

func TestSubmitGuess(t *testing.T) {
	db := test_utils.DatabaseTest(t)

	user, album := createGuessUser(t, db, "family")
	stranger, _ := createGuessUser(t, db, "neighbours")
	createGuessPhotos(t, db, album,
		eligiblePhoto("berlin", 52.52, 13.40, 2015, time.March),
		eligiblePhoto("tokyo", 35.68, 139.65, 2019, time.July),
		eligiblePhoto("sydney", -33.87, 151.21, 2021, time.December),
	)

	t.Run("plays a game to the end", func(t *testing.T) {
		game := startGuessGame(t, db, user, "device-a", 3)
		rounds := guessRounds(t, db, game)

		_, err := actions.SubmitGuess(db, user, rounds[1].PublicID, exactGuess(rounds[1]))
		assert.ErrorIs(t, err, actions.ErrGuessOutOfOrder)

		// Round 0: exact.
		result, err := actions.SubmitGuess(db, user, rounds[0].PublicID, exactGuess(rounds[0]))
		require.NoError(t, err)
		assert.Equal(t, 10000, result.Score)
		assert.Equal(t, 0, result.MonthsOff)
		assert.InDelta(t, 0, result.DistanceKm, 1e-6)

		// Round 1: 1 degree of latitude (about 111 km) and 12 months off.
		guess := exactGuess(rounds[1])
		guess.Latitude += 1
		guess.Year -= 1
		result, err = actions.SubmitGuess(db, user, rounds[1].PublicID, guess)
		require.NoError(t, err)
		assert.InDelta(t, 111.2, result.DistanceKm, 0.5)
		assert.Equal(t, 12, result.MonthsOff)
		assert.Equal(t, guessgame.LocationScore(result.DistanceKm, game.Scales()), result.LocationScore)
		assert.Equal(t, 4501, result.TimeScore, "a year off keeps 90%")
		assert.Equal(t, result.LocationScore+result.TimeScore, result.Score)
		assert.Equal(t, guess.Latitude, result.Guess.Latitude)
		assert.Equal(t, rounds[1].AnswerLatitude, result.Answer.Latitude)
		total := 10000 + result.Score

		// Guessing a round again returns the stored result.
		again, err := actions.SubmitGuess(db, user, rounds[1].PublicID, exactGuess(rounds[1]))
		require.NoError(t, err)
		assert.Equal(t, result, again)

		reloaded, err := actions.GetGuessGame(db, user, game.PublicID)
		require.NoError(t, err)
		assert.Equal(t, total, reloaded.TotalScore)
		assert.Nil(t, reloaded.FinishedAt)

		// Round 2 finishes the game.
		result, err = actions.SubmitGuess(db, user, rounds[2].PublicID, exactGuess(rounds[2]))
		require.NoError(t, err)
		total += result.Score

		reloaded, err = actions.GetGuessGame(db, user, game.PublicID)
		require.NoError(t, err)
		assert.Equal(t, total, reloaded.TotalScore)
		assert.NotNil(t, reloaded.FinishedAt)

		for _, round := range guessRounds(t, db, reloaded) {
			assert.NotNil(t, round.Result())
		}
	})

	t.Run("other accounts cannot reach the game", func(t *testing.T) {
		game := startGuessGame(t, db, user, "device-a", 1)
		rounds := guessRounds(t, db, game)

		_, err := actions.GetGuessGame(db, stranger, game.PublicID)
		assert.ErrorIs(t, err, actions.ErrGuessGameNotFound)

		_, err = actions.SubmitGuess(db, stranger, rounds[0].PublicID, exactGuess(rounds[0]))
		assert.ErrorIs(t, err, actions.ErrGuessNotFound)

		_, err = actions.GetGuessGame(db, user, "unknown")
		assert.ErrorIs(t, err, actions.ErrGuessGameNotFound)

		_, err = actions.GetGuessGame(db, user, fmt.Sprint(game.ID))
		assert.ErrorIs(t, err, actions.ErrGuessGameNotFound, "the database id is not an address")
	})

	t.Run("invalid guesses", func(t *testing.T) {
		game := startGuessGame(t, db, user, "device-a", 1)
		round := guessRounds(t, db, game)[0]

		for name, change := range map[string]func(*models.GuessInput){
			"latitude":  func(g *models.GuessInput) { g.Latitude = 91 },
			"longitude": func(g *models.GuessInput) { g.Longitude = -181 },
			"month":     func(g *models.GuessInput) { g.Month = 13 },
			"month 0":   func(g *models.GuessInput) { g.Month = 0 },
			"year":      func(g *models.GuessInput) { g.Year = 0 },
		} {
			guess := exactGuess(round)
			change(&guess)
			_, err := actions.SubmitGuess(db, user, round.PublicID, guess)
			assert.Error(t, err, name)
		}
	})

	t.Run("a replaced game is closed", func(t *testing.T) {
		old := startGuessGame(t, db, user, "device-a", 1)
		other := startGuessGame(t, db, user, "device-b", 1)

		replaceID := old.PublicID
		_, err := actions.StartGuessGame(db, user, actions.StartGuessGameArgs{
			DeviceKey:     "device-a",
			Rounds:        1,
			ReplaceGameID: &replaceID,
		})
		require.NoError(t, err)

		old, err = actions.GetGuessGame(db, user, old.PublicID)
		require.NoError(t, err)
		assert.True(t, old.Abandoned)

		round := guessRounds(t, db, old)[0]
		_, err = actions.SubmitGuess(db, user, round.PublicID, exactGuess(round))
		assert.ErrorIs(t, err, actions.ErrGuessGameClosed)

		// Other devices' games are untouched.
		other, err = actions.GetGuessGame(db, user, other.PublicID)
		require.NoError(t, err)
		assert.False(t, other.Abandoned)
	})

	t.Run("deleted photo keeps its answer", func(t *testing.T) {
		game := startGuessGame(t, db, user, "device-c", 1)
		round := guessRounds(t, db, game)[0]

		require.NoError(t, db.Delete(&models.Media{}, *round.MediaID).Error)

		result, err := actions.SubmitGuess(db, user, round.PublicID, exactGuess(round))
		require.NoError(t, err)
		assert.Equal(t, 10000, result.Score)

		round = guessRounds(t, db, game)[0]
		assert.Nil(t, round.MediaID)
	})
}

func TestSimultaneousGuessGames(t *testing.T) {
	db := test_utils.DatabaseTest(t)

	user, album := createGuessUser(t, db, "family")
	createGuessPhotos(t, db, album,
		eligiblePhoto("berlin", 52.52, 13.40, 2015, time.March),
		eligiblePhoto("tokyo", 35.68, 139.65, 2019, time.July),
		eligiblePhoto("sydney", -33.87, 151.21, 2021, time.December),
	)

	t.Run("interleaved games do not touch each other", func(t *testing.T) {
		gameA := startGuessGame(t, db, user, "device-a", 2)
		gameB := startGuessGame(t, db, user, "device-b", 2)
		roundsA := guessRounds(t, db, gameA)
		roundsB := guessRounds(t, db, gameB)

		far := func(round *models.GuessRound) models.GuessInput {
			guess := exactGuess(round)
			guess.Latitude = -guess.Latitude
			guess.Year -= 10
			return guess
		}

		_, err := actions.SubmitGuess(db, user, roundsA[0].PublicID, exactGuess(roundsA[0]))
		require.NoError(t, err)
		resultB, err := actions.SubmitGuess(db, user, roundsB[0].PublicID, far(roundsB[0]))
		require.NoError(t, err)
		_, err = actions.SubmitGuess(db, user, roundsA[1].PublicID, exactGuess(roundsA[1]))
		require.NoError(t, err)

		gameA, err = actions.GetGuessGame(db, user, gameA.PublicID)
		require.NoError(t, err)
		gameB, err = actions.GetGuessGame(db, user, gameB.PublicID)
		require.NoError(t, err)

		assert.Equal(t, 20000, gameA.TotalScore)
		assert.NotNil(t, gameA.FinishedAt)
		assert.Equal(t, resultB.Score, gameB.TotalScore)
		assert.Nil(t, gameB.FinishedAt)
	})

	t.Run("concurrent guesses on one round store one guess", func(t *testing.T) {
		game := startGuessGame(t, db, user, "device-c", 1)
		round := guessRounds(t, db, game)[0]

		const tabs = 4
		results := make([]*models.GuessRoundResult, tabs)
		errs := make([]error, tabs)

		var wg sync.WaitGroup
		for i := range tabs {
			wg.Add(1)
			go func() {
				defer wg.Done()
				guess := exactGuess(round)
				guess.Year -= i // every tab guesses a different year
				results[i], errs[i] = actions.SubmitGuess(db, user, round.PublicID, guess)
			}()
		}
		wg.Wait()

		var stored *models.GuessRoundResult
		for i := range tabs {
			if errs[i] != nil && drivers.GetDatabaseDriverType(db) == drivers.SQLITE {
				// SQLite refuses a write that overlaps the winning one, before the winner commits.
				// The tab then reports an error and its next request gets the stored result.
				assert.ErrorContains(t, errs[i], "database is locked", "tab %d", i)
				continue
			}
			require.NoError(t, errs[i], "tab %d", i)
			if stored == nil {
				stored = results[i]
			}
			assert.Equal(t, stored, results[i], "every tab gets the same stored result")
		}
		require.NotNil(t, stored, "at least one tab succeeds")

		again, err := actions.SubmitGuess(db, user, round.PublicID, exactGuess(round))
		require.NoError(t, err)
		assert.Equal(t, stored, again)

		game, err = actions.GetGuessGame(db, user, game.PublicID)
		require.NoError(t, err)
		assert.Equal(t, stored.Score, game.TotalScore)
	})
}

func TestGuessGameOpenLimit(t *testing.T) {
	db := test_utils.DatabaseTest(t)

	user, album := createGuessUser(t, db, "family")
	createGuessPhotos(t, db, album, eligiblePhoto("berlin", 52.52, 13.40, 2015, time.March))

	// Fill the account with open games, half of them idle for longer than the timeout.
	games := make([]models.GuessGame, actions.GuessMaxOpenGames)
	for i := range games {
		updatedAt := time.Now()
		if i%2 == 0 {
			updatedAt = updatedAt.Add(-actions.GuessIdleTimeout - time.Hour)
		}
		games[i] = models.GuessGame{
			PublicID:   guessgame.NewPublicID(),
			UserID:     user.ID,
			DeviceKey:  fmt.Sprintf("device-%d", i),
			RoundCount: 1,
			Difficulty: string(models.GuessDifficultyAuto),
			ModelTimestamps: models.ModelTimestamps{
				CreatedAt: updatedAt,
				UpdatedAt: updatedAt,
			},
		}
	}
	require.NoError(t, db.Create(&games).Error)

	startGuessGame(t, db, user, "device-new", 1)

	var abandoned int64
	require.NoError(t, db.Model(&models.GuessGame{}).Where("abandoned = ?", true).Count(&abandoned).Error)
	assert.Equal(t, int64(actions.GuessMaxOpenGames/2), abandoned, "only idle games are abandoned")

	// Fill up again with active games: now nothing is idle, so a new game is refused.
	more := make([]models.GuessGame, actions.GuessMaxOpenGames/2-1)
	for i := range more {
		more[i] = models.GuessGame{
			PublicID:   guessgame.NewPublicID(),
			UserID:     user.ID,
			DeviceKey:  "device-more",
			RoundCount: 1,
			Difficulty: string(models.GuessDifficultyAuto),
		}
	}
	require.NoError(t, db.Create(&more).Error)

	_, err := actions.StartGuessGame(db, user, actions.StartGuessGameArgs{DeviceKey: "device-new", Rounds: 1})
	assert.ErrorIs(t, err, actions.ErrGuessTooManyGames)
}
