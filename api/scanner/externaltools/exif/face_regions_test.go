package exif

import (
	"reflect"
	"testing"
)

// Fork-specific: not part of upstream Photoview.
func TestFaceRegionFields(t *testing.T) {
	resetForTest()

	cleanup, err := Initialize()
	if err != nil {
		t.Fatalf("Initialize() error: %v", err)
	}
	defer cleanup()

	fields, err := FaceRegionFields("./test_data/faces.jpg")
	if err != nil {
		t.Fatalf("FaceRegionFields() error: %v", err)
	}

	want := map[string]any{
		"RegionName":  []any{"Anna Smith", 2024.0},
		"RegionType":  []any{"Face", "Face"},
		"RegionAreaX": []any{0.25, 0.6},
		"RegionAreaY": []any{0.5, 0.4},
		"RegionAreaW": []any{0.1, 0.1},
		"RegionAreaH": []any{0.2, 0.1},
	}
	for key, value := range want {
		if got := fields[key]; !reflect.DeepEqual(got, value) {
			t.Errorf("FaceRegionFields()[%q] = %#v, want %#v", key, got, value)
		}
	}

	fields, err = FaceRegionFields("./test_data/bird.jpg")
	if err != nil {
		t.Fatalf("FaceRegionFields() error: %v", err)
	}
	if _, ok := fields["RegionName"]; ok {
		t.Errorf("FaceRegionFields(bird.jpg) has a RegionName: %#v", fields["RegionName"])
	}
}
