package models

import "gorm.io/gorm"

// FilterMinRating limits a media query to media rated at least minRating, when it is set.
//
// Fork-specific: not part of upstream Photoview.
func FilterMinRating(query *gorm.DB, minRating *int) *gorm.DB {
	if minRating == nil {
		return query
	}

	return query.Where("media.rating >= ?", *minRating)
}
