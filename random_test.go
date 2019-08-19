package sbm

import (
	"testing"

	"github.com/vault-thirteen/tester"
)

func Test_createRandomValuePair(t *testing.T) {

	var err error
	var tst *tester.Test
	var valueLeft uint
	var valueRight uint
	var valueSum uint

	tst = tester.New(t)

	valueSum = 100
	iMax := 1000
	for i := 1; i <= iMax; i++ {
		valueLeft, valueRight, err = createRandomValuePair(valueSum)
		tst.MustBeNoError(err)
		tst.MustBeEqual(valueLeft+valueRight, valueSum)
	}
}
