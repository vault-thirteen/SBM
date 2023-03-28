package sbm

import (
	"errors"
)

const (
	ErrVersion = "version error"
)

func validateFormat(headerFormat HeaderDataVersion) (err error) {
	if headerFormat.version != SbmFormatVersion1 {
		return errors.New(ErrVersion)
	}

	return nil
}
