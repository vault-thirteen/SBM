package sbm

import "github.com/vault-thirteen/auxie/errors"

// Errors.
const (
	Err_HeaderSize   = "header is too short"
	Err_HeaderEnding = "header ending syntax error"
)

func trimHeader(rawHeader []byte) (trimmedHeader []byte, err error) {
	// Check header's size.
	if len(rawHeader) < 2 {
		return trimmedHeader, errors.ErrorFromText(Err_HeaderSize)
	}

	// Check header's ending.
	idxLast := len(rawHeader) - 1
	if (rawHeader[idxLast-1] != CR) ||
		(rawHeader[idxLast] != LF) {
		return trimmedHeader, errors.ErrorFromText(Err_HeaderEnding)
	}

	return rawHeader[0 : idxLast-1], nil
}
