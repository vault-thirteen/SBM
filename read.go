// read.go.

package sbm

import (
	"errors"
	"fmt"

	"github.com/vault-thirteen/bit"
	rdr "github.com/vault-thirteen/reader"
)

// Errors.
const (
	ErrBottomHeaderMismatch = "Bottom Header Mismatch"
	ErrBadSeparator         = "Bad Separator"
	ErrAreaMismatch         = "Area Mismatch"
)

func (sbm *Sbm) readTopHeaders(
	lineReader *rdr.Reader,
) (err error) {

	var headerFormat HeaderDataVersion
	var headerSize HeaderDataSize
	var line []byte

	// 1. Format.
	line, err = lineReader.ReadLineEndingWithCRLF()
	if err != nil {
		return
	}
	err = sbm.parseHeaderFormat(line)
	if err != nil {
		return
	}

	// 2. Version.
	line, err = lineReader.ReadLineEndingWithCRLF()
	if err != nil {
		return
	}
	headerFormat, err = parseHeaderVersion(line)
	if err != nil {
		return
	}
	err = validateFormat(headerFormat)
	if err != nil {
		return
	}
	sbm.format.version = headerFormat.version

	// 3. Width.
	line, err = lineReader.ReadLineEndingWithCRLF()
	if err != nil {
		return
	}
	headerSize, err = parseHeaderWidth(line)
	if err != nil {
		return
	}
	sbm.pixelArray.metaData.width = headerSize.sizeFixed
	sbm.pixelArray.metaData.header.width.topLeft = headerSize.sizeRandomLeft
	sbm.pixelArray.metaData.header.width.topRight = headerSize.sizeRandomRight

	// 4. Height.
	line, err = lineReader.ReadLineEndingWithCRLF()
	if err != nil {
		return
	}
	headerSize, err = parseHeaderHeight(line)
	if err != nil {
		return
	}
	sbm.pixelArray.metaData.height = headerSize.sizeFixed
	sbm.pixelArray.metaData.header.height.topLeft = headerSize.sizeRandomLeft
	sbm.pixelArray.metaData.header.height.topRight = headerSize.sizeRandomRight

	// 5. Area...

	// 5.1. Get the Data.
	line, err = lineReader.ReadLineEndingWithCRLF()
	if err != nil {
		return
	}
	headerSize, err = parseHeaderArea(line)
	if err != nil {
		return
	}

	// 5.2. Verify the Data.
	if sbm.pixelArray.metaData.width*sbm.pixelArray.metaData.height != headerSize.sizeFixed {
		err = errors.New(ErrAreaMismatch)
		return
	}

	// 5.3. Save the Data.
	sbm.pixelArray.metaData.area = headerSize.sizeFixed
	sbm.pixelArray.metaData.header.area.topLeft = headerSize.sizeRandomLeft
	sbm.pixelArray.metaData.header.area.topRight = headerSize.sizeRandomRight

	return
}

func (sbm *Sbm) readArrayData(
	lineReader *rdr.Reader,
) (err error) {

	var bytesArray []byte
	var bitsArray []bit.Bit
	var lastByteIsPartial bool
	var readBytesSize uint
	var xErr interface{}

	// Exception Handler.
	defer func() {
		xErr = recover()
		if xErr != nil {
			err = fmt.Errorf("xErr:%v.", xErr)
		}
	}()

	// Get the Size of Array to know how many Bytes to read.
	if sbm.pixelArray.metaData.area%bit.BitsPerByte == 0 {
		readBytesSize = sbm.pixelArray.metaData.area / bit.BitsPerByte
	} else {
		lastByteIsPartial = true
		readBytesSize = (sbm.pixelArray.metaData.area / bit.BitsPerByte) + 1
	}

	// Read the Array.
	bytesArray, err = lineReader.ReadBytes(int(readBytesSize))
	if err != nil {
		return
	}

	// Read the Separator.
	err = readSeparator(lineReader)
	if err != nil {
		return
	}

	// Save the Array...

	// 1. Bits.
	bitsArray = bit.ConvertBytesToBits(bytesArray)
	if lastByteIsPartial {
		bitsArray = bitsArray[:sbm.pixelArray.metaData.area]
	}
	sbm.pixelArray.data.bits = bitsArray

	// 2. Bytes.
	sbm.pixelArray.data.bytes, _ = bit.ConvertBitsToBytes(bitsArray)

	return
}

func (sbm *Sbm) readBottomHeaders(
	lineReader *rdr.Reader,
) (err error) {

	var headerSize HeaderDataSize
	var line []byte

	// 1. Width...

	// 1.1. Get the Data.
	line, err = lineReader.ReadLineEndingWithCRLF()
	if err != nil {
		return
	}
	headerSize, err = parseHeaderWidth(line)
	if err != nil {
		return
	}

	// 1.2. Verify the Data.
	if headerSize.sizeFixed != sbm.pixelArray.metaData.width {
		err = errors.New(ErrBottomHeaderMismatch)
		return
	}

	// 1.3. Save the Data.
	sbm.pixelArray.metaData.header.width.bottomLeft = headerSize.sizeRandomLeft
	sbm.pixelArray.metaData.header.width.bottomRight = headerSize.sizeRandomRight

	// 2. Height...

	// 2.1. Get the Data.
	line, err = lineReader.ReadLineEndingWithCRLF()
	if err != nil {
		return
	}
	headerSize, err = parseHeaderHeight(line)
	if err != nil {
		return
	}

	// 2.2. Verify the Data.
	if headerSize.sizeFixed != sbm.pixelArray.metaData.height {
		err = errors.New(ErrBottomHeaderMismatch)
		return
	}

	// 2.3. Save the Data.
	sbm.pixelArray.metaData.header.height.bottomLeft = headerSize.sizeRandomLeft
	sbm.pixelArray.metaData.header.height.bottomRight = headerSize.sizeRandomRight

	// 3. Area...

	// 3.1. Get the Data.
	line, err = lineReader.ReadLineEndingWithCRLF()
	if err != nil {
		return
	}
	headerSize, err = parseHeaderArea(line)
	if err != nil {
		return
	}

	// 3.2. Verify the Data.
	if headerSize.sizeFixed != sbm.pixelArray.metaData.area {
		err = errors.New(ErrBottomHeaderMismatch)
		return
	}
	if sbm.pixelArray.metaData.width*sbm.pixelArray.metaData.height != headerSize.sizeFixed {
		err = errors.New(ErrAreaMismatch)
		return
	}

	// 3.3. Save the Data.
	sbm.pixelArray.metaData.header.area.bottomLeft = headerSize.sizeRandomLeft
	sbm.pixelArray.metaData.header.area.bottomRight = headerSize.sizeRandomRight

	return
}

func readSeparator(
	lineReader *rdr.Reader,
) (err error) {

	var ba []byte

	ba, err = lineReader.ReadBytes(2)
	if err != nil {
		return
	}

	if (len(ba) != 2) ||
		(ba[0] != '\r') ||
		(ba[1] != '\n') {
		err = errors.New(ErrBadSeparator)
		return
	}

	return
}
