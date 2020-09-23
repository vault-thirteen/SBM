// string.go.

////////////////////////////////////////////////////////////////////////////////
//
// Copyright © 2019..2020 by Vault Thirteen.
//
// All rights reserved. No part of this publication may be reproduced,
// distributed, or transmitted in any form or by any means, including
// photocopying, recording, or other electronic or mechanical methods,
// without the prior written permission of the publisher, except in the case
// of brief quotations embodied in critical reviews and certain other
// noncommercial uses permitted by copyright law. For permission requests,
// write to the publisher, addressed “Copyright Protected Material” at the
// address below.
//
////////////////////////////////////////////////////////////////////////////////
//
// Web Site Address:	https://github.com/vault-thirteen.
//
////////////////////////////////////////////////////////////////////////////////

package sbm

import (
	"errors"
)

// Errors.
const (
	ErrHeaderSize   = "Header is too short"
	ErrHeaderEnding = "Header Ending Syntax Error"
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
