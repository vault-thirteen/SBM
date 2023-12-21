package sbm

import (
	"testing"

	"github.com/vault-thirteen/auxie/bit"
	"github.com/vault-thirteen/auxie/tester"
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

	// Get the format.
	format = sbm.GetFormat()

	// Check the result.
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

	// Get the format.
	arrayBytes = sbm.GetArrayBytes()

	// Check the result.
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

	// Get the format.
	arrayBits = sbm.GetArrayBits()

	// Check the result.
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

	// Get the format.
	arrayWidth = sbm.GetArrayWidth()

	// Check the result.
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

	// Get the format.
	arrayHeight = sbm.GetArrayHeight()

	// Check the result.
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

	// Get the format.
	arrayArea = sbm.GetArrayArea()

	// Check the result.
	tst.MustBeEqual(arrayArea, arrayAreaExpected)
}
