package sbm

import (
	"github.com/vault-thirteen/auxie/random"
)

// createRandomValuePair creates a pair of random values, which have the sum
// equal to the sum specified.
func createRandomValuePair(valueSum uint) (valueLeft uint, valueRight uint, err error) {
	valueLeft, err = random.Uint(0, valueSum)
	if err != nil {
		return valueLeft, valueRight, err
	}

	valueRight = valueSum - valueLeft
	return valueLeft, valueRight, nil
}
