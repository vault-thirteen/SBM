package sbm

import (
	"errors"
	"io"

	"github.com/vault-thirteen/auxie/bit"
	rdr "github.com/vault-thirteen/auxie/reader"
)

// NewFromBitsArray creates a new SBM from an array of bits.
// Performs the fool checks.
func NewFromBitsArray(
	arrayBits []bit.Bit, // List of all bits of a 2D-array in a 1D-array.
	arrayWidth uint, // Width of a 2D-array.
	arrayHeight uint, // Height of a 2D-array.
) (sbm *Sbm, err error) {
	var arrayArea = arrayWidth * arrayHeight

	// Checks.
	if (arrayWidth == 0) ||
		(arrayHeight == 0) ||
		(uint(len(arrayBits)) != arrayArea) {
		return nil, errors.New(ErrDimension)
	}

	return newFromBitsArray(arrayBits, arrayWidth, arrayHeight)
}

// NewFromBytesArray creates a new SBM from an array of bytes.
// Performs the fool checks.
func NewFromBytesArray(
	arrayBytes []byte, // List of all bytes in a 1D-array.
	arrayWidth uint, // Width of a 2D-array.
	arrayHeight uint, // Height of a 2D-array.
) (sbm *Sbm, err error) {
	arrayAreaReal := arrayWidth * arrayHeight
	arrayBitsCountMax := uint(len(arrayBytes) * bit.BitsPerByte)
	arrayBitsCountMin := arrayBitsCountMax - (bit.BitsPerByte - 1)

	// Checks.
	if (arrayWidth == 0) ||
		(arrayHeight == 0) ||
		(arrayAreaReal > arrayBitsCountMax) ||
		(arrayAreaReal < arrayBitsCountMin) {
		return nil, errors.New(ErrDimension)
	}

	return newFromBytesArray(arrayBytes, arrayWidth, arrayHeight)
}

// newFromBitsArray creates a new SBM from an array of bits.
// Does not perform the fool checks.
func newFromBitsArray(
	arrayBits []bit.Bit, // List of all bits of a 2D-array in a 1D-array.
	arrayWidth uint, // Width of a 2D-array.
	arrayHeight uint, // Height of a 2D-array.
) (sbm *Sbm, err error) {
	arrayArea := arrayWidth * arrayHeight
	arrayBytes, _ := bit.ConvertBitsToBytes(arrayBits)

	return newFromBitsAndBytesArrays(arrayBits, arrayBytes, arrayWidth, arrayHeight, arrayArea)
}

// newFromBytesArray creates a new SBM from an array of bits.
// Does not perform the fool checks.
func newFromBytesArray(
	arrayBytes []byte, // List of all bytes in a 1D-array.
	arrayWidth uint, // Width of a 2D-array.
	arrayHeight uint, // Height of a 2D-array.
) (sbm *Sbm, err error) {
	arrayArea := arrayWidth * arrayHeight

	// Convert all bytes to bits and remove the redundant bits from the end.
	arrayBits := bit.ConvertBytesToBits(arrayBytes)
	if arrayArea != uint(len(arrayBits)) {
		arrayBits = arrayBits[0:arrayArea]
	}

	return newFromBitsAndBytesArrays(arrayBits, arrayBytes, arrayWidth, arrayHeight, arrayArea)
}

// newFromBitsAndBytesArrays creates a new SBM from the arrays of bits and
// bytes. Does not perform the fool checks.
func newFromBitsAndBytesArrays(
	arrayBits []bit.Bit, // List of all bits of a 2D-array in a 1D-array.
	arrayBytes []byte, // List of all bytes in a 1D-array.
	arrayWidth uint, // Width of a 2D-array.
	arrayHeight uint, // Height of a 2D-array.
	arrayArea uint,
) (sbm *Sbm, err error) {

	// Create a new Object.
	sbm = new(Sbm)

	// Fill the array data and dimensions.
	sbm.format.version = SbmFormatVersion1
	sbm.pixelArray = SbmPixelArray{
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

	// Fill the random header fields ...

	// 1. Width.
	sbm.pixelArray.metaData.header.width.topLeft,
		sbm.pixelArray.metaData.header.width.topRight,
		err = createRandomValuePair(arrayWidth)
	if err != nil {
		return nil, err
	}
	sbm.pixelArray.metaData.header.width.bottomLeft,
		sbm.pixelArray.metaData.header.width.bottomRight,
		err = createRandomValuePair(arrayWidth)
	if err != nil {
		return nil, err
	}

	// 2. Height.
	sbm.pixelArray.metaData.header.height.topLeft,
		sbm.pixelArray.metaData.header.height.topRight,
		err = createRandomValuePair(arrayHeight)
	if err != nil {
		return nil, err
	}
	sbm.pixelArray.metaData.header.height.bottomLeft,
		sbm.pixelArray.metaData.header.height.bottomRight,
		err = createRandomValuePair(arrayHeight)
	if err != nil {
		return nil, err
	}

	// 3. Area.
	sbm.pixelArray.metaData.header.area.topLeft,
		sbm.pixelArray.metaData.header.area.topRight,
		err = createRandomValuePair(arrayArea)
	if err != nil {
		return nil, err
	}
	sbm.pixelArray.metaData.header.area.bottomLeft,
		sbm.pixelArray.metaData.header.area.bottomRight,
		err = createRandomValuePair(arrayArea)
	if err != nil {
		return nil, err
	}

	return sbm, nil
}

// NewFromStream reads an SBM object from the stream.
func NewFromStream(reader io.Reader) (sbm *Sbm, err error) {
	sbm = new(Sbm)
	lineReader := rdr.New(reader)

	// Read the top headers.
	err = sbm.readTopHeaders(lineReader)
	if err != nil {
		return nil, err
	}

	// Read the binary array of bits with the 'NewFromBitsArray Line' at the end.
	err = sbm.readArrayData(lineReader)
	if err != nil {
		return nil, err
	}

	// Read the bottom headers.
	err = sbm.readBottomHeaders(lineReader)
	if err != nil {
		return nil, err
	}

	return sbm, nil
}
