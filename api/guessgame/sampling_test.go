package guessgame

import (
	mathrand "math/rand/v2"
	"testing"
	"time"

	_ "github.com/photoview/photoview/api/test_utils/flags"
	"github.com/stretchr/testify/assert"
)

func newRng(seed uint64) *mathrand.Rand {
	return mathrand.New(mathrand.NewPCG(seed, seed))
}

func TestPickMedia(t *testing.T) {
	now := time.Date(2026, 11, 6, 12, 0, 0, 0, time.UTC)
	cutoff := now.AddDate(0, 0, -30)

	t.Run("distinct ids from the pool", func(t *testing.T) {
		pool := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		for seed := range uint64(50) {
			picked := PickMedia(pool, nil, cutoff, 5, newRng(seed))
			assert.Len(t, picked, 5)
			assert.Subset(t, pool, picked)

			seen := map[int]bool{}
			for _, id := range picked {
				assert.False(t, seen[id], "duplicate id %d", id)
				seen[id] = true
			}
		}
	})

	t.Run("deterministic with a seeded generator", func(t *testing.T) {
		pool := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		assert.Equal(t, PickMedia(pool, nil, cutoff, 5, newRng(7)), PickMedia(pool, nil, cutoff, 5, newRng(7)))
	})

	t.Run("recently played media are left out", func(t *testing.T) {
		pool := []int{1, 2, 3, 4, 5, 6}
		lastPlayed := map[int]time.Time{
			1: now.AddDate(0, 0, -1),
			2: now.AddDate(0, 0, -2),
			3: now.AddDate(0, 0, -40), // played before the cutoff counts as fresh
		}
		for seed := range uint64(50) {
			picked := PickMedia(pool, lastPlayed, cutoff, 4, newRng(seed))
			assert.ElementsMatch(t, []int{3, 4, 5, 6}, picked)
		}
	})

	t.Run("a small pool falls back to the longest ago played", func(t *testing.T) {
		pool := []int{1, 2, 3, 4}
		lastPlayed := map[int]time.Time{
			1: now.AddDate(0, 0, -1),
			2: now.AddDate(0, 0, -10),
			3: now.AddDate(0, 0, -5),
		}
		picked := PickMedia(pool, lastPlayed, cutoff, 3, newRng(1))
		assert.ElementsMatch(t, []int{4, 2, 3}, picked)
	})

	t.Run("never more than the pool", func(t *testing.T) {
		picked := PickMedia([]int{1, 2}, nil, cutoff, 5, newRng(1))
		assert.ElementsMatch(t, []int{1, 2}, picked)
		assert.Empty(t, PickMedia(nil, nil, cutoff, 5, newRng(1)))
	})

	t.Run("the pool is not modified", func(t *testing.T) {
		pool := []int{1, 2, 3, 4, 5}
		PickMedia(pool, nil, cutoff, 3, newRng(3))
		assert.Equal(t, []int{1, 2, 3, 4, 5}, pool)
	})
}

func TestNewPublicID(t *testing.T) {
	seen := map[string]bool{}
	for range 1000 {
		id := NewPublicID()
		assert.Len(t, id, 22)
		assert.Regexp(t, `^[A-Za-z0-9_-]+$`, id)
		assert.False(t, seen[id])
		seen[id] = true
	}
}
