package face_detection

// Face regions read from image metadata (XMP MWG regions), used by the
// faces_exif face detector. Kept free of build tags, so it is always compiled and tested.
//
// Fork-specific: not part of upstream Photoview.

import (
	"strconv"
	"strings"
)

// FaceRegion is a named face in normalized image coordinates, centered on (X, Y).
type FaceRegion struct {
	Name   string
	X, Y   float64
	Width  float64
	Height float64
}

// ParseFaceRegions returns the named face regions in exiftool's metadata `fields`,
// rotated from the stored image to the displayed one by the EXIF `orientation`.
// Regions without a name, regions of other types and incomplete regions are skipped.
func ParseFaceRegions(fields map[string]interface{}, orientation int) []FaceRegion {
	names := toStringSlice(fields["RegionName"])
	types := toStringSlice(fields["RegionType"])
	xs := toFloatSlice(fields["RegionAreaX"])
	ys := toFloatSlice(fields["RegionAreaY"])
	ws := toFloatSlice(fields["RegionAreaW"])
	hs := toFloatSlice(fields["RegionAreaH"])

	count := min(len(names), len(types), len(xs), len(ys), len(ws), len(hs))

	var faces []FaceRegion
	for i := 0; i < count; i++ {
		name := strings.TrimSpace(names[i])
		if types[i] != "Face" || name == "" {
			continue
		}

		x, y, h, w := RotatePoint(xs[i], ys[i], hs[i], ws[i], orientation)
		faces = append(faces, FaceRegion{Name: name, X: x, Y: y, Width: w, Height: h})
	}

	return faces
}

// exiftool reports a single region as a scalar and several regions as a list.
func toStringSlice(raw interface{}) []string {
	var result []string
	switch v := raw.(type) {
	case []interface{}:
		for _, item := range v {
			result = append(result, toString(item))
		}
	case nil:
	default:
		result = []string{toString(v)}
	}
	return result
}

// toString also accepts numbers, since exiftool -n reports a name such as "2024" as one.
func toString(raw interface{}) string {
	switch v := raw.(type) {
	case string:
		return v
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	}
	return ""
}

func toFloatSlice(raw interface{}) []float64 {
	var result []float64
	switch v := raw.(type) {
	case float64:
		result = []float64{v}
	case []interface{}:
		for _, item := range v {
			f, ok := item.(float64)
			if !ok {
				// Keep the positions of the other lists aligned
				return nil
			}
			result = append(result, f)
		}
	}
	return result
}

// RotatePoint transforms normalized coordinates (x, y) and size (h, w)
// according to EXIF Orientation (1–8).
func RotatePoint(x, y, h, w float64, orientation int) (float64, float64, float64, float64) {
	var rotations int
	var flip bool

	switch orientation {
	case 2:
		rotations, flip = 0, true
	case 3:
		rotations, flip = 2, false
	case 4:
		rotations, flip = 2, true
	case 5:
		rotations, flip = 3, true
	case 6:
		rotations, flip = 1, false
	case 7:
		rotations, flip = 1, true
	case 8:
		rotations, flip = 3, false
	default:
		rotations, flip = 0, false
	}

	for i := 0; i < rotations; i++ {
		x, y = 1-y, x
		h, w = w, h
	}

	// Apply horizontal flip if needed
	if flip {
		x = 1 - x
	}

	return x, y, h, w
}
