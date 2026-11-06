package guessgame

import (
	"crypto/rand"
	"encoding/base64"
	mathrand "math/rand/v2"
	"slices"
	"time"
)

// PickMedia picks n distinct media ids from pool for a new game.
//
// Media last played before cutoff, or never, are preferred. When there are fewer than n of
// those, the rest is filled with recently played media, the longest ago first.
// lastPlayed maps a media id to when it was last played on the device starting the game.
// The result is in random order; it is shorter than n only when the pool is.
func PickMedia(pool []int, lastPlayed map[int]time.Time, cutoff time.Time, n int, rng *mathrand.Rand) []int {
	var fresh, recent []int
	for _, id := range pool {
		if played, ok := lastPlayed[id]; ok && !played.Before(cutoff) {
			recent = append(recent, id)
		} else {
			fresh = append(fresh, id)
		}
	}

	picked := sample(fresh, n, rng)

	if missing := n - len(picked); missing > 0 {
		slices.SortStableFunc(recent, func(a, b int) int {
			return lastPlayed[a].Compare(lastPlayed[b])
		})
		picked = append(picked, recent[:min(missing, len(recent))]...)
	}

	rng.Shuffle(len(picked), func(i, j int) {
		picked[i], picked[j] = picked[j], picked[i]
	})
	return picked
}

// sample returns up to n random elements of ids, using a partial Fisher-Yates shuffle on a copy.
func sample(ids []int, n int, rng *mathrand.Rand) []int {
	ids = slices.Clone(ids)
	n = min(n, len(ids))
	for i := range n {
		j := i + rng.IntN(len(ids)-i)
		ids[i], ids[j] = ids[j], ids[i]
	}
	return ids[:n]
}

// NewPublicID returns a random 128-bit id, encoded as 22 URL-safe characters.
// Games and rounds are addressed by these ids, so people sharing an account cannot
// reach each other's games by counting.
func NewPublicID() string {
	b := make([]byte, 16)
	// rand.Read never returns an error, it crashes the program instead.
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}
