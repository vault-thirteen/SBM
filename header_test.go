package sbm

import (
	"testing"

	"github.com/vault-thirteen/auxie/tester"
)

func Test_const(t *testing.T) {
	var tst = tester.New(t)

	// Test #1.
	tst.MustBeEqual(MimeType, "image/x-portable-bitmap")
}
