package exif

// ForgetCachedMetadata drops the exiftool metadata cached for a file, so the next
// parse reads the file again. Needed when a file changed on disk.
//
// Fork-specific: not part of upstream Photoview.
func ForgetCachedMetadata(path string) {
	if parser, ok := globalExifParser.(*externalExifParser); ok {
		parser.dataLoader.Clear(path)
	}
}
