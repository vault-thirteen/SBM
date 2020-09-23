// read_test.go.

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
	"fmt"
	"io"
	"math"
	"testing"

	"github.com/vault-thirteen/bit"
	rdr "github.com/vault-thirteen/reader"
	"github.com/vault-thirteen/tester"
)

func Test_readTopHeaders(t *testing.T) {

	// Test Type (Kind).
	const (
		TestKindMustBeNoError    = 1
		TestKindMustBeAnyError   = 2
		TestKindMustBeExactError = 3
	)

	type Test struct {
		kind              byte
		data              []byte
		sbm               *Sbm
		sbmExpected       Sbm
		errorTextExpected string
	}

	var err error
	var lineReader *rdr.Reader
	var reader io.Reader
	var test Test
	var testIdx int
	var tests []Test
	var tst *tester.Test

	tst = tester.New(t)
	tests = make([]Test, 0)

	// Test #1. Positive.
	test = Test{
		kind: TestKindMustBeNoError,
		data: []byte(
			"SBM (SIMPLE BIT MAP)" + NL +
				"VERSION 1" + NL +
				"WIDTH 123 (100 + 23)" + NL +
				"HEIGHT 456 (450 + 6)" + NL +
				"AREA 56088 (56000 + 88)" + NL +
				"SOME JUNK IS HERE",
		),
		sbm: new(Sbm),
		sbmExpected: Sbm{
			format: SbmFormat{
				version: SbmFormatVersion1,
			},
			pixelArray: SbmPixelArray{
				data: SbmPixelArrayData{},
				metaData: SbmPixelArrayMetaData{
					width:  123,
					height: 456,
					area:   56088,
					header: SbmPixelArrayMetaDataHeader{
						width: SbmPixelArrayMetaDataHeaderData{
							topLeft:  100,
							topRight: 23,
						},
						height: SbmPixelArrayMetaDataHeaderData{
							topLeft:  450,
							topRight: 6,
						},
						area: SbmPixelArrayMetaDataHeaderData{
							topLeft:  56000,
							topRight: 88,
						},
					},
				},
			},
		},
		errorTextExpected: "",
	}
	tests = append(tests, test)

	// Test #2. Negative. Bad Format Header.
	test = Test{
		kind: TestKindMustBeExactError,
		data: []byte(
			"SOME JUNK IS HERE",
		),
		sbm:               new(Sbm),
		sbmExpected:       Sbm{},
		errorTextExpected: io.EOF.Error(),
	}
	tests = append(tests, test)

	// Test #3. Negative. Bad Format Header.
	test = Test{
		kind: TestKindMustBeExactError,
		data: []byte(
			"SBM (SIMPLE HACKED MAP)\r\n" +
				"JUNK",
		),
		sbm:               new(Sbm),
		sbmExpected:       Sbm{},
		errorTextExpected: ErrFormat,
	}
	tests = append(tests, test)

	// Test #4. Negative. Bad Version Header.
	test = Test{
		kind: TestKindMustBeExactError,
		data: []byte(
			"SBM (SIMPLE BIT MAP)\r\n" +
				"JUNK",
		),
		sbm:               new(Sbm),
		sbmExpected:       Sbm{},
		errorTextExpected: io.EOF.Error(),
	}
	tests = append(tests, test)

	// Test #5. Negative. Bad Version Header.
	test = Test{
		kind: TestKindMustBeAnyError,
		data: []byte(
			"SBM (SIMPLE BIT MAP)\r\n" +
				"VERSION JUNK\r\n" +
				"XYZ",
		),
		sbm:               new(Sbm),
		sbmExpected:       Sbm{},
		errorTextExpected: "",
	}
	tests = append(tests, test)

	// Test #6. Negative. Bad Version Header.
	test = Test{
		kind: TestKindMustBeExactError,
		data: []byte(
			"SBM (SIMPLE BIT MAP)\r\n" +
				"VERSION 123\r\n" +
				"JUNK",
		),
		sbm:               new(Sbm),
		sbmExpected:       Sbm{},
		errorTextExpected: ErrVersion,
	}
	tests = append(tests, test)

	// Test #7. Negative. Bad Width Header.
	test = Test{
		kind: TestKindMustBeExactError,
		data: []byte(
			"SBM (SIMPLE BIT MAP)\r\n" +
				"VERSION 1\r\n" +
				"JUNK",
		),
		sbm: new(Sbm),
		sbmExpected: Sbm{
			format: SbmFormat{
				version: SbmFormatVersion1,
			},
		},
		errorTextExpected: io.EOF.Error(),
	}
	tests = append(tests, test)

	// Test #8. Negative. Bad Width Header.
	test = Test{
		kind: TestKindMustBeExactError,
		data: []byte(
			"SBM (SIMPLE BIT MAP)\r\n" +
				"VERSION 1\r\n" +
				"WIDTH JUNK\r\n" +
				"XYZ",
		),
		sbm: new(Sbm),
		sbmExpected: Sbm{
			format: SbmFormat{
				version: SbmFormatVersion1,
			},
		},
		errorTextExpected: ErrHeaderSyntax,
	}
	tests = append(tests, test)

	// Test #9. Negative. Bad Height Header.
	test = Test{
		kind: TestKindMustBeExactError,
		data: []byte(
			"SBM (SIMPLE BIT MAP)\r\n" +
				"VERSION 1\r\n" +
				"WIDTH 10 (9 + 1)\r\n" +
				"XYZ",
		),
		sbm: new(Sbm),
		sbmExpected: Sbm{
			format: SbmFormat{
				version: SbmFormatVersion1,
			},
			pixelArray: SbmPixelArray{
				metaData: SbmPixelArrayMetaData{
					width: 10,
					header: SbmPixelArrayMetaDataHeader{
						width: SbmPixelArrayMetaDataHeaderData{
							topLeft:  9,
							topRight: 1,
						},
					},
				},
			},
		},
		errorTextExpected: io.EOF.Error(),
	}
	tests = append(tests, test)

	// Test #10. Negative. Bad Height Header.
	test = Test{
		kind: TestKindMustBeExactError,
		data: []byte(
			"SBM (SIMPLE BIT MAP)\r\n" +
				"VERSION 1\r\n" +
				"WIDTH 10 (9 + 1)\r\n" +
				"HEIGHT XYZ\r\n" +
				"QWERTY",
		),
		sbm: new(Sbm),
		sbmExpected: Sbm{
			format: SbmFormat{
				version: SbmFormatVersion1,
			},
			pixelArray: SbmPixelArray{
				metaData: SbmPixelArrayMetaData{
					width: 10,
					header: SbmPixelArrayMetaDataHeader{
						width: SbmPixelArrayMetaDataHeaderData{
							topLeft:  9,
							topRight: 1,
						},
					},
				},
			},
		},
		errorTextExpected: ErrHeaderSyntax,
	}
	tests = append(tests, test)

	// Test #11. Negative. Bad Area Header.
	test = Test{
		kind: TestKindMustBeExactError,
		data: []byte(
			"SBM (SIMPLE BIT MAP)\r\n" +
				"VERSION 1\r\n" +
				"WIDTH 10 (9 + 1)\r\n" +
				"HEIGHT 15 (14 + 1)\r\n" +
				"QWERTY",
		),
		sbm: new(Sbm),
		sbmExpected: Sbm{
			format: SbmFormat{
				version: SbmFormatVersion1,
			},
			pixelArray: SbmPixelArray{
				metaData: SbmPixelArrayMetaData{
					width:  10,
					height: 15,
					header: SbmPixelArrayMetaDataHeader{
						width: SbmPixelArrayMetaDataHeaderData{
							topLeft:  9,
							topRight: 1,
						},
						height: SbmPixelArrayMetaDataHeaderData{
							topLeft:  14,
							topRight: 1,
						},
					},
				},
			},
		},
		errorTextExpected: io.EOF.Error(),
	}
	tests = append(tests, test)

	// Test #12. Negative. Bad Area Header.
	test = Test{
		kind: TestKindMustBeExactError,
		data: []byte(
			"SBM (SIMPLE BIT MAP)\r\n" +
				"VERSION 1\r\n" +
				"WIDTH 10 (9 + 1)\r\n" +
				"HEIGHT 15 (14 + 1)\r\n" +
				"AREA XYZ\r\n" +
				"QWERTY",
		),
		sbm: new(Sbm),
		sbmExpected: Sbm{
			format: SbmFormat{
				version: SbmFormatVersion1,
			},
			pixelArray: SbmPixelArray{
				metaData: SbmPixelArrayMetaData{
					width:  10,
					height: 15,
					header: SbmPixelArrayMetaDataHeader{
						width: SbmPixelArrayMetaDataHeaderData{
							topLeft:  9,
							topRight: 1,
						},
						height: SbmPixelArrayMetaDataHeaderData{
							topLeft:  14,
							topRight: 1,
						},
					},
				},
			},
		},
		errorTextExpected: ErrHeaderSyntax,
	}
	tests = append(tests, test)

	// Test #13. Negative. Bad Area.
	test = Test{
		kind: TestKindMustBeExactError,
		data: []byte(
			"SBM (SIMPLE BIT MAP)" + NL +
				"VERSION 1" + NL +
				"WIDTH 123 (100 + 23)" + NL +
				"HEIGHT 456 (450 + 6)" + NL +
				"AREA 56089 (56000 + 89)" + NL +
				"SOME JUNK IS HERE",
		),
		sbm: new(Sbm),
		sbmExpected: Sbm{
			format: SbmFormat{
				version: SbmFormatVersion1,
			},
			pixelArray: SbmPixelArray{
				metaData: SbmPixelArrayMetaData{
					width:  123,
					height: 456,
					header: SbmPixelArrayMetaDataHeader{
						width: SbmPixelArrayMetaDataHeaderData{
							topLeft:  100,
							topRight: 23,
						},
						height: SbmPixelArrayMetaDataHeaderData{
							topLeft:  450,
							topRight: 6,
						},
						area: SbmPixelArrayMetaDataHeaderData{},
					},
				},
			},
		},
		errorTextExpected: ErrAreaMismatch,
	}
	tests = append(tests, test)

	// Run the Tests.
	for testIdx, test = range tests {

		// Run the Action.
		reader = bytes.NewReader(test.data)
		lineReader = rdr.NewReader(reader)
		err = test.sbm.readTopHeaders(lineReader)

		// Check.
		switch test.kind {

		case TestKindMustBeNoError:
			tst.MustBeNoError(err)

		case TestKindMustBeAnyError:
			tst.MustBeAnError(err)

		case TestKindMustBeExactError:
			tst.MustBeAnError(err)
			tst.MustBeEqual(err.Error(), test.errorTextExpected)

		default:
			t.FailNow()
		}

		tst.MustBeDifferent(test.sbm, nil)
		tst.MustBeEqual(*(test.sbm), test.sbmExpected)

		fmt.Printf("[%v]", testIdx+1)
	}
	fmt.Println()
}

func Test_readArrayData(t *testing.T) {

	var data []byte
	var err error
	var lineReader *rdr.Reader
	var reader io.Reader
	var sbm *Sbm
	var sbmExpected Sbm
	var tst *tester.Test

	tst = tester.New(t)

	// Test #1. Positive. One Byte, eight Bits.
	data = []byte{255}
	data = append(data, []byte(NL+"SOME JUNK")...)
	sbm = new(Sbm)
	sbm.pixelArray.metaData.area = 8
	sbmExpected = Sbm{
		pixelArray: SbmPixelArray{
			data: SbmPixelArrayData{
				bits: []bit.Bit{
					bit.One,
					bit.One,
					bit.One,
					bit.One,
					bit.One,
					bit.One,
					bit.One,
					bit.One,
				},
				bytes: []byte{
					255,
				},
			},
			metaData: SbmPixelArrayMetaData{
				area: 8,
			},
		},
	}
	reader = bytes.NewReader(data)
	lineReader = rdr.NewReader(reader)
	err = sbm.readArrayData(lineReader)
	tst.MustBeNoError(err)
	tst.MustBeEqual(*sbm, sbmExpected)

	// Test #2. Positive. Two Bytes, nine Bits.
	data = []byte{255, 1}
	data = append(data, []byte(NL+"SOME JUNK")...)
	sbm = new(Sbm)
	sbm.pixelArray.metaData.area = 9
	sbmExpected = Sbm{
		pixelArray: SbmPixelArray{
			data: SbmPixelArrayData{
				bits: []bit.Bit{
					bit.One,
					bit.One,
					bit.One,
					bit.One,
					bit.One,
					bit.One,
					bit.One,
					bit.One,
					bit.One,
				},
				bytes: []byte{
					255,
					1,
				},
			},
			metaData: SbmPixelArrayMetaData{
				area: 9,
			},
		},
	}
	reader = bytes.NewReader(data)
	lineReader = rdr.NewReader(reader)
	err = sbm.readArrayData(lineReader)
	tst.MustBeNoError(err)
	tst.MustBeEqual(*sbm, sbmExpected)

	// Test #3. Negative. Bad Separator.
	data = []byte{
		255,
		1,
	}
	data = append(data, []byte("SOME JUNK")...)
	sbm = new(Sbm)
	sbm.pixelArray.metaData.area = 9
	sbmExpected = Sbm{
		pixelArray: SbmPixelArray{
			metaData: SbmPixelArrayMetaData{
				area: 9,
			},
		},
	}
	reader = bytes.NewReader(data)
	lineReader = rdr.NewReader(reader)
	err = sbm.readArrayData(lineReader)
	tst.MustBeAnError(err)
	tst.MustBeEqual(err.Error(), ErrBadSeparator)
	tst.MustBeEqual(*sbm, sbmExpected)

	// Test #4. Negative. Area is too big.
	// Actually, this Test can not be performed with the now-a-days Hardware.
	// It is likely to raise an Exception during the Memory Allocation.
	data = []byte{
		255,
		1,
	}
	data = append(data, []byte("SOME JUNK")...)
	sbm = new(Sbm)
	sbm.pixelArray.metaData.area = math.MaxUint64
	sbmExpected = Sbm{
		pixelArray: SbmPixelArray{
			metaData: SbmPixelArrayMetaData{
				area: math.MaxUint64,
			},
		},
	}
	reader = bytes.NewReader(data)
	lineReader = rdr.NewReader(reader)
	err = sbm.readArrayData(lineReader)
	tst.MustBeAnError(err)
	fmt.Println(err)
	tst.MustBeEqual(*sbm, sbmExpected)

	// Test #5. Negative. Data is not enough.
	data = []byte{255}
	data = append(data, []byte(NL)...)
	sbm = new(Sbm)
	sbm.pixelArray.metaData.area = 100
	sbmExpected = Sbm{
		pixelArray: SbmPixelArray{
			metaData: SbmPixelArrayMetaData{
				area: 100,
			},
		},
	}
	reader = bytes.NewReader(data)
	lineReader = rdr.NewReader(reader)
	err = sbm.readArrayData(lineReader)
	tst.MustBeAnError(err)
	tst.MustBeEqual(err.Error(), io.ErrUnexpectedEOF.Error())
	tst.MustBeEqual(*sbm, sbmExpected)
}

func Test_readBottomHeaders(t *testing.T) {

	// Test Type (Kind).
	const (
		TestKindMustBeNoError    = 1
		TestKindMustBeAnyError   = 2
		TestKindMustBeExactError = 3
	)

	type Test struct {
		kind              byte
		data              []byte
		sbm               *Sbm
		sbmExpected       Sbm
		errorTextExpected string
	}

	var err error
	var lineReader *rdr.Reader
	var reader io.Reader
	var test Test
	var testIdx int
	var tests []Test
	var tst *tester.Test

	tst = tester.New(t)
	tests = make([]Test, 0)

	// Test #1. Positive.
	test = Test{
		kind: TestKindMustBeNoError,
		data: []byte(
			"WIDTH 123 (100 + 23)" + NL +
				"HEIGHT 456 (450 + 6)" + NL +
				"AREA 56088 (56000 + 88)" + NL +
				"SOME JUNK IS HERE",
		),
		sbm: &Sbm{
			pixelArray: SbmPixelArray{
				metaData: SbmPixelArrayMetaData{
					width:  123,
					height: 456,
					area:   56088,
				},
			},
		},
		sbmExpected: Sbm{
			pixelArray: SbmPixelArray{
				data: SbmPixelArrayData{},
				metaData: SbmPixelArrayMetaData{
					width:  123,
					height: 456,
					area:   56088,
					header: SbmPixelArrayMetaDataHeader{
						width: SbmPixelArrayMetaDataHeaderData{
							bottomLeft:  100,
							bottomRight: 23,
						},
						height: SbmPixelArrayMetaDataHeaderData{
							bottomLeft:  450,
							bottomRight: 6,
						},
						area: SbmPixelArrayMetaDataHeaderData{
							bottomLeft:  56000,
							bottomRight: 88,
						},
					},
				},
			},
		},
		errorTextExpected: "",
	}
	tests = append(tests, test)

	// Test #2. Negative. Bad Width Header.
	test = Test{
		kind: TestKindMustBeExactError,
		data: []byte(
			"JUNK",
		),
		sbm:               new(Sbm),
		sbmExpected:       Sbm{},
		errorTextExpected: io.EOF.Error(),
	}
	tests = append(tests, test)

	// Test #3. Negative. Bad Width Header.
	test = Test{
		kind: TestKindMustBeExactError,
		data: []byte(
			"WIDTH JUNK\r\n" +
				"XYZ",
		),
		sbm:               new(Sbm),
		sbmExpected:       Sbm{},
		errorTextExpected: ErrHeaderSyntax,
	}
	tests = append(tests, test)

	// Test #4. Negative.
	// Conflict with existing Width (e.g. taken from the Top header).
	// The old Value of the Object is not changed.
	test = Test{
		kind: TestKindMustBeExactError,
		data: []byte(
			"WIDTH 100 (99 + 1)\r\n" +
				"XYZ",
		),
		sbm: &Sbm{
			pixelArray: SbmPixelArray{
				metaData: SbmPixelArrayMetaData{
					width: 10,
				},
			},
		},
		sbmExpected: Sbm{
			pixelArray: SbmPixelArray{
				metaData: SbmPixelArrayMetaData{
					width: 10,
				},
			},
		},
		errorTextExpected: ErrBottomHeaderMismatch,
	}
	tests = append(tests, test)

	// Test #5. Negative. Bad Height Header.
	test = Test{
		kind: TestKindMustBeExactError,
		data: []byte(
			"WIDTH 10 (9 + 1)\r\n" +
				"XYZ",
		),
		sbm: &Sbm{
			pixelArray: SbmPixelArray{
				metaData: SbmPixelArrayMetaData{
					width: 10,
					header: SbmPixelArrayMetaDataHeader{
						width: SbmPixelArrayMetaDataHeaderData{
							bottomLeft:  9,
							bottomRight: 1,
						},
					},
				},
			},
		},
		sbmExpected: Sbm{
			pixelArray: SbmPixelArray{
				metaData: SbmPixelArrayMetaData{
					width: 10,
					header: SbmPixelArrayMetaDataHeader{
						width: SbmPixelArrayMetaDataHeaderData{
							bottomLeft:  9,
							bottomRight: 1,
						},
					},
				},
			},
		},
		errorTextExpected: io.EOF.Error(),
	}
	tests = append(tests, test)

	// Test #6. Negative. Bad Height Header.
	test = Test{
		kind: TestKindMustBeExactError,
		data: []byte(
			"WIDTH 10 (9 + 1)\r\n" +
				"HEIGHT XYZ\r\n" +
				"JUNK",
		),
		sbm: &Sbm{
			pixelArray: SbmPixelArray{
				metaData: SbmPixelArrayMetaData{
					width: 10,
					header: SbmPixelArrayMetaDataHeader{
						width: SbmPixelArrayMetaDataHeaderData{
							bottomLeft:  9,
							bottomRight: 1,
						},
					},
				},
			},
		},
		sbmExpected: Sbm{
			pixelArray: SbmPixelArray{
				metaData: SbmPixelArrayMetaData{
					width: 10,
					header: SbmPixelArrayMetaDataHeader{
						width: SbmPixelArrayMetaDataHeaderData{
							bottomLeft:  9,
							bottomRight: 1,
						},
					},
				},
			},
		},
		errorTextExpected: ErrHeaderSyntax,
	}
	tests = append(tests, test)

	// Test #7. Negative.
	// Conflict with existing Height (e.g. taken from the Top header).
	// The old Value of the Object is not changed.
	test = Test{
		kind: TestKindMustBeExactError,
		data: []byte(
			"WIDTH 10 (9 + 1)\r\n" +
				"HEIGHT 100 (99 + 1)\r\n" +
				"JUNK",
		),
		sbm: &Sbm{
			pixelArray: SbmPixelArray{
				metaData: SbmPixelArrayMetaData{
					width:  10,
					height: 15,
					header: SbmPixelArrayMetaDataHeader{
						width: SbmPixelArrayMetaDataHeaderData{
							bottomLeft:  9,
							bottomRight: 1,
						},
					},
				},
			},
		},
		sbmExpected: Sbm{
			pixelArray: SbmPixelArray{
				metaData: SbmPixelArrayMetaData{
					width:  10,
					height: 15,
					header: SbmPixelArrayMetaDataHeader{
						width: SbmPixelArrayMetaDataHeaderData{
							bottomLeft:  9,
							bottomRight: 1,
						},
					},
				},
			},
		},
		errorTextExpected: ErrBottomHeaderMismatch,
	}
	tests = append(tests, test)

	// Test #8. Negative. Bad Area Header.
	test = Test{
		kind: TestKindMustBeExactError,
		data: []byte(
			"WIDTH 10 (9 + 1)\r\n" +
				"HEIGHT 15 (14 + 1)\r\n" +
				"QWERTY",
		),
		sbm: &Sbm{
			pixelArray: SbmPixelArray{
				metaData: SbmPixelArrayMetaData{
					width:  10,
					height: 15,
					header: SbmPixelArrayMetaDataHeader{
						width: SbmPixelArrayMetaDataHeaderData{
							bottomLeft:  9,
							bottomRight: 1,
						},
						height: SbmPixelArrayMetaDataHeaderData{
							bottomLeft:  14,
							bottomRight: 1,
						},
					},
				},
			},
		},
		sbmExpected: Sbm{
			pixelArray: SbmPixelArray{
				metaData: SbmPixelArrayMetaData{
					width:  10,
					height: 15,
					header: SbmPixelArrayMetaDataHeader{
						width: SbmPixelArrayMetaDataHeaderData{
							bottomLeft:  9,
							bottomRight: 1,
						},
						height: SbmPixelArrayMetaDataHeaderData{
							bottomLeft:  14,
							bottomRight: 1,
						},
					},
				},
			},
		},
		errorTextExpected: io.EOF.Error(),
	}
	tests = append(tests, test)

	// Test #9. Negative. Bad Area Header.
	test = Test{
		kind: TestKindMustBeExactError,
		data: []byte(
			"WIDTH 10 (9 + 1)\r\n" +
				"HEIGHT 15 (14 + 1)\r\n" +
				"AREA XYZ\r\n" +
				"QWERTY",
		),
		sbm: &Sbm{
			pixelArray: SbmPixelArray{
				metaData: SbmPixelArrayMetaData{
					width:  10,
					height: 15,
					header: SbmPixelArrayMetaDataHeader{
						width: SbmPixelArrayMetaDataHeaderData{
							bottomLeft:  9,
							bottomRight: 1,
						},
						height: SbmPixelArrayMetaDataHeaderData{
							bottomLeft:  14,
							bottomRight: 1,
						},
					},
				},
			},
		},
		sbmExpected: Sbm{
			pixelArray: SbmPixelArray{
				metaData: SbmPixelArrayMetaData{
					width:  10,
					height: 15,
					header: SbmPixelArrayMetaDataHeader{
						width: SbmPixelArrayMetaDataHeaderData{
							bottomLeft:  9,
							bottomRight: 1,
						},
						height: SbmPixelArrayMetaDataHeaderData{
							bottomLeft:  14,
							bottomRight: 1,
						},
					},
				},
			},
		},
		errorTextExpected: ErrHeaderSyntax,
	}
	tests = append(tests, test)

	// Test #10. Negative.
	// Conflict with existing Area (e.g. taken from the Top header).
	// The old Value of the Object is not changed.
	test = Test{
		kind: TestKindMustBeExactError,
		data: []byte(
			"WIDTH 10 (9 + 1)\r\n" +
				"HEIGHT 15 (14 + 1)\r\n" +
				"AREA 999 (909 + 90)\r\n" +
				"QWERTY",
		),
		sbm: &Sbm{
			pixelArray: SbmPixelArray{
				metaData: SbmPixelArrayMetaData{
					width:  10,
					height: 15,
					area:   150,
					header: SbmPixelArrayMetaDataHeader{
						width: SbmPixelArrayMetaDataHeaderData{
							bottomLeft:  9,
							bottomRight: 1,
						},
						height: SbmPixelArrayMetaDataHeaderData{
							bottomLeft:  14,
							bottomRight: 1,
						},
					},
				},
			},
		},
		sbmExpected: Sbm{
			pixelArray: SbmPixelArray{
				metaData: SbmPixelArrayMetaData{
					width:  10,
					height: 15,
					area:   150,
					header: SbmPixelArrayMetaDataHeader{
						width: SbmPixelArrayMetaDataHeaderData{
							bottomLeft:  9,
							bottomRight: 1,
						},
						height: SbmPixelArrayMetaDataHeaderData{
							bottomLeft:  14,
							bottomRight: 1,
						},
					},
				},
			},
		},
		errorTextExpected: ErrBottomHeaderMismatch,
	}
	tests = append(tests, test)

	// Test #11. Negative.
	// Conflict with real Size calculated from Width and Height.
	test = Test{
		kind: TestKindMustBeExactError,
		data: []byte(
			"WIDTH 10 (9 + 1)\r\n" +
				"HEIGHT 15 (14 + 1)\r\n" +
				"AREA 9999 (9009 + 990)\r\n" +
				"QWERTY",
		),
		sbm: &Sbm{
			pixelArray: SbmPixelArray{
				metaData: SbmPixelArrayMetaData{
					width:  10,
					height: 15,
					area:   9999,
					header: SbmPixelArrayMetaDataHeader{
						width: SbmPixelArrayMetaDataHeaderData{
							bottomLeft:  9,
							bottomRight: 1,
						},
						height: SbmPixelArrayMetaDataHeaderData{
							bottomLeft:  14,
							bottomRight: 1,
						},
					},
				},
			},
		},
		sbmExpected: Sbm{
			pixelArray: SbmPixelArray{
				metaData: SbmPixelArrayMetaData{
					width:  10,
					height: 15,
					area:   9999,
					header: SbmPixelArrayMetaDataHeader{
						width: SbmPixelArrayMetaDataHeaderData{
							bottomLeft:  9,
							bottomRight: 1,
						},
						height: SbmPixelArrayMetaDataHeaderData{
							bottomLeft:  14,
							bottomRight: 1,
						},
					},
				},
			},
		},
		errorTextExpected: ErrAreaMismatch,
	}
	tests = append(tests, test)

	// Run the Tests.
	for testIdx, test = range tests {

		// Run the Action.
		reader = bytes.NewReader(test.data)
		lineReader = rdr.NewReader(reader)
		err = test.sbm.readBottomHeaders(lineReader)

		// Check.
		switch test.kind {

		case TestKindMustBeNoError:
			tst.MustBeNoError(err)

		case TestKindMustBeAnyError:
			tst.MustBeAnError(err)

		case TestKindMustBeExactError:
			tst.MustBeAnError(err)
			tst.MustBeEqual(err.Error(), test.errorTextExpected)

		default:
			t.FailNow()
		}

		tst.MustBeDifferent(test.sbm, nil)
		tst.MustBeEqual(*(test.sbm), test.sbmExpected)

		fmt.Printf("[%v]", testIdx+1)
	}
	fmt.Println()
}

func Test_readSeparator(t *testing.T) {

	var err error
	var lineReader *rdr.Reader
	var reader io.Reader
	var tst *tester.Test

	tst = tester.New(t)

	// Test #1. Positive.
	reader = bytes.NewReader([]byte(NL))
	lineReader = rdr.NewReader(reader)
	err = readSeparator(lineReader)
	tst.MustBeNoError(err)

	// Test #2. Negative.
	reader = bytes.NewReader([]byte("ABC"))
	lineReader = rdr.NewReader(reader)
	err = readSeparator(lineReader)
	tst.MustBeAnError(err)

	// Test #3. Negative.
	reader = bytes.NewReader([]byte{})
	lineReader = rdr.NewReader(reader)
	err = readSeparator(lineReader)
	tst.MustBeAnError(err)
}
