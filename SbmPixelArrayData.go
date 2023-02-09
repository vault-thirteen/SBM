// SbmPixelArrayData.go.

package sbm

import (
	"github.com/vault-thirteen/auxie/bit"
)

// SBM Internal Data Model: Pixel Array Data.
type SbmPixelArrayData struct {
	bytes []byte
	bits  []bit.Bit
}
