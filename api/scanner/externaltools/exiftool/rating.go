package exiftool

// Rating stores the user-assigned star rating (EXIF/XMP `Rating`).
//
// Fork-specific: not part of upstream Photoview.
type Rating struct {
	Rating *int64
}
