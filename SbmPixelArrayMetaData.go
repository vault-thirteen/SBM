package sbm

// SbmPixelArrayMetaData is meta-data for a pixel array.
type SbmPixelArrayMetaData struct {

	// Fixed values.
	width  uint
	height uint
	area   uint

	// Values used in top and bottom headers.
	header SbmPixelArrayMetaDataHeader
}
