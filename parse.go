package sbm

import (
	"bytes"
	"math"
	"strconv"
	"strings"

	"github.com/vault-thirteen/auxie/errors"
)

// Errors.
const (
	Err_Format       = "format is unrecognized"
	Err_HeaderSyntax = "header syntax error"
	Err_Overflow     = "overflow"
	Err_Integrity    = "integrity failure"

	Errf_HeaderUnexpected = "unexpected header: '%v'"
)

// parseHeaderFormat parses the format header.
func (sbm *Sbm) parseHeaderFormat(rawHeader []byte) (err error) {
	if !bytes.Equal(rawHeader, []byte(Header_FormatName)) {
		return errors.ErrorFromText(Err_Format)
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
		return headerData, errors.ErrorFromText(Err_HeaderSyntax)
	}

	// Check the header name.
	if headerParts[0] != HeaderPrefix_Version {
		return headerData, errors.ErrorFromFormattedText(Errf_HeaderUnexpected, headerParts[0])
	}

	// Parse the version number.
	var versionNumberTmp uint64
	versionNumberTmp, err = strconv.ParseUint(headerParts[1], 10, 64)
	if err != nil {
		return headerData, err
	}
	if versionNumberTmp > math.MaxUint8 {
		return headerData, errors.ErrorFromText(Err_Overflow)
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
		return headerData, errors.ErrorFromText(Err_HeaderSyntax)
	}

	// Check the header name.
	if headerParts[0] != headerNameExpected {
		return headerData, errors.ErrorFromFormattedText(Errf_HeaderUnexpected, headerParts[0])
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
		return headerData, errors.ErrorFromText(Err_HeaderSyntax)
	}
	headerPartSizeLeft = strings.TrimLeft(headerPartSizeLeft, HeaderPartsBracketLeft)
	sizeTmp, err = strconv.ParseUint(headerPartSizeLeft, 10, 64)
	if err != nil {
		return headerData, err
	}
	sizeRandomLeft := uint(sizeTmp)

	// Plus sign.
	if headerParts[3] != HeaderPartsPlus {
		return headerData, errors.ErrorFromText(Err_HeaderSyntax)
	}

	// Parse the random right size.
	headerPartSizeRight := headerParts[4]
	if !strings.HasSuffix(headerPartSizeRight, HeaderPartsBracketRight) {
		return headerData, errors.ErrorFromText(Err_HeaderSyntax)
	}
	headerPartSizeRight = strings.TrimRight(headerPartSizeRight, HeaderPartsBracketRight)
	sizeTmp, err = strconv.ParseUint(headerPartSizeRight, 10, 64)
	if err != nil {
		return headerData, err
	}
	sizeRandomRight := uint(sizeTmp)

	// Verify the integrity of the size.
	if sizeFixed != (sizeRandomLeft + sizeRandomRight) {
		return headerData, errors.ErrorFromText(Err_Integrity)
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
