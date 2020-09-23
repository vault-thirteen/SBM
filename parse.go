// parse.go.

////////////////////////////////////////////////////////////////////////////////
//
// Copyright © 2019..2020 by Vault Thirteen.
//
// All rights reserved. No part of this publication may be reproduced,
// distributed, or transmitted in any form or by any means, including
// photocopying, recording, or other electronic or mechanical methods,
// without the prior written permission of the publisher, except in the case
// of brief quotations embodied in critical reviews and certain other
// noncommercial uses permitted by copyright law. For permission requests,
// write to the publisher, addressed “Copyright Protected Material” at the
// address below.
//
////////////////////////////////////////////////////////////////////////////////
//
// Web Site Address:	https://github.com/vault-thirteen.
//
////////////////////////////////////////////////////////////////////////////////

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
	ErrFormat            = "Format is unrecognized"
	ErrHeaderSyntax      = "Header Syntax Error"
	ErrfHeaderUnexpected = "Unexpected Header: '%v'"
	ErrOverflow          = "Overflow"
	ErrIntegrity         = "Integrity Failure"
)

// Parses the Format Header.
func (sbm Sbm) parseHeaderFormat(
	rawHeader []byte,
) (err error) {

	if !bytes.Equal(rawHeader, []byte(Header_FormatName)) {
		err = errors.New(ErrFormat)
		return
	}

	return
}

// Parses the 'Version' Header.
func parseHeaderVersion(
	rawHeader []byte,
) (headerData HeaderDataVersion, err error) {

	var headerParts []string
	var headerPartsCountExpected int
	var rawHeaderTrimmed []byte
	var versionNumber byte
	var versionNumberTmp uint64

	// Trim.
	rawHeaderTrimmed, err = removeCRLF(rawHeader)
	if err != nil {
		return
	}

	// Split.
	headerPartsCountExpected = 2
	headerParts = strings.Split(string(rawHeaderTrimmed), HeaderPartsSeparator)
	if len(headerParts) != headerPartsCountExpected {
		err = errors.New(ErrHeaderSyntax)
		return
	}

	// Check the Header Name.
	if headerParts[0] != HeaderPrefix_Version {
		err = fmt.Errorf(ErrfHeaderUnexpected, headerParts[0])
		return
	}

	// Parse the Version Number.
	versionNumberTmp, err = strconv.ParseUint(headerParts[1], 10, 64)
	if err != nil {
		return
	}
	if versionNumberTmp > math.MaxUint8 {
		err = errors.New(ErrOverflow)
		return
	}
	versionNumber = byte(versionNumberTmp)

	// Save the Version Number.
	headerData.version = versionNumber

	return
}

// Parses the Size Header.
func parseHeaderSize(
	rawHeader []byte,
	headerNameExpected string,
) (headerData HeaderDataSize, err error) {

	var headerParts []string
	var headerPartSizeLeft string
	var headerPartSizeRight string
	var headerPartsCountExpected int
	var rawHeaderTrimmed []byte
	var sizeFixed uint
	var sizeRandomLeft uint
	var sizeRandomRight uint
	var sizeTmp uint64

	// Trim.
	rawHeaderTrimmed, err = removeCRLF(rawHeader)
	if err != nil {
		return
	}

	// Split.
	headerPartsCountExpected = 5
	headerParts = strings.Split(string(rawHeaderTrimmed), HeaderPartsSeparator)
	if len(headerParts) != headerPartsCountExpected {
		err = errors.New(ErrHeaderSyntax)
		return
	}

	// Check the Header Name.
	if headerParts[0] != headerNameExpected {
		err = fmt.Errorf(ErrfHeaderUnexpected, headerParts[0])
		return
	}

	// Parse the fixed Size.
	sizeTmp, err = strconv.ParseUint(headerParts[1], 10, 64)
	if err != nil {
		return
	}
	sizeFixed = uint(sizeTmp)

	// Parse the random left Size.
	headerPartSizeLeft = headerParts[2]
	if !strings.HasPrefix(headerPartSizeLeft, HeaderPartsBracketLeft) {
		err = errors.New(ErrHeaderSyntax)
		return
	}
	headerPartSizeLeft = strings.TrimLeft(headerPartSizeLeft, HeaderPartsBracketLeft)
	sizeTmp, err = strconv.ParseUint(headerPartSizeLeft, 10, 64)
	if err != nil {
		return
	}
	sizeRandomLeft = uint(sizeTmp)

	// Plus Sign.
	if headerParts[3] != HeaderPartsPlus {
		err = errors.New(ErrHeaderSyntax)
		return
	}

	// Parse the random right Size.
	headerPartSizeRight = headerParts[4]
	if !strings.HasSuffix(headerPartSizeRight, HeaderPartsBracketRight) {
		err = errors.New(ErrHeaderSyntax)
		return
	}
	headerPartSizeRight = strings.TrimRight(headerPartSizeRight, HeaderPartsBracketRight)
	sizeTmp, err = strconv.ParseUint(headerPartSizeRight, 10, 64)
	if err != nil {
		return
	}
	sizeRandomRight = uint(sizeTmp)

	// Verify the Integrity of the Size.
	if sizeFixed != (sizeRandomLeft + sizeRandomRight) {
		err = errors.New(ErrIntegrity)
		return
	}

	// Save the Data.
	headerData.sizeFixed = sizeFixed
	headerData.sizeRandomLeft = sizeRandomLeft
	headerData.sizeRandomRight = sizeRandomRight

	return
}

// Parses the 'Width' Header.
func parseHeaderWidth(
	rawHeader []byte,
) (headerData HeaderDataSize, err error) {
	return parseHeaderSize(rawHeader, HeaderPrefix_Width)
}

// Parses the 'Height' Header.
func parseHeaderHeight(
	rawHeader []byte,
) (headerData HeaderDataSize, err error) {
	return parseHeaderSize(rawHeader, HeaderPrefix_Height)
}

// Parses the 'Area' Header.
func parseHeaderArea(
	rawHeader []byte,
) (headerData HeaderDataSize, err error) {
	return parseHeaderSize(rawHeader, HeaderPrefix_Area)
}
