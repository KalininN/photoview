package exif

// Fork-specific: not part of upstream Photoview.

import (
	"fmt"
)

// FaceRegionFields returns the metadata tags of the file at `filepath` by name, with numeric values,
// for face_detection.ParseFaceRegions to read the named face regions (XMP MWG `RegionName`, `RegionArea*`) from.
func FaceRegionFields(filepath string) (map[string]any, error) {
	globalMu.Lock()
	defer globalMu.Unlock()

	if globalExifParser == nil {
		return nil, fmt.Errorf("no exif parser initialized")
	}

	fields := map[string]any{}
	if err := globalExifParser.QueryJSONTagsByNumber(filepath, &fields); err != nil {
		return nil, err
	}

	return fields, nil
}
