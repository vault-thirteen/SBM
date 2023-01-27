// new_test.go.

package sbm

import (
	"bytes"
	"io"
	"testing"

	"github.com/vault-thirteen/bit"
	"github.com/vault-thirteen/tester"
)

func Test_NewFromBitsArray(t *testing.T) {

	var err error
	var tst *tester.Test

	tst = tester.New(t)

	// Test #1. Zero Width.
	_, err = NewFromBitsArray(
		[]bit.Bit{},
		0,
		123,
	)
	tst.MustBeAnError(err)

	// Test #2. Zero Height.
	_, err = NewFromBitsArray(
		[]bit.Bit{},
		123,
		0,
	)
	tst.MustBeAnError(err)

	// Test #3. Bits Array Size is bad.
	_, err = NewFromBitsArray(
		[]bit.Bit{
			bit.One,
		},
		3,
		3,
	)
	tst.MustBeAnError(err)

	// Test #4. No Error.
	_, err = NewFromBitsArray(
		[]bit.Bit{
			// Byte #1.
			bit.One,
			bit.One,
			bit.One,
			bit.One,
			bit.One,
			bit.One,
			bit.One,
			bit.One,
			// Byte #2.
			bit.One,
			bit.One,
			bit.One,
			bit.One,
		},
		3,
		4,
	)
	tst.MustBeNoError(err)
}

func Test_NewFromBytesArray(t *testing.T) {

	var err error
	var tst *tester.Test

	tst = tester.New(t)

	// Test #1. Zero Width.
	_, err = NewFromBytesArray(
		[]byte{},
		0,
		123,
	)
	tst.MustBeAnError(err)

	// Test #2. Zero Height.
	_, err = NewFromBytesArray(
		[]byte{},
		123,
		0,
	)
	tst.MustBeAnError(err)

	// Test #3. Bytes Array Size is too small.
	_, err = NewFromBytesArray(
		[]byte{ // 1 < 2.
			255,
		},
		3, // 3 x 3 = 9 Bits => 2 Bytes are needed to store an Array.
		3,
	)
	tst.MustBeAnError(err)

	// Test #4. Bytes Array Size is too big.
	_, err = NewFromBytesArray(
		[]byte{ // 3 > 2.
			255,
			255,
			255,
		},
		3, // 3 x 3 = 9 Bits => 2 Bytes are needed to store an Array.
		3,
	)
	tst.MustBeAnError(err)

	// Test #5. No Error.
	_, err = NewFromBytesArray(
		[]byte{ // 2 = 2.
			255,
			255,
		},
		3, // 3 x 3 = 9 Bits => 2 Bytes are needed to store an Array.
		3,
	)
	tst.MustBeNoError(err)
}

func Test_newFromBitsArray(t *testing.T) {

	var err error
	var result *Sbm
	var resultExpected *Sbm
	var tst *tester.Test

	tst = tester.New(t)

	// Test #1. Normal Test.
	resultExpected = &Sbm{
		format: SbmFormat{
			version: SbmFormatVersion1,
		},
		pixelArray: SbmPixelArray{
			data: SbmPixelArrayData{
				bytes: []byte{
					255,
					15,
				},
				bits: []bit.Bit{
					// Byte #1.
					bit.One,
					bit.One,
					bit.One,
					bit.One,
					bit.One,
					bit.One,
					bit.One,
					bit.One,
					// Byte #2.
					bit.One,
					bit.One,
					bit.One,
					bit.One,
				},
			},
			metaData: SbmPixelArrayMetaData{
				width:  3,
				height: 4,
				area:   12,
				header: SbmPixelArrayMetaDataHeader{
					width:  SbmPixelArrayMetaDataHeaderData{}, // Random.
					height: SbmPixelArrayMetaDataHeaderData{}, // Random.
					area:   SbmPixelArrayMetaDataHeaderData{}, // Random.
				},
			},
		},
	}
	result, err = newFromBitsArray(
		[]bit.Bit{
			// Byte #1.
			bit.One,
			bit.One,
			bit.One,
			bit.One,
			bit.One,
			bit.One,
			bit.One,
			bit.One,
			// Byte #2.
			bit.One,
			bit.One,
			bit.One,
			bit.One,
		},
		3,
		4,
	)
	tst.MustBeNoError(err)
	tst.MustBeEqual(result.format, resultExpected.format)
	tst.MustBeEqual(result.pixelArray.data, resultExpected.pixelArray.data)
	tst.MustBeEqual(result.pixelArray.metaData.width, resultExpected.pixelArray.metaData.width)
	tst.MustBeEqual(result.pixelArray.metaData.height, resultExpected.pixelArray.metaData.height)
	tst.MustBeEqual(result.pixelArray.metaData.area, resultExpected.pixelArray.metaData.area)

	// Check random Meta-Data.
	tst.MustBeEqual(result.pixelArray.metaData.header.width.topLeft+
		result.pixelArray.metaData.header.width.topRight,
		resultExpected.pixelArray.metaData.width)
	tst.MustBeEqual(result.pixelArray.metaData.header.width.bottomLeft+
		result.pixelArray.metaData.header.width.bottomRight,
		resultExpected.pixelArray.metaData.width)
	tst.MustBeEqual(result.pixelArray.metaData.header.height.topLeft+
		result.pixelArray.metaData.header.height.topRight,
		resultExpected.pixelArray.metaData.height)
	tst.MustBeEqual(result.pixelArray.metaData.header.height.bottomLeft+
		result.pixelArray.metaData.header.height.bottomRight,
		resultExpected.pixelArray.metaData.height)
	tst.MustBeEqual(result.pixelArray.metaData.header.area.topLeft+
		result.pixelArray.metaData.header.area.topRight,
		resultExpected.pixelArray.metaData.area)
	tst.MustBeEqual(result.pixelArray.metaData.header.area.bottomLeft+
		result.pixelArray.metaData.header.area.bottomRight,
		resultExpected.pixelArray.metaData.area)
}

func Test_newFromBytesArray(t *testing.T) {

	var err error
	var result *Sbm
	var resultExpected *Sbm
	var tst *tester.Test

	tst = tester.New(t)

	// Test #1. Normal Test.
	resultExpected = &Sbm{
		format: SbmFormat{
			version: SbmFormatVersion1,
		},
		pixelArray: SbmPixelArray{
			data: SbmPixelArrayData{
				bytes: []byte{
					255,
					1,
				},
				bits: []bit.Bit{
					// Byte #1.
					bit.One,
					bit.One,
					bit.One,
					bit.One,
					bit.One,
					bit.One,
					bit.One,
					bit.One,
					// Byte #2.
					bit.One,
				},
			},
			metaData: SbmPixelArrayMetaData{
				width:  3,
				height: 3,
				area:   9,
				header: SbmPixelArrayMetaDataHeader{
					width:  SbmPixelArrayMetaDataHeaderData{}, // Random.
					height: SbmPixelArrayMetaDataHeaderData{}, // Random.
					area:   SbmPixelArrayMetaDataHeaderData{}, // Random.
				},
			},
		},
	}
	result, err = newFromBytesArray(
		[]byte{
			255,
			1,
		},
		3,
		3,
	)
	tst.MustBeNoError(err)
	tst.MustBeEqual(result.format, resultExpected.format)
	tst.MustBeEqual(result.pixelArray.data, resultExpected.pixelArray.data)
	tst.MustBeEqual(result.pixelArray.metaData.width, resultExpected.pixelArray.metaData.width)
	tst.MustBeEqual(result.pixelArray.metaData.height, resultExpected.pixelArray.metaData.height)
	tst.MustBeEqual(result.pixelArray.metaData.area, resultExpected.pixelArray.metaData.area)

	// Check random Meta-Data.
	tst.MustBeEqual(result.pixelArray.metaData.header.width.topLeft+
		result.pixelArray.metaData.header.width.topRight,
		resultExpected.pixelArray.metaData.width)
	tst.MustBeEqual(result.pixelArray.metaData.header.width.bottomLeft+
		result.pixelArray.metaData.header.width.bottomRight,
		resultExpected.pixelArray.metaData.width)
	tst.MustBeEqual(result.pixelArray.metaData.header.height.topLeft+
		result.pixelArray.metaData.header.height.topRight,
		resultExpected.pixelArray.metaData.height)
	tst.MustBeEqual(result.pixelArray.metaData.header.height.bottomLeft+
		result.pixelArray.metaData.header.height.bottomRight,
		resultExpected.pixelArray.metaData.height)
	tst.MustBeEqual(result.pixelArray.metaData.header.area.topLeft+
		result.pixelArray.metaData.header.area.topRight,
		resultExpected.pixelArray.metaData.area)
	tst.MustBeEqual(result.pixelArray.metaData.header.area.bottomLeft+
		result.pixelArray.metaData.header.area.bottomRight,
		resultExpected.pixelArray.metaData.area)
}

func Test_newFromBitsAndBytesArrays(t *testing.T) {

	var err error
	var result *Sbm
	var resultExpected *Sbm
	var tst *tester.Test

	tst = tester.New(t)

	// Test #1. Normal Test.
	resultExpected = &Sbm{
		format: SbmFormat{
			version: SbmFormatVersion1,
		},
		pixelArray: SbmPixelArray{
			data: SbmPixelArrayData{
				bytes: []byte{
					255,
					15,
				},
				bits: []bit.Bit{
					// Byte #1.
					bit.One,
					bit.One,
					bit.One,
					bit.One,
					bit.One,
					bit.One,
					bit.One,
					bit.One,
					// Byte #2.
					bit.One,
					bit.One,
					bit.One,
					bit.One,
				},
			},
			metaData: SbmPixelArrayMetaData{
				width:  3,
				height: 4,
				area:   12,
				header: SbmPixelArrayMetaDataHeader{
					width:  SbmPixelArrayMetaDataHeaderData{}, // Random.
					height: SbmPixelArrayMetaDataHeaderData{}, // Random.
					area:   SbmPixelArrayMetaDataHeaderData{}, // Random.
				},
			},
		},
	}
	result, err = newFromBitsAndBytesArrays(
		[]bit.Bit{
			// Byte #1.
			bit.One,
			bit.One,
			bit.One,
			bit.One,
			bit.One,
			bit.One,
			bit.One,
			bit.One,
			// Byte #2.
			bit.One,
			bit.One,
			bit.One,
			bit.One,
		},
		[]byte{
			255,
			15,
		},
		3,
		4,
		12,
	)
	tst.MustBeNoError(err)
	tst.MustBeEqual(result.format, resultExpected.format)
	tst.MustBeEqual(result.pixelArray.data, resultExpected.pixelArray.data)
	tst.MustBeEqual(result.pixelArray.metaData.width, resultExpected.pixelArray.metaData.width)
	tst.MustBeEqual(result.pixelArray.metaData.height, resultExpected.pixelArray.metaData.height)
	tst.MustBeEqual(result.pixelArray.metaData.area, resultExpected.pixelArray.metaData.area)

	// Check random Meta-Data.
	tst.MustBeEqual(result.pixelArray.metaData.header.width.topLeft+
		result.pixelArray.metaData.header.width.topRight,
		resultExpected.pixelArray.metaData.width)
	tst.MustBeEqual(result.pixelArray.metaData.header.width.bottomLeft+
		result.pixelArray.metaData.header.width.bottomRight,
		resultExpected.pixelArray.metaData.width)
	tst.MustBeEqual(result.pixelArray.metaData.header.height.topLeft+
		result.pixelArray.metaData.header.height.topRight,
		resultExpected.pixelArray.metaData.height)
	tst.MustBeEqual(result.pixelArray.metaData.header.height.bottomLeft+
		result.pixelArray.metaData.header.height.bottomRight,
		resultExpected.pixelArray.metaData.height)
	tst.MustBeEqual(result.pixelArray.metaData.header.area.topLeft+
		result.pixelArray.metaData.header.area.topRight,
		resultExpected.pixelArray.metaData.area)
	tst.MustBeEqual(result.pixelArray.metaData.header.area.bottomLeft+
		result.pixelArray.metaData.header.area.bottomRight,
		resultExpected.pixelArray.metaData.area)
}

func Test_NewFromStream(t *testing.T) {

	var data []byte
	var err error
	var reader io.Reader
	var sbm *Sbm
	var sbmExpected Sbm
	var tst *tester.Test

	tst = tester.New(t)

	// Test #1. Positive.
	data = []byte(
		"SBM (SIMPLE BIT MAP)" + NL +
			"VERSION 1" + NL +
			"WIDTH 3 (2 + 1)" + NL +
			"HEIGHT 4 (3 + 1)" + NL +
			"AREA 12 (11 + 1)" + NL,
	)
	data = append(data, []byte{
		255,
		255,
	}...)
	data = append(data, []byte(NL)...)
	data = append(data, []byte("WIDTH 3 (0 + 3)"+NL+
		"HEIGHT 4 (0 + 4)"+NL+
		"AREA 12 (10 + 2)"+NL,
	)...)
	data = append(data, []byte("JUNK")...)
	sbmExpected = Sbm{
		format: SbmFormat{
			version: SbmFormatVersion1,
		},
		pixelArray: SbmPixelArray{
			data: SbmPixelArrayData{
				bits: []bit.Bit{
					bit.One, // 1.
					bit.One,
					bit.One,
					bit.One,
					bit.One,
					bit.One,
					bit.One,
					bit.One,
					bit.One,
					bit.One,
					bit.One,
					bit.One, // 12.
				},
				bytes: []byte{
					255,
					15,
				},
			},
			metaData: SbmPixelArrayMetaData{
				width:  3,
				height: 4,
				area:   12,
				header: SbmPixelArrayMetaDataHeader{
					width: SbmPixelArrayMetaDataHeaderData{
						topLeft:     2,
						topRight:    1,
						bottomLeft:  0,
						bottomRight: 3,
					},
					height: SbmPixelArrayMetaDataHeaderData{
						topLeft:     3,
						topRight:    1,
						bottomLeft:  0,
						bottomRight: 4,
					},
					area: SbmPixelArrayMetaDataHeaderData{
						topLeft:     11,
						topRight:    1,
						bottomLeft:  10,
						bottomRight: 2,
					},
				},
			},
		},
	}
	reader = bytes.NewReader(data)
	sbm, err = NewFromStream(reader)
	tst.MustBeNoError(err)
	tst.MustBeEqual(*sbm, sbmExpected)

	// Test #2. Negative. Total Junk.
	data = []byte(
		"SOME JUNK IS HERE",
	)
	reader = bytes.NewReader(data)
	sbm, err = NewFromStream(reader)
	tst.MustBeAnError(err)
	tst.MustBeEqual(sbm, (*Sbm)(nil))

	// Test #3. Negative. Bad Array.
	data = []byte(
		"SBM (SIMPLE BIT MAP)" + NL +
			"VERSION 1" + NL +
			"WIDTH 3 (2 + 1)" + NL +
			"HEIGHT 4 (3 + 1)" + NL +
			"AREA 12 (11 + 1)" + NL,
	)
	data = append(data, []byte{}...)
	data = append(data, []byte(NL)...)
	reader = bytes.NewReader(data)
	sbm, err = NewFromStream(reader)
	tst.MustBeAnError(err)
	tst.MustBeEqual(sbm, (*Sbm)(nil))

	// Test #3. Negative. Bad Bottom Headers.
	data = []byte(
		"SBM (SIMPLE BIT MAP)" + NL +
			"VERSION 1" + NL +
			"WIDTH 3 (2 + 1)" + NL +
			"HEIGHT 4 (3 + 1)" + NL +
			"AREA 12 (11 + 1)" + NL,
	)
	data = append(data, []byte{
		255,
		255,
	}...)
	data = append(data, []byte(NL)...)
	data = append(data, []byte("WIDTH Q (QQ + QQ)"+NL)...)
	reader = bytes.NewReader(data)
	sbm, err = NewFromStream(reader)
	tst.MustBeAnError(err)
	tst.MustBeEqual(sbm, (*Sbm)(nil))
}
