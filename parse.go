package sbm

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Errors.
const (
	ErrFormat            = "format is unrecognized"
	ErrHeaderSyntax      = "header syntax error"
	ErrfHeaderUnexpected = "unexpected header: '%v'"
	ErrOverflow          = "overflow"
	ErrIntegrity         = "integrity failure"
)

// parseHeaderFormat parses the format header.
func (sbm *Sbm) parseHeaderFormat(rawHeader []byte) (err error) {
	if !bytes.Equal(rawHeader, []byte(Header_FormatName)) {
		return errors.New(ErrFormat)
	}

	return nil
}

// parseHeaderVersion parses the version header.
func parseHeaderVersion(rawHeader []byte) (headerData HeaderDataVersion, err error) {

	// Trim.
	var rawHeaderTrimmed []byte
	rawHeaderTrimmed, err = trimHeader(rawHeader)
	if err != nil {
		return headerData, err
	}

	// Split.
	headerPartsCountExpected := 2
	headerParts := strings.Split(string(rawHeaderTrimmed), HeaderPartsSeparator)
	if len(headerParts) != headerPartsCountExpected {
		return headerData, errors.New(ErrHeaderSyntax)
	}

	// Check the header name.
	if headerParts[0] != HeaderPrefix_Version {
		return headerData, fmt.Errorf(ErrfHeaderUnexpected, headerParts[0])
	}

	// Parse the version number.
	var versionNumberTmp uint64
	versionNumberTmp, err = strconv.ParseUint(headerParts[1], 10, 64)
	if err != nil {
		return headerData, err
	}
	if versionNumberTmp > math.MaxUint8 {
		return headerData, errors.New(ErrOverflow)
	}

	// Save the version number.
	headerData.version = byte(versionNumberTmp)

	return headerData, nil
}

// parseHeaderSize parses the size header.
func parseHeaderSize(rawHeader []byte, headerNameExpected string) (headerData HeaderDataSize, err error) {

	// Trim.
	var rawHeaderTrimmed []byte
	rawHeaderTrimmed, err = trimHeader(rawHeader)
	if err != nil {
		return headerData, err
	}

	// Split.
	headerPartsCountExpected := 5
	headerParts := strings.Split(string(rawHeaderTrimmed), HeaderPartsSeparator)
	if len(headerParts) != headerPartsCountExpected {
		return headerData, errors.New(ErrHeaderSyntax)
	}

	// Check the header name.
	if headerParts[0] != headerNameExpected {
		return headerData, fmt.Errorf(ErrfHeaderUnexpected, headerParts[0])
	}

	// Parse the fixed size.
	var sizeTmp uint64
	sizeTmp, err = strconv.ParseUint(headerParts[1], 10, 64)
	if err != nil {
		return headerData, err
	}
	sizeFixed := uint(sizeTmp)

	// Parse the random left size.
	headerPartSizeLeft := headerParts[2]
	if !strings.HasPrefix(headerPartSizeLeft, HeaderPartsBracketLeft) {
		return headerData, errors.New(ErrHeaderSyntax)
	}
	headerPartSizeLeft = strings.TrimLeft(headerPartSizeLeft, HeaderPartsBracketLeft)
	sizeTmp, err = strconv.ParseUint(headerPartSizeLeft, 10, 64)
	if err != nil {
		return headerData, err
	}
	sizeRandomLeft := uint(sizeTmp)

	// Plus sign.
	if headerParts[3] != HeaderPartsPlus {
		return headerData, errors.New(ErrHeaderSyntax)
	}

	// Parse the random right size.
	headerPartSizeRight := headerParts[4]
	if !strings.HasSuffix(headerPartSizeRight, HeaderPartsBracketRight) {
		return headerData, errors.New(ErrHeaderSyntax)
	}
	headerPartSizeRight = strings.TrimRight(headerPartSizeRight, HeaderPartsBracketRight)
	sizeTmp, err = strconv.ParseUint(headerPartSizeRight, 10, 64)
	if err != nil {
		return headerData, err
	}
	sizeRandomRight := uint(sizeTmp)

	// Verify the integrity of the size.
	if sizeFixed != (sizeRandomLeft + sizeRandomRight) {
		return headerData, errors.New(ErrIntegrity)
	}

	// Save the data.
	headerData.sizeFixed = sizeFixed
	headerData.sizeRandomLeft = sizeRandomLeft
	headerData.sizeRandomRight = sizeRandomRight

	return headerData, nil
}

// parseHeaderWidth parses the width header.
func parseHeaderWidth(rawHeader []byte) (headerData HeaderDataSize, err error) {
	return parseHeaderSize(rawHeader, HeaderPrefix_Width)
}

// parseHeaderHeight parses the height header.
func parseHeaderHeight(rawHeader []byte) (headerData HeaderDataSize, err error) {
	return parseHeaderSize(rawHeader, HeaderPrefix_Height)
}

// parseHeaderArea parses the area header.
func parseHeaderArea(rawHeader []byte) (headerData HeaderDataSize, err error) {
	return parseHeaderSize(rawHeader, HeaderPrefix_Area)
}
