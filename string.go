package sbm

import (
	"errors"
)

// Errors.
const (
	ErrHeaderSize   = "header is too short"
	ErrHeaderEnding = "header ending syntax error"
)

func trimHeader(rawHeader []byte) (trimmedHeader []byte, err error) {
	// Check header's size.
	if len(rawHeader) < 2 {
		return trimmedHeader, errors.New(ErrHeaderSize)
	}

	// Check header's ending.
	idxLast := len(rawHeader) - 1
	if (rawHeader[idxLast-1] != CR) ||
		(rawHeader[idxLast] != LF) {
		return trimmedHeader, errors.New(ErrHeaderEnding)
	}

	return rawHeader[0 : idxLast-1], nil
}
