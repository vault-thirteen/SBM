// new.go.

package sbm

import (
	"errors"
	"io"

	"github.com/vault-thirteen/bit"
	rdr "github.com/vault-thirteen/reader"
)

// NewFromBitsArray creates a new SBM from an Array of Bits.
// Performs the Fool Checks.
func NewFromBitsArray(
	arrayBits []bit.Bit, // List of all Bits of a 2D-Array in a 1D-Array.
	arrayWidth uint, // Width of a 2D-Array.
	arrayHeight uint, // Height of a 2D-Array.
) (result *Sbm, err error) {

	var arrayArea uint

	// Preparation.
	arrayArea = arrayWidth * arrayHeight

	// Checks.
	if (arrayWidth == 0) ||
		(arrayHeight == 0) ||
		(uint(len(arrayBits)) != arrayArea) {
		err = errors.New(ErrDimension)
		return
	}

	return newFromBitsArray(arrayBits, arrayWidth, arrayHeight)
}

// NewFromBytesArray creates a new SBM from an Array of Bytes.
// Performs the Fool Checks.
func NewFromBytesArray(
	arrayBytes []byte, // List of all Bytes in a 1D-Array.
	arrayWidth uint, // Width of a 2D-Array.
	arrayHeight uint, // Height of a 2D-Array.
) (result *Sbm, err error) {

	var arrayAreaReal uint
	var arrayBitsCountMax uint
	var arrayBitsCountMin uint

	// Preparation.
	arrayAreaReal = arrayWidth * arrayHeight
	arrayBitsCountMax = uint(len(arrayBytes) * bit.BitsPerByte)
	arrayBitsCountMin = arrayBitsCountMax - (bit.BitsPerByte - 1)

	// Checks.
	if (arrayWidth == 0) ||
		(arrayHeight == 0) ||
		(arrayAreaReal > arrayBitsCountMax) ||
		(arrayAreaReal < arrayBitsCountMin) {
		err = errors.New(ErrDimension)
		return
	}

	return newFromBytesArray(arrayBytes, arrayWidth, arrayHeight)
}

// newFromBitsArray creates a new SBM from an Array of Bits.
// Does not perform the Fool Checks.
func newFromBitsArray(
	arrayBits []bit.Bit, // List of all Bits of a 2D-Array in a 1D-Array.
	arrayWidth uint, // Width of a 2D-Array.
	arrayHeight uint, // Height of a 2D-Array.
) (result *Sbm, err error) {

	var arrayArea uint
	var arrayBytes []byte

	// Preparation.
	arrayArea = arrayWidth * arrayHeight
	arrayBytes, _ = bit.ConvertBitsToBytes(arrayBits)

	return newFromBitsAndBytesArrays(
		arrayBits,
		arrayBytes,
		arrayWidth,
		arrayHeight,
		arrayArea,
	)
}

// newFromBytesArray creates a new SBM from an Array of Bits.
// Does not perform the Fool Checks.
func newFromBytesArray(
	arrayBytes []byte, // List of all Bytes in a 1D-Array.
	arrayWidth uint, // Width of a 2D-Array.
	arrayHeight uint, // Height of a 2D-Array.
) (result *Sbm, err error) {

	var arrayArea uint
	var arrayBits []bit.Bit

	// Preparation.
	arrayArea = arrayWidth * arrayHeight

	// Convert all Bytes to Bits and remove the redundant Bits from the End.
	arrayBits = bit.ConvertBytesToBits(arrayBytes)
	if arrayArea != uint(len(arrayBits)) {
		arrayBits = arrayBits[0:arrayArea]
	}

	return newFromBitsAndBytesArrays(
		arrayBits,
		arrayBytes,
		arrayWidth,
		arrayHeight,
		arrayArea,
	)
}

// newFromBitsAndBytesArrays creates a new SBM from the Arrays of Bits and
// Bytes. Does not perform the Fool Checks.
func newFromBitsAndBytesArrays(
	arrayBits []bit.Bit, // List of all Bits of a 2D-Array in a 1D-Array.
	arrayBytes []byte, // List of all Bytes in a 1D-Array.
	arrayWidth uint, // Width of a 2D-Array.
	arrayHeight uint, // Height of a 2D-Array.
	arrayArea uint,
) (result *Sbm, err error) {

	// Create a new Object.
	result = new(Sbm)

	// Fill the Array Data and Dimensions.
	result.format.version = SbmFormatVersion1
	result.pixelArray = SbmPixelArray{
		data: SbmPixelArrayData{
			bits:  arrayBits,
			bytes: arrayBytes,
		},
		metaData: SbmPixelArrayMetaData{
			width:  arrayWidth,
			height: arrayHeight,
			area:   arrayArea,
			header: SbmPixelArrayMetaDataHeader{},
		},
	}

	// Fill the random Header Fields...

	// 1. Width.
	result.pixelArray.metaData.header.width.topLeft,
		result.pixelArray.metaData.header.width.topRight,
		err = createRandomValuePair(arrayWidth)
	if err != nil {
		return
	}
	result.pixelArray.metaData.header.width.bottomLeft,
		result.pixelArray.metaData.header.width.bottomRight,
		err = createRandomValuePair(arrayWidth)
	if err != nil {
		return
	}

	// 2. Height.
	result.pixelArray.metaData.header.height.topLeft,
		result.pixelArray.metaData.header.height.topRight,
		err = createRandomValuePair(arrayHeight)
	if err != nil {
		return
	}
	result.pixelArray.metaData.header.height.bottomLeft,
		result.pixelArray.metaData.header.height.bottomRight,
		err = createRandomValuePair(arrayHeight)
	if err != nil {
		return
	}

	// 3. Area.
	result.pixelArray.metaData.header.area.topLeft,
		result.pixelArray.metaData.header.area.topRight,
		err = createRandomValuePair(arrayArea)
	if err != nil {
		return
	}
	result.pixelArray.metaData.header.area.bottomLeft,
		result.pixelArray.metaData.header.area.bottomRight,
		err = createRandomValuePair(arrayArea)
	if err != nil {
		return
	}

	return
}

// NewFromStream reads an SBM Object from the Stream.
func NewFromStream(
	reader io.Reader,
) (*Sbm, error) {

	var err error
	var lineReader *rdr.Reader
	var sbm *Sbm

	sbm = new(Sbm)
	lineReader = rdr.NewReader(reader)

	// Read the top Headers.
	err = sbm.readTopHeaders(lineReader)
	if err != nil {
		return nil, err
	}

	// Read the binary Array of Bits with the 'NewFromBitsArray Line' at the End.
	err = sbm.readArrayData(lineReader)
	if err != nil {
		return nil, err
	}

	// Read the bottom Headers.
	err = sbm.readBottomHeaders(lineReader)
	if err != nil {
		return nil, err
	}

	return sbm, nil
}
