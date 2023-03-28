package sbm

const (
	ErrDimension = "array dimension error"
)

// Sbm is a Simple Bit Map.
type Sbm struct {
	format     SbmFormat
	pixelArray SbmPixelArray
}
