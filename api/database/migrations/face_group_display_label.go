package migrations

import (
	"fmt"

	"github.com/photoview/photoview/api/graphql/models"
	"gorm.io/gorm"
)

// MigrateFaceGroupDisplayLabels gives face groups without a display label their label,
// for example groups created before display labels existed.
//
// Fork-specific: not part of upstream Photoview.
func MigrateFaceGroupDisplayLabels(db *gorm.DB) error {
	if err := db.Model(&models.FaceGroup{}).
		Where("display_label IS NULL AND label IS NOT NULL").
		UpdateColumn("display_label", gorm.Expr("label")).Error; err != nil {
		return fmt.Errorf("failed to set display labels of face groups: %w", err)
	}
	return nil
}
