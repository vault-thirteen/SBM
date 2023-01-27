// validate.go.

package sbm

import (
	"errors"
)

// Errors.
const (
	ErrVersion = "Version Error"
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
