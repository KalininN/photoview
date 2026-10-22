package migrations_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/photoview/photoview/api/database/migrations"
	"github.com/photoview/photoview/api/graphql/models"
	"github.com/photoview/photoview/api/test_utils"
)

func TestFaceGroupDisplayLabelMigration(t *testing.T) {
	db := test_utils.DatabaseTest(t)
	defer db.Exec("DELETE FROM face_groups")

	label, display := "Anna Smith", "Anna"
	groups := []models.FaceGroup{
		{Label: &label},
		{Label: &label, DisplayLabel: &display},
		{},
	}
	require.NoError(t, db.Create(&groups).Error)

	require.NoError(t, migrations.MigrateFaceGroupDisplayLabels(db))

	var results []models.FaceGroup
	require.NoError(t, db.Order("id").Find(&results).Error)
	require.Len(t, results, 3)

	assert.Equal(t, "Anna Smith", *results[0].DisplayLabel)
	assert.Equal(t, "Anna", *results[1].DisplayLabel)
	assert.Nil(t, results[2].DisplayLabel)
}
