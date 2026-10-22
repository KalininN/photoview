//go:build no_face_detection || !faces_exif

package face_detection

// FromMetadata reports that faces are read from image metadata rather than detected.
//
// Fork-specific: not part of upstream Photoview.
const FromMetadata = false
