package sbm

import (
	"github.com/vault-thirteen/bit"
)

// Returns the Format.
func (sbm Sbm) GetFormat() SbmFormat {
	return sbm.format
}

func (sbm Sbm) GetArrayBytes() []byte {
	return sbm.pixelArray.data.bytes
}

func (sbm Sbm) GetArrayBits() []bit.Bit {
	return sbm.pixelArray.data.bits
}

func (sbm Sbm) GetArrayWidth() uint {
	return sbm.pixelArray.metaData.width
}

func (sbm Sbm) GetArrayHeight() uint {
	return sbm.pixelArray.metaData.height
}

func (sbm Sbm) GetArrayArea() uint {
	return sbm.pixelArray.metaData.area
}
