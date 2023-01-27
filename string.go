// string.go.

package sbm

import (
	"errors"
)

// Errors.
const (
	ErrHeaderSize   = "header is too short"
	ErrHeaderEnding = "header ending syntax error"
)

func removeCRLF(
	rawHeader []byte,
) ([]byte, error) {

	var err error
	var idxLast int

	// Check Header Size.
	if len(rawHeader) < 2 {
		err = errors.New(ErrHeaderSize)
		return []byte{}, err
	}

	// Check Header Ending.
	idxLast = len(rawHeader) - 1
	if (rawHeader[idxLast-1] != '\r') ||
		(rawHeader[idxLast] != '\n') {
		err = errors.New(ErrHeaderEnding)
		return []byte{}, err
	}

	return rawHeader[0 : idxLast-1], nil
}
