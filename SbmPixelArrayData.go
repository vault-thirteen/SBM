package sbm

import (
	"github.com/vault-thirteen/auxie/bit"
)

// SbmPixelArrayData contains pixel array data.
type SbmPixelArrayData struct {
	bits  []bit.Bit
	bytes []byte
}
