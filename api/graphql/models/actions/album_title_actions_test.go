package actions_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/photoview/photoview/api/graphql/models/actions"
)

func TestSanitizeTitle(t *testing.T) {
	assert.Equal(t, "Summer 2024: Paris (day 1)", actions.SanitizeTitle("Summer 2024: Paris (day 1)"))
	assert.Equal(t, "Лето в Москве", actions.SanitizeTitle("Лето в Москве"))
	assert.Equal(t, "scriptalert(1)script", actions.SanitizeTitle("<script>alert(1)</script>"))
	assert.Equal(t, "", actions.SanitizeTitle("<>/\\"))
}
