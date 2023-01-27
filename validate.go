// validate.go.

package sbm

import (
	"errors"
)

// Errors.
const (
	ErrVersion = "version error"
)

func validateFormat(
	headerFormat HeaderDataVersion,
) (err error) {

	if headerFormat.version != SbmFormatVersion1 {
		err = errors.New(ErrVersion)
		return
	}

	return
}
