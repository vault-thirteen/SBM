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
		kind: TestKindMustBeAnyError,
		data: []byte(
			"SBM (SIMPLE HACKED MAP)\r\n" +
				"JUNK",
		),
		sbm:               new(Sbm),
		sbmExpected:       Sbm{},
		errorTextExpected: "",
	}
	tests = append(tests, test)

	// Test #3. Negative. Bad Format Header.
	test = Test{
		kind: TestKindMustBeAnyError,
		data: []byte(
			"SOME JUNK IS HERE",
		),
		sbm:               new(Sbm),
		sbmExpected:       Sbm{},
		errorTextExpected: "",
	}
	tests = append(tests, test)

	// Test #4. Negative. Bad Version Header.
	test = Test{
		kind: TestKindMustBeAnyError,
		data: []byte(
			"SBM (SIMPLE BIT MAP)\r\n" +
				"VERSION JUNK",
		),
		sbm:               new(Sbm),
		sbmExpected:       Sbm{},
		errorTextExpected: "",
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
		kind: TestKindMustBeAnyError,
		data: []byte(
			"SBM (SIMPLE BIT MAP)\r\n" +
				"VERSION 123\r\n" +
				"JUNK",
		),
		sbm:               new(Sbm),
		sbmExpected:       Sbm{},
		errorTextExpected: "",
	}
	tests = append(tests, test)

	// Test #7. Negative. Bad Width Header.
	test = Test{
		kind: TestKindMustBeAnyError,
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
		errorTextExpected: "",
	}
	tests = append(tests, test)

	// Test #8. Negative. Bad Width Header.
	test = Test{
		kind: TestKindMustBeAnyError,
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
		errorTextExpected: "",
	}
	tests = append(tests, test)

	// Test #9. Negative. Bad Area.
	test = Test{
		kind: TestKindMustBeAnyError,
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
		errorTextExpected: "",
	}
	tests = append(tests, test)

	// Run the Tests.
	for testIdx, test = range tests {

		// Run the Action.
		reader = bytes.NewReader(test.data)
		lineReader = rdr.New(reader)
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
	lineReader = rdr.New(reader)
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
	lineReader = rdr.New(reader)
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
	lineReader = rdr.New(reader)
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
	lineReader = rdr.New(reader)
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
	lineReader = rdr.New(reader)
	err = sbm.readArrayData(lineReader)
	tst.MustBeAnError(err)
	tst.MustBeEqual(err.Error(), io.ErrUnexpectedEOF.Error())
	tst.MustBeEqual(*sbm, sbmExpected)
}

func Test_readBottomHeaders(t *testing.T) {

	var data []byte
	var err error
	var lineReader *rdr.Reader
	var reader io.Reader
	var sbm *Sbm
	var sbmExpected Sbm
	var tst *tester.Test

	tst = tester.New(t)

	// Test #1. Positive.
	data = []byte(
		"WIDTH 123 (100 + 23)" + NL +
			"HEIGHT 456 (450 + 6)" + NL +
			"AREA 56088 (56000 + 88)" + NL +
			"SOME JUNK IS HERE",
	)
	sbm = new(Sbm)
	sbm.pixelArray.metaData.width = 123
	sbm.pixelArray.metaData.height = 456
	sbm.pixelArray.metaData.area = 56088
	sbmExpected = Sbm{
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
	}
	reader = bytes.NewReader(data)
	lineReader = rdr.New(reader)
	err = sbm.readBottomHeaders(lineReader)
	tst.MustBeNoError(err)
	tst.MustBeEqual(*sbm, sbmExpected)

	// Test #2. Negative.
	// Conflict with existing Width (e.g. taken from the Top header).
	// The old Value of the Object is not changed.
	data = []byte(
		"WIDTH 123 (100 + 23)" + NL +
			"HEIGHT 456 (450 + 6)" + NL +
			"AREA 789 (780 + 9)" + NL +
			"SOME JUNK IS HERE",
	)
	sbm = new(Sbm)
	sbm.pixelArray.metaData.width = 999
	sbm.pixelArray.metaData.height = 456
	sbm.pixelArray.metaData.area = 789
	sbmExpected = Sbm{
		pixelArray: SbmPixelArray{
			data: SbmPixelArrayData{},
			metaData: SbmPixelArrayMetaData{
				width:  999,
				height: 456,
				area:   789,
			},
		},
	}
	reader = bytes.NewReader(data)
	lineReader = rdr.New(reader)
	err = sbm.readBottomHeaders(lineReader)
	tst.MustBeAnError(err)
	tst.MustBeEqual(err.Error(), ErrBottomHeaderMismatch)
	tst.MustBeEqual(*sbm, sbmExpected)

	// Test #3. Negative.
	// Conflict with existing Height (e.g. taken from the Top header).
	// The old Value of the Object has Changes only in the following Fields:
	//	*	Random Bottom Width Headers.
	data = []byte(
		"WIDTH 123 (100 + 23)" + NL +
			"HEIGHT 456 (450 + 6)" + NL +
			"AREA 789 (780 + 9)" + NL +
			"SOME JUNK IS HERE",
	)
	sbm = new(Sbm)
	sbm.pixelArray.metaData.width = 123
	sbm.pixelArray.metaData.height = 999
	sbm.pixelArray.metaData.area = 789
	sbmExpected = Sbm{
		pixelArray: SbmPixelArray{
			data: SbmPixelArrayData{},
			metaData: SbmPixelArrayMetaData{
				width:  123,
				height: 999,
				area:   789,
				header: SbmPixelArrayMetaDataHeader{
					width: SbmPixelArrayMetaDataHeaderData{
						bottomLeft:  100,
						bottomRight: 23,
					},
					height: SbmPixelArrayMetaDataHeaderData{},
					area:   SbmPixelArrayMetaDataHeaderData{},
				},
			},
		},
	}
	reader = bytes.NewReader(data)
	lineReader = rdr.New(reader)
	err = sbm.readBottomHeaders(lineReader)
	tst.MustBeAnError(err)
	tst.MustBeEqual(err.Error(), ErrBottomHeaderMismatch)
	tst.MustBeEqual(*sbm, sbmExpected)

	// Test #4. Negative.
	// Conflict with existing Area (e.g. taken from the Top header).
	// The old Value of the Object has Changes only in the following Fields:
	//	*	Random Bottom Width Headers;
	//	*	Random Bottom Height Headers.
	data = []byte(
		"WIDTH 123 (100 + 23)" + NL +
			"HEIGHT 456 (450 + 6)" + NL +
			"AREA 789 (780 + 9)" + NL +
			"SOME JUNK IS HERE",
	)
	sbm = new(Sbm)
	sbm.pixelArray.metaData.width = 123
	sbm.pixelArray.metaData.height = 456
	sbm.pixelArray.metaData.area = 999
	sbmExpected = Sbm{
		pixelArray: SbmPixelArray{
			data: SbmPixelArrayData{},
			metaData: SbmPixelArrayMetaData{
				width:  123,
				height: 456,
				area:   999,
				header: SbmPixelArrayMetaDataHeader{
					width: SbmPixelArrayMetaDataHeaderData{
						bottomLeft:  100,
						bottomRight: 23,
					},
					height: SbmPixelArrayMetaDataHeaderData{
						bottomLeft:  450,
						bottomRight: 6,
					},
					area: SbmPixelArrayMetaDataHeaderData{},
				},
			},
		},
	}
	reader = bytes.NewReader(data)
	lineReader = rdr.New(reader)
	err = sbm.readBottomHeaders(lineReader)
	tst.MustBeAnError(err)
	tst.MustBeEqual(err.Error(), ErrBottomHeaderMismatch)
	tst.MustBeEqual(*sbm, sbmExpected)

	// Test #5. Negative.
	// Area is not equal to Width multiplied by Height.
	data = []byte(
		"WIDTH 123 (100 + 23)" + NL +
			"HEIGHT 456 (450 + 6)" + NL +
			"AREA 56089 (56000 + 89)" + NL +
			"SOME JUNK IS HERE",
	)
	sbm = new(Sbm)
	sbm.pixelArray.metaData.width = 123
	sbm.pixelArray.metaData.height = 456
	sbm.pixelArray.metaData.area = 56089
	sbmExpected = Sbm{
		pixelArray: SbmPixelArray{
			data: SbmPixelArrayData{},
			metaData: SbmPixelArrayMetaData{
				width:  123,
				height: 456,
				area:   56089,
				header: SbmPixelArrayMetaDataHeader{
					width: SbmPixelArrayMetaDataHeaderData{
						bottomLeft:  100,
						bottomRight: 23,
					},
					height: SbmPixelArrayMetaDataHeaderData{
						bottomLeft:  450,
						bottomRight: 6,
					},
					area: SbmPixelArrayMetaDataHeaderData{},
				},
			},
		},
	}
	reader = bytes.NewReader(data)
	lineReader = rdr.New(reader)
	err = sbm.readBottomHeaders(lineReader)
	tst.MustBeAnError(err)
	tst.MustBeEqual(err.Error(), ErrAreaMismatch)
	tst.MustBeEqual(*sbm, sbmExpected)
}

func Test_readSeparator(t *testing.T) {

	var err error
	var lineReader *rdr.Reader
	var reader io.Reader
	var tst *tester.Test

	tst = tester.New(t)

	// Test #1. Positive.
	reader = bytes.NewReader([]byte(NL))
	lineReader = rdr.New(reader)
	err = readSeparator(lineReader)
	tst.MustBeNoError(err)

	// Test #2. Negative.
	reader = bytes.NewReader([]byte("ABC"))
	lineReader = rdr.New(reader)
	err = readSeparator(lineReader)
	tst.MustBeAnError(err)

	// Test #3. Negative.
	reader = bytes.NewReader([]byte{})
	lineReader = rdr.New(reader)
	err = readSeparator(lineReader)
	tst.MustBeAnError(err)
}
