// SbmPixelArrayMetaData.go.

package sbm

// SBM Internal Data Model: Pixel Array Meta-Data.
type SbmPixelArrayMetaData struct {

	// Fixed Values.
	width  uint
	height uint
	area   uint

	// Values used in top and bottom Headers.
	header SbmPixelArrayMetaDataHeader
}
