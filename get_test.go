// get_test.go.

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
	"testing"

	"github.com/vault-thirteen/bit"
	"github.com/vault-thirteen/tester"
)

func Test_GetFormat(t *testing.T) {

	var format SbmFormat
	var formatExpected SbmFormat
	var sbm *Sbm
	var tst *tester.Test

	tst = tester.New(t)

	sbm = &Sbm{
		format: SbmFormat{
			version: SbmFormatVersion1,
		},
	}
	formatExpected = SbmFormat{
		version: SbmFormatVersion1,
	}

	// Get the Format.
	format = sbm.GetFormat()

	// Check the Result.
	tst.MustBeEqual(format, formatExpected)
}

func Test_GetArrayBytes(t *testing.T) {

	var arrayBytes []byte
	var arrayBytesExpected []byte
	var sbm *Sbm
	var tst *tester.Test

	tst = tester.New(t)

	sbm = &Sbm{
		pixelArray: SbmPixelArray{
			data: SbmPixelArrayData{
				bytes: []byte{
					1,
					2,
					3,
				},
			},
		},
	}
	arrayBytesExpected = []byte{
		1,
		2,
		3,
	}

	// Get the Format.
	arrayBytes = sbm.GetArrayBytes()

	// Check the Result.
	tst.MustBeEqual(arrayBytes, arrayBytesExpected)
}

func Test_GetArrayBits(t *testing.T) {

	var arrayBits []bit.Bit
	var arrayBitsExpected []bit.Bit
	var sbm *Sbm
	var tst *tester.Test

	tst = tester.New(t)

	sbm = &Sbm{
		pixelArray: SbmPixelArray{
			data: SbmPixelArrayData{
				bits: []bit.Bit{
					bit.One,
					bit.Zero,
					bit.One,
				},
			},
		},
	}
	arrayBitsExpected = []bit.Bit{
		bit.One,
		bit.Zero,
		bit.One,
	}

	// Get the Format.
	arrayBits = sbm.GetArrayBits()

	// Check the Result.
	tst.MustBeEqual(arrayBits, arrayBitsExpected)
}

func Test_GetArrayWidth(t *testing.T) {

	var arrayWidth uint
	var arrayWidthExpected uint
	var sbm *Sbm
	var tst *tester.Test

	tst = tester.New(t)

	sbm = &Sbm{
		pixelArray: SbmPixelArray{
			metaData: SbmPixelArrayMetaData{
				width: 123,
			},
		},
	}
	arrayWidthExpected = 123

	// Get the Format.
	arrayWidth = sbm.GetArrayWidth()

	// Check the Result.
	tst.MustBeEqual(arrayWidth, arrayWidthExpected)
}

func Test_GetArrayHeight(t *testing.T) {

	var arrayHeight uint
	var arrayHeightExpected uint
	var sbm *Sbm
	var tst *tester.Test

	tst = tester.New(t)

	sbm = &Sbm{
		pixelArray: SbmPixelArray{
			metaData: SbmPixelArrayMetaData{
				height: 456,
			},
		},
	}
	arrayHeightExpected = 456

	// Get the Format.
	arrayHeight = sbm.GetArrayHeight()

	// Check the Result.
	tst.MustBeEqual(arrayHeight, arrayHeightExpected)
}

func Test_GetArrayArea(t *testing.T) {

	var arrayArea uint
	var arrayAreaExpected uint
	var sbm *Sbm
	var tst *tester.Test

	tst = tester.New(t)

	sbm = &Sbm{
		pixelArray: SbmPixelArray{
			metaData: SbmPixelArrayMetaData{
				area: 789,
			},
		},
	}
	arrayAreaExpected = 789

	// Get the Format.
	arrayArea = sbm.GetArrayArea()

	// Check the Result.
	tst.MustBeEqual(arrayArea, arrayAreaExpected)
}
