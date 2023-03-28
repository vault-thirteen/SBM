package sbm

import (
	"testing"

	"github.com/vault-thirteen/tester"
)

func Test_validateFormat(t *testing.T) {
	const VersionNone = 0

	var tst = tester.New(t)
	var err error
	var headerFormat HeaderDataVersion

	// Test #1. Positive.
	headerFormat = HeaderDataVersion{
		version: SbmFormatVersion1,
	}
	err = validateFormat(headerFormat)
	tst.MustBeNoError(err)

	// Test #1. Negative.
	headerFormat = HeaderDataVersion{
		version: VersionNone,
	}
	err = validateFormat(headerFormat)
	tst.MustBeAnError(err)
	tst.MustBeEqual(err.Error(), ErrVersion)
}
