package face_detection

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFaceRegions(t *testing.T) {
	t.Run("single region as scalars", func(t *testing.T) {
		fields := map[string]interface{}{
			"RegionName":  "Anna Smith",
			"RegionType":  "Face",
			"RegionAreaX": 0.25,
			"RegionAreaY": 0.5,
			"RegionAreaW": 0.1,
			"RegionAreaH": 0.2,
		}

		assert.Equal(t, []FaceRegion{{Name: "Anna Smith", X: 0.25, Y: 0.5, Width: 0.1, Height: 0.2}},
			ParseFaceRegions(fields, 1))
	})

	t.Run("skips unnamed, non-face and incomplete regions", func(t *testing.T) {
		fields := map[string]interface{}{
			"RegionName":  []interface{}{"Anna", "", "Pet", "  ", 2024.0, "Bob"},
			"RegionType":  []interface{}{"Face", "Face", "Pet", "Face", "Face", "Face"},
			"RegionAreaX": []interface{}{0.1, 0.2, 0.3, 0.4, 0.5},
			"RegionAreaY": []interface{}{0.1, 0.2, 0.3, 0.4, 0.5},
			"RegionAreaW": []interface{}{0.1, 0.1, 0.1, 0.1, 0.1},
			"RegionAreaH": []interface{}{0.1, 0.1, 0.1, 0.1, 0.1},
		}

		var names []string
		for _, face := range ParseFaceRegions(fields, 1) {
			names = append(names, face.Name)
		}
		// "Bob" has no area
		assert.Equal(t, []string{"Anna", "2024"}, names)
	})

	t.Run("no regions", func(t *testing.T) {
		assert.Empty(t, ParseFaceRegions(map[string]interface{}{}, 1))
		assert.Empty(t, ParseFaceRegions(map[string]interface{}{"RegionType": "Face"}, 1))
	})

	t.Run("rotated", func(t *testing.T) {
		fields := map[string]interface{}{
			"RegionName":  "Anna",
			"RegionType":  "Face",
			"RegionAreaX": 0.25,
			"RegionAreaY": 0.5,
			"RegionAreaW": 0.1,
			"RegionAreaH": 0.2,
		}

		// Orientation 6: rotated 90° clockwise for display
		assert.Equal(t, []FaceRegion{{Name: "Anna", X: 0.5, Y: 0.25, Width: 0.2, Height: 0.1}},
			ParseFaceRegions(fields, 6))
	})
}

func TestRotatePoint(t *testing.T) {
	tests := []struct {
		orientation int
		x, y, h, w  float64
	}{
		{orientation: 0, x: 0.2, y: 0.3, h: 0.4, w: 0.1},
		{orientation: 1, x: 0.2, y: 0.3, h: 0.4, w: 0.1},
		{orientation: 2, x: 0.8, y: 0.3, h: 0.4, w: 0.1},
		{orientation: 3, x: 0.8, y: 0.7, h: 0.4, w: 0.1},
		{orientation: 6, x: 0.7, y: 0.2, h: 0.1, w: 0.4},
		{orientation: 8, x: 0.3, y: 0.8, h: 0.1, w: 0.4},
	}

	for _, test := range tests {
		x, y, h, w := RotatePoint(0.2, 0.3, 0.4, 0.1, test.orientation)
		assert.InDelta(t, test.x, x, 1e-9, "x, orientation %d", test.orientation)
		assert.InDelta(t, test.y, y, 1e-9, "y, orientation %d", test.orientation)
		assert.InDelta(t, test.h, h, 1e-9, "h, orientation %d", test.orientation)
		assert.InDelta(t, test.w, w, 1e-9, "w, orientation %d", test.orientation)
	}
}
