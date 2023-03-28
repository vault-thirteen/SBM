package sbm

import (
	"errors"
	"fmt"

	"github.com/vault-thirteen/auxie/bit"
	rdr "github.com/vault-thirteen/auxie/reader"
)

// Errors.
const (
	ErrBottomHeaderMismatch = "bottom header mismatch"
	ErrBadSeparator         = "bad separator"
	ErrAreaMismatch         = "area mismatch"
)

func (sbm *Sbm) readTopHeaders(lineReader *rdr.Reader) (err error) {

	// 1. Format.
	var line []byte
	line, err = lineReader.ReadLineEndingWithCRLF()
	if err != nil {
		return err
	}
	err = sbm.parseHeaderFormat(line)
	if err != nil {
		return err
	}

	// 2. Version.
	line, err = lineReader.ReadLineEndingWithCRLF()
	if err != nil {
		return err
	}
	var headerFormat HeaderDataVersion
	headerFormat, err = parseHeaderVersion(line)
	if err != nil {
		return err
	}
	err = validateFormat(headerFormat)
	if err != nil {
		return err
	}
	sbm.format.version = headerFormat.version

	// 3. Width.
	line, err = lineReader.ReadLineEndingWithCRLF()
	if err != nil {
		return err
	}
	var headerSize HeaderDataSize
	headerSize, err = parseHeaderWidth(line)
	if err != nil {
		return err
	}
	sbm.pixelArray.metaData.width = headerSize.sizeFixed
	sbm.pixelArray.metaData.header.width.topLeft = headerSize.sizeRandomLeft
	sbm.pixelArray.metaData.header.width.topRight = headerSize.sizeRandomRight

	// 4. Height.
	line, err = lineReader.ReadLineEndingWithCRLF()
	if err != nil {
		return err
	}
	headerSize, err = parseHeaderHeight(line)
	if err != nil {
		return err
	}
	sbm.pixelArray.metaData.height = headerSize.sizeFixed
	sbm.pixelArray.metaData.header.height.topLeft = headerSize.sizeRandomLeft
	sbm.pixelArray.metaData.header.height.topRight = headerSize.sizeRandomRight

	// 5. Area...

	// 5.1. Get the Data.
	line, err = lineReader.ReadLineEndingWithCRLF()
	if err != nil {
		return err
	}
	headerSize, err = parseHeaderArea(line)
	if err != nil {
		return err
	}

	// 5.2. Verify the Data.
	if sbm.pixelArray.metaData.width*sbm.pixelArray.metaData.height != headerSize.sizeFixed {
		return errors.New(ErrAreaMismatch)
	}

	// 5.3. Save the Data.
	sbm.pixelArray.metaData.area = headerSize.sizeFixed
	sbm.pixelArray.metaData.header.area.topLeft = headerSize.sizeRandomLeft
	sbm.pixelArray.metaData.header.area.topRight = headerSize.sizeRandomRight

	return nil
}

func (sbm *Sbm) readArrayData(lineReader *rdr.Reader) (err error) {

	// Exception Handler.
	defer func() {
		xErr := recover()
		if xErr != nil {
			err = fmt.Errorf("xErr:%v.", xErr)
		}
	}()

	// Get the Size of Array to know how many Bytes to read.
	var lastByteIsPartial bool
	var readBytesSize uint
	if sbm.pixelArray.metaData.area%bit.BitsPerByte == 0 {
		readBytesSize = sbm.pixelArray.metaData.area / bit.BitsPerByte
	} else {
		lastByteIsPartial = true
		readBytesSize = (sbm.pixelArray.metaData.area / bit.BitsPerByte) + 1
	}

	// Read the Array.
	var bytesArray []byte
	bytesArray, err = lineReader.ReadBytes(int(readBytesSize))
	if err != nil {
		return err
	}

	// Read the Separator.
	err = readSeparator(lineReader)
	if err != nil {
		return err
	}

	// Save the Array...

	// 1. Bits.
	var bitsArray []bit.Bit
	bitsArray = bit.ConvertBytesToBits(bytesArray)
	if lastByteIsPartial {
		bitsArray = bitsArray[:sbm.pixelArray.metaData.area]
	}
	sbm.pixelArray.data.bits = bitsArray

	// 2. Bytes.
	sbm.pixelArray.data.bytes, _ = bit.ConvertBitsToBytes(bitsArray)

	return nil
}

func (sbm *Sbm) readBottomHeaders(lineReader *rdr.Reader) (err error) {

	// 1. Width ...

	// 1.1. Get the Data.
	var line []byte
	line, err = lineReader.ReadLineEndingWithCRLF()
	if err != nil {
		return err
	}
	var headerSize HeaderDataSize
	headerSize, err = parseHeaderWidth(line)
	if err != nil {
		return err
	}

	// 1.2. Verify the Data.
	if headerSize.sizeFixed != sbm.pixelArray.metaData.width {
		return errors.New(ErrBottomHeaderMismatch)
	}

	// 1.3. Save the Data.
	sbm.pixelArray.metaData.header.width.bottomLeft = headerSize.sizeRandomLeft
	sbm.pixelArray.metaData.header.width.bottomRight = headerSize.sizeRandomRight

	// 2. Height ...

	// 2.1. Get the Data.
	line, err = lineReader.ReadLineEndingWithCRLF()
	if err != nil {
		return err
	}
	headerSize, err = parseHeaderHeight(line)
	if err != nil {
		return err
	}

	// 2.2. Verify the Data.
	if headerSize.sizeFixed != sbm.pixelArray.metaData.height {
		return errors.New(ErrBottomHeaderMismatch)
	}

	// 2.3. Save the Data.
	sbm.pixelArray.metaData.header.height.bottomLeft = headerSize.sizeRandomLeft
	sbm.pixelArray.metaData.header.height.bottomRight = headerSize.sizeRandomRight

	// 3. Area ...

	// 3.1. Get the Data.
	line, err = lineReader.ReadLineEndingWithCRLF()
	if err != nil {
		return err
	}
	headerSize, err = parseHeaderArea(line)
	if err != nil {
		return err
	}

	// 3.2. Verify the Data.
	if headerSize.sizeFixed != sbm.pixelArray.metaData.area {
		return errors.New(ErrBottomHeaderMismatch)
	}
	if sbm.pixelArray.metaData.width*sbm.pixelArray.metaData.height != headerSize.sizeFixed {
		return errors.New(ErrAreaMismatch)
	}

	// 3.3. Save the Data.
	sbm.pixelArray.metaData.header.area.bottomLeft = headerSize.sizeRandomLeft
	sbm.pixelArray.metaData.header.area.bottomRight = headerSize.sizeRandomRight

	return nil
}

func readSeparator(lineReader *rdr.Reader) (err error) {
	var ba []byte
	ba, err = lineReader.ReadBytes(2)
	if err != nil {
		return err
	}

	if (len(ba) != 2) || (ba[0] != CR) || (ba[1] != LF) {
		return errors.New(ErrBadSeparator)
	}

	return nil
}
