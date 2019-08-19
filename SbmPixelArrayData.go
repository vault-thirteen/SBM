package sbm

import (
	"github.com/vault-thirteen/bit"
)

// SBM Internal Data Model: Pixel Array Data.
type SbmPixelArrayData struct {
	bytes []byte
	bits  []bit.Bit
}
