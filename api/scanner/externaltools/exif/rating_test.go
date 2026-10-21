package exif

import (
	"fmt"
	"testing"
)

// Fork-specific: not part of upstream Photoview.
func TestParseRating(t *testing.T) {
	resetForTest()

	cleanup, err := Initialize()
	if err != nil {
		t.Fatalf("Initialize() error: %v", err)
	}
	defer cleanup()

	tests := []struct {
		file string
		want *int64
	}{
		{"./test_data/rated.jpg", new(int64(4))},
		{"./test_data/bird.jpg", nil},
	}

	for _, tc := range tests {
		t.Run(tc.file, func(t *testing.T) {
			metadata, err := Parse(tc.file)
			if err != nil {
				t.Fatalf("Parse(%q) returns an error: %v", tc.file, err)
			}

			got := metadata.Rating
			if (got == nil) != (tc.want == nil) || (got != nil && *got != *tc.want) {
				t.Errorf("Parse(%q).Rating = %v, want %v", tc.file, ptrString(got), ptrString(tc.want))
			}
		})
	}
}

func ptrString(v *int64) string {
	if v == nil {
		return "nil"
	}
	return fmt.Sprint(*v)
}
