package sbm

import "github.com/vault-thirteen/auxie/errors"

const (
	ErrVersion = "version error"
)

func validateFormat(headerFormat HeaderDataVersion) (err error) {
	if headerFormat.version != SbmFormatVersion1 {
		return errors.ErrorFromText(ErrVersion)
	}

	return nil
}
