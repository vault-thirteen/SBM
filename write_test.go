package sbm

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/vault-thirteen/auxie/bit"
	"github.com/vault-thirteen/tester"
)

func Test_Write(t *testing.T) {

	var buffer *bytes.Buffer
	var bufferReference []byte
	var err error
	var sbm *Sbm
	var tst *tester.Test

	tst = tester.New(t)

	sbm = &Sbm{
		format: SbmFormat{
			version: SbmFormatVersion1,
		},
		pixelArray: SbmPixelArray{
			data: SbmPixelArrayData{
				bytes: []byte{
					255,
					2,
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
					bit.Zero,
					bit.One,
				},
			},
			metaData: SbmPixelArrayMetaData{
				width:  3,
				height: 4,
				area:   12,
				header: SbmPixelArrayMetaDataHeader{
					width: SbmPixelArrayMetaDataHeaderData{
						topLeft:     1,
						topRight:    2,
						bottomLeft:  3,
						bottomRight: 0,
					},
					height: SbmPixelArrayMetaDataHeaderData{
						topLeft:     1,
						topRight:    3,
						bottomLeft:  4,
						bottomRight: 0,
					},
					area: SbmPixelArrayMetaDataHeaderData{
						topLeft:     1,
						topRight:    11,
						bottomLeft:  12,
						bottomRight: 0,
					},
				},
			},
		},
	}

	// Write.
	buffer = bytes.NewBuffer([]byte{})
	err = sbm.Write(buffer)
	tst.MustBeNoError(err)

	// Check the Data written.
	bufferReference = []byte(
		"SBM (SIMPLE BIT MAP)" + NL +
			"VERSION 1" + NL +
			"WIDTH 3 (1 + 2)" + NL +
			"HEIGHT 4 (1 + 3)" + NL +
			"AREA 12 (1 + 11)" + NL,
	)
	bufferReference = append(bufferReference,
		sbm.pixelArray.data.bytes...,
	)
	bufferReference = append(bufferReference,
		[]byte(NL)...,
	)
	bufferReference = append(bufferReference,
		[]byte(
			"WIDTH 3 (3 + 0)"+NL+
				"HEIGHT 4 (4 + 0)"+NL+
				"AREA 12 (12 + 0)"+NL,
		)...,
	)
	fmt.Println(string(bufferReference))
	fmt.Println(string(buffer.Bytes()))
	tst.MustBeEqual(buffer.Bytes(), bufferReference)
}

func Test_writeTopHeaders(t *testing.T) {

	var buffer *bytes.Buffer
	var bufferReference string
	var err error
	var sbm *Sbm
	var tst *tester.Test

	tst = tester.New(t)

	sbm = &Sbm{
		format: SbmFormat{
			version: SbmFormatVersion1,
		},
		pixelArray: SbmPixelArray{
			data: SbmPixelArrayData{
				bytes: []byte{},
				bits:  []bit.Bit{},
			},
			metaData: SbmPixelArrayMetaData{
				width:  3,
				height: 4,
				area:   12,
				header: SbmPixelArrayMetaDataHeader{
					width: SbmPixelArrayMetaDataHeaderData{
						topLeft:     1,
						topRight:    2,
						bottomLeft:  3,
						bottomRight: 0,
					},
					height: SbmPixelArrayMetaDataHeaderData{
						topLeft:     1,
						topRight:    3,
						bottomLeft:  4,
						bottomRight: 0,
					},
					area: SbmPixelArrayMetaDataHeaderData{
						topLeft:     1,
						topRight:    11,
						bottomLeft:  12,
						bottomRight: 0,
					},
				},
			},
		},
	}

	// Write.
	buffer = bytes.NewBufferString("")
	err = sbm.writeTopHeaders(buffer)
	tst.MustBeNoError(err)

	// Check the Data written.
	fmt.Println(buffer.String())
	bufferReference = "SBM (SIMPLE BIT MAP)" + NL +
		"VERSION 1" + NL +
		"WIDTH 3 (1 + 2)" + NL +
		"HEIGHT 4 (1 + 3)" + NL +
		"AREA 12 (1 + 11)" + NL
	tst.MustBeEqual(buffer.String(), bufferReference)
}

func Test_writeArrayData(t *testing.T) {

	var buffer *bytes.Buffer
	var bufferReference []byte
	var err error
	var sbm *Sbm
	var tst *tester.Test

	tst = tester.New(t)

	sbm = &Sbm{
		format: SbmFormat{
			version: SbmFormatVersion1,
		},
		pixelArray: SbmPixelArray{
			data: SbmPixelArrayData{
				bytes: []byte{
					255,
					2,
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
					bit.Zero,
					bit.One,
				},
			},
			metaData: SbmPixelArrayMetaData{
				width:  3,
				height: 4,
				area:   12,
				header: SbmPixelArrayMetaDataHeader{
					width: SbmPixelArrayMetaDataHeaderData{
						topLeft:     1,
						topRight:    2,
						bottomLeft:  3,
						bottomRight: 0,
					},
					height: SbmPixelArrayMetaDataHeaderData{
						topLeft:     1,
						topRight:    3,
						bottomLeft:  4,
						bottomRight: 0,
					},
					area: SbmPixelArrayMetaDataHeaderData{
						topLeft:     1,
						topRight:    11,
						bottomLeft:  12,
						bottomRight: 0,
					},
				},
			},
		},
	}

	// Write.
	buffer = bytes.NewBuffer([]byte{})
	err = sbm.writeArrayData(buffer)
	tst.MustBeNoError(err)

	// Check the Data written.
	bufferReference = sbm.pixelArray.data.bytes
	bufferReference = append(bufferReference,
		[]byte(NL)...,
	)
	fmt.Println(string(bufferReference))
	fmt.Println(string(buffer.Bytes()))
	tst.MustBeEqual(buffer.Bytes(), bufferReference)
}

func Test_writeBottomHeaders(t *testing.T) {

	var buffer *bytes.Buffer
	var bufferReference string
	var err error
	var sbm *Sbm
	var tst *tester.Test

	tst = tester.New(t)

	sbm = &Sbm{
		format: SbmFormat{
			version: SbmFormatVersion1,
		},
		pixelArray: SbmPixelArray{
			data: SbmPixelArrayData{
				bytes: []byte{},
				bits:  []bit.Bit{},
			},
			metaData: SbmPixelArrayMetaData{
				width:  3,
				height: 4,
				area:   12,
				header: SbmPixelArrayMetaDataHeader{
					width: SbmPixelArrayMetaDataHeaderData{
						topLeft:     1,
						topRight:    2,
						bottomLeft:  3,
						bottomRight: 0,
					},
					height: SbmPixelArrayMetaDataHeaderData{
						topLeft:     1,
						topRight:    3,
						bottomLeft:  4,
						bottomRight: 0,
					},
					area: SbmPixelArrayMetaDataHeaderData{
						topLeft:     1,
						topRight:    11,
						bottomLeft:  12,
						bottomRight: 0,
					},
				},
			},
		},
	}

	// Write.
	buffer = bytes.NewBufferString("")
	err = sbm.writeBottomHeaders(buffer)
	tst.MustBeNoError(err)

	// Check the Data written.
	fmt.Println(buffer.String())
	bufferReference = "WIDTH 3 (3 + 0)" + NL +
		"HEIGHT 4 (4 + 0)" + NL +
		"AREA 12 (12 + 0)" + NL
	tst.MustBeEqual(buffer.String(), bufferReference)
}
