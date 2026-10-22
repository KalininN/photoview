package models

import "gorm.io/gorm"

// Bounds that a minRating filter is clamped to. Media rated below
// LowestMinRating are never shown by a rating filter.
const (
	LowestMinRating  = 3
	HighestMinRating = 5
)

// MapMinRating is the minimum rating of the media shown on the map.
const MapMinRating = 4

// ClampMinRating limits a minRating filter to [LowestMinRating, HighestMinRating].
func ClampMinRating(minRating int) int {
	return max(LowestMinRating, min(minRating, HighestMinRating))
}

// FilterMinRating limits a media query to media rated at least minRating, when it is set.
//
// Fork-specific: not part of upstream Photoview.
func FilterMinRating(query *gorm.DB, minRating *int) *gorm.DB {
	return FilterMinRatingColumn(query, "media.rating", minRating)
}

// FilterMinRatingColumn is FilterMinRating for queries where the rating column has another name.
func FilterMinRatingColumn(query *gorm.DB, column string, minRating *int) *gorm.DB {
	if minRating == nil {
		return query
	}

	return query.Where(column+" >= ?", ClampMinRating(*minRating))
}
