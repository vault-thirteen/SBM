package sbm

// This Package provides a Model and Methods to work with the SBM Format.
// "SBM" is an acronym for "Simple Bit Map".

// The SBM Format is a Format for two-level (monochrome, black-&-white) raster
// graphical Images.

// The SBM Format is used to store a raw uncompressed Array of binary Pixels
// together with basic meta-data describing the size of the Pixel Array.
// SBM Format is stored in a mixed Encoding. This means that meta-data is
// encoded as plain ASCII Text Symbols and Pixel Array is encoded using the
// binary Format. The Array of Pixels is composed of Pixels Row by Row,
// starting with the top Row, and having the bottom Row at the End. Each Row
// is composed of W Pixels, where W is the Array's Width. Total Number of Rows
// is H, where H is the Array's Height. The Array contains A Pixels total,
// where A is the Array's Area, the multiple of W and H. Each Pixel in the
// Array is a separate Bit, where Zero Bit (0) is black (dark color) and One
// Bit (1) is white (light color).

// Due to the Limitations of current Hardware, the Order of Bits in each Byte
// is not controlled by this Library (Package). The least significant Bit is
// considered to be the first Bit, the most significant Bit is the last Bit.

// Errors.
const (
	ErrDimension = "Array Dimension Error"
)

// SBM Internal Data Model.
type Sbm struct {

	// Format Parameters.
	format SbmFormat

	// Pixel Array Data & Parameters.
	pixelArray SbmPixelArray
}
