package sbm

import (
	"fmt"
	"testing"

	"github.com/vault-thirteen/auxie/tester"
)

func Test_parseHeaderFormat(t *testing.T) {

	var err error
	var rawHeader []byte
	var sbm *Sbm
	var tst *tester.Test

	tst = tester.New(t)

	// Test #1. Good Header.
	sbm = &Sbm{}
	rawHeader = []byte("SBM (SIMPLE BIT MAP)" + NL)
	err = sbm.parseHeaderFormat(rawHeader)

	// Check the Result.
	tst.MustBeNoError(err)

	// Test #2. Bad Header.
	sbm = &Sbm{}
	rawHeader = []byte("ABC")
	err = sbm.parseHeaderFormat(rawHeader)

	// Check the Result.
	tst.MustBeAnError(err)
}

func Test_parseHeaderVersion(t *testing.T) {

	const VersionNone = 0

	// Test Type (Kind).
	const (
		TestKindMustBeNoError    = 1
		TestKindMustBeAnyError   = 2
		TestKindMustBeExactError = 3
	)

	type Test struct {
		kind               byte
		rawHeader          []byte
		errorTextExpected  string
		headerDataExpected HeaderDataVersion
	}

	var err error
	var headerData HeaderDataVersion
	var test Test
	var testIdx int
	var tests []Test
	var tst *tester.Test

	tst = tester.New(t)
	tests = make([]Test, 0)

	// Test #1. Good Header.
	test = Test{
		kind:              TestKindMustBeNoError,
		rawHeader:         []byte("VERSION 1" + NL),
		errorTextExpected: "",
		headerDataExpected: HeaderDataVersion{
			SbmFormatVersion1,
		},
	}
	tests = append(tests, test)

	// Test #2. Bad Header: No CR+LF End.
	test = Test{
		kind:              TestKindMustBeExactError,
		rawHeader:         []byte("VERSION X" + "\r"),
		errorTextExpected: ErrHeaderEnding,
		headerDataExpected: HeaderDataVersion{
			VersionNone,
		},
	}
	tests = append(tests, test)

	// Test #3. Bad Header: Not two Parts.
	test = Test{
		kind:              TestKindMustBeExactError,
		rawHeader:         []byte("VERSION X Y" + NL),
		errorTextExpected: ErrHeaderSyntax,
		headerDataExpected: HeaderDataVersion{
			VersionNone,
		},
	}
	tests = append(tests, test)

	// Test #4. Bad Header: Header Name.
	test = Test{
		kind:              TestKindMustBeExactError,
		rawHeader:         []byte("VERSIOM X" + NL),
		errorTextExpected: fmt.Sprintf(ErrfHeaderUnexpected, "VERSIOM"),
		headerDataExpected: HeaderDataVersion{
			VersionNone,
		},
	}
	tests = append(tests, test)

	// Test #5. Bad Header: Version is not numeric.
	test = Test{
		kind:              TestKindMustBeAnyError,
		rawHeader:         []byte("VERSION X" + NL),
		errorTextExpected: "",
		headerDataExpected: HeaderDataVersion{
			VersionNone,
		},
	}
	tests = append(tests, test)

	// Test #6. Bad Header: Version Overflow.
	test = Test{
		kind:              TestKindMustBeAnyError,
		rawHeader:         []byte("VERSION 256" + NL),
		errorTextExpected: ErrOverflow,
		headerDataExpected: HeaderDataVersion{
			VersionNone,
		},
	}
	tests = append(tests, test)

	// Run the Tests.
	for testIdx, test = range tests {
		headerData, err = parseHeaderVersion(test.rawHeader)
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

		tst.MustBeEqual(headerData, test.headerDataExpected)
		fmt.Printf("[%v]", testIdx+1)
	}
	fmt.Println()
}

func Test_parseHeaderSize(t *testing.T) {

	const HeaderName = "SIZE"

	// Test Type (Kind).
	const (
		TestKindMustBeNoError    = 1
		TestKindMustBeAnyError   = 2
		TestKindMustBeExactError = 3
	)

	type Test struct {
		kind               byte
		rawHeader          []byte
		errorTextExpected  string
		headerDataExpected HeaderDataSize
	}

	var err error
	var headerData HeaderDataSize
	var test Test
	var testIdx int
	var tests []Test
	var tst *tester.Test

	tst = tester.New(t)
	tests = make([]Test, 0)

	// Test #1. Good Header.
	test = Test{
		kind:              TestKindMustBeNoError,
		rawHeader:         []byte("SIZE 123 (100 + 23)" + NL),
		errorTextExpected: "",
		headerDataExpected: HeaderDataSize{
			sizeFixed:       123,
			sizeRandomLeft:  100,
			sizeRandomRight: 23,
		},
	}
	tests = append(tests, test)

	// Test #2. Bad Header: No CR+LF End.
	test = Test{
		kind:               TestKindMustBeExactError,
		rawHeader:          []byte("JUNKY TOWN" + "\r"),
		errorTextExpected:  ErrHeaderEnding,
		headerDataExpected: HeaderDataSize{},
	}
	tests = append(tests, test)

	// Test #3. Bad Header: Parts Count is wrong.
	test = Test{
		kind:               TestKindMustBeExactError,
		rawHeader:          []byte("SIZE 123 (100 + 23) QUAKE" + NL),
		errorTextExpected:  ErrHeaderSyntax,
		headerDataExpected: HeaderDataSize{},
	}
	tests = append(tests, test)

	// Test #4. Bad Header: Header Name.
	test = Test{
		kind:               TestKindMustBeExactError,
		rawHeader:          []byte("ALPHA 123 (100 + 23)" + NL),
		errorTextExpected:  fmt.Sprintf(ErrfHeaderUnexpected, "ALPHA"),
		headerDataExpected: HeaderDataSize{},
	}
	tests = append(tests, test)

	// Test #5. Bad Header: Fixed Size is not numeric.
	test = Test{
		kind:               TestKindMustBeAnyError,
		rawHeader:          []byte("SIZE XYZ (100 + 23)" + NL),
		errorTextExpected:  "",
		headerDataExpected: HeaderDataSize{},
	}
	tests = append(tests, test)

	// Test #6. Bad Header: Left Bracket is absent.
	test = Test{
		kind:               TestKindMustBeExactError,
		rawHeader:          []byte("SIZE 123 |100 + 23)" + NL),
		errorTextExpected:  ErrHeaderSyntax,
		headerDataExpected: HeaderDataSize{},
	}
	tests = append(tests, test)

	// Test #7. Bad Header: Left Size is not numeric.
	test = Test{
		kind:               TestKindMustBeAnyError,
		rawHeader:          []byte("SIZE 123 (XYZ + 23)" + NL),
		errorTextExpected:  "",
		headerDataExpected: HeaderDataSize{},
	}
	tests = append(tests, test)

	// Test #8. Bad Header: No Plus Sign.
	test = Test{
		kind:               TestKindMustBeExactError,
		rawHeader:          []byte("SIZE 123 (100 * 23)" + NL),
		errorTextExpected:  ErrHeaderSyntax,
		headerDataExpected: HeaderDataSize{},
	}
	tests = append(tests, test)

	// Test #9. Bad Header: Right Bracket is absent.
	test = Test{
		kind:               TestKindMustBeExactError,
		rawHeader:          []byte("SIZE 123 (100 + 23|" + NL),
		errorTextExpected:  ErrHeaderSyntax,
		headerDataExpected: HeaderDataSize{},
	}
	tests = append(tests, test)

	// Test #10. Bad Header: Right Size is not numeric.
	test = Test{
		kind:               TestKindMustBeAnyError,
		rawHeader:          []byte("SIZE 123 (100 + ABC)" + NL),
		errorTextExpected:  "",
		headerDataExpected: HeaderDataSize{},
	}
	tests = append(tests, test)

	// Test #11. Bad Sum.
	test = Test{
		kind:               TestKindMustBeExactError,
		rawHeader:          []byte("SIZE 123 (100 + 24)" + NL),
		errorTextExpected:  ErrIntegrity,
		headerDataExpected: HeaderDataSize{},
	}
	tests = append(tests, test)

	// Run the Tests.
	for testIdx, test = range tests {
		headerData, err = parseHeaderSize(test.rawHeader, HeaderName)
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

		tst.MustBeEqual(headerData, test.headerDataExpected)
		fmt.Printf("[%v]", testIdx+1)
	}
	fmt.Println()
}

func Test_parseHeaderWidth(t *testing.T) {

	var err error
	var headerData HeaderDataSize
	var headerDataExpected HeaderDataSize
	var rawHeader []byte
	var tst *tester.Test

	tst = tester.New(t)

	// Test #1. Good Header.
	rawHeader = []byte("WIDTH 123 (100 + 23)" + NL)
	headerDataExpected = HeaderDataSize{
		sizeFixed:       123,
		sizeRandomLeft:  100,
		sizeRandomRight: 23,
	}
	headerData, err = parseHeaderWidth(rawHeader)
	tst.MustBeNoError(err)
	tst.MustBeEqual(headerData, headerDataExpected)

	// Test #2. Bad Header.
	rawHeader = []byte("XYZ 123 (100 + 23)" + NL)
	headerDataExpected = HeaderDataSize{}
	headerData, err = parseHeaderWidth(rawHeader)
	tst.MustBeAnError(err)
	tst.MustBeEqual(headerData, headerDataExpected)
}

func Test_parseHeaderHeight(t *testing.T) {

	var err error
	var headerData HeaderDataSize
	var headerDataExpected HeaderDataSize
	var rawHeader []byte
	var tst *tester.Test

	tst = tester.New(t)

	// Test #1. Good Header.
	rawHeader = []byte("HEIGHT 123 (100 + 23)" + NL)
	headerDataExpected = HeaderDataSize{
		sizeFixed:       123,
		sizeRandomLeft:  100,
		sizeRandomRight: 23,
	}
	headerData, err = parseHeaderHeight(rawHeader)
	tst.MustBeNoError(err)
	tst.MustBeEqual(headerData, headerDataExpected)

	// Test #2. Bad Header.
	rawHeader = []byte("XYZ 123 (100 + 23)" + NL)
	headerDataExpected = HeaderDataSize{}
	headerData, err = parseHeaderHeight(rawHeader)
	tst.MustBeAnError(err)
	tst.MustBeEqual(headerData, headerDataExpected)
}

func Test_parseHeaderArea(t *testing.T) {

	var err error
	var headerData HeaderDataSize
	var headerDataExpected HeaderDataSize
	var rawHeader []byte
	var tst *tester.Test

	tst = tester.New(t)

	// Test #1. Good Header.
	rawHeader = []byte("AREA 123 (100 + 23)" + NL)
	headerDataExpected = HeaderDataSize{
		sizeFixed:       123,
		sizeRandomLeft:  100,
		sizeRandomRight: 23,
	}
	headerData, err = parseHeaderArea(rawHeader)
	tst.MustBeNoError(err)
	tst.MustBeEqual(headerData, headerDataExpected)

	// Test #2. Bad Header.
	rawHeader = []byte("XYZ 123 (100 + 23)" + NL)
	headerDataExpected = HeaderDataSize{}
	headerData, err = parseHeaderArea(rawHeader)
	tst.MustBeAnError(err)
	tst.MustBeEqual(headerData, headerDataExpected)
}
