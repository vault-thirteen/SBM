////////////////////////////////////////////////////////////////////////////////
//
// Copyright © 2019 by Vault Thirteen.
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
	"testing"

	"github.com/vault-thirteen/tester"
)

func Test_removeCRLF(t *testing.T) {

	var err error
	var rawHeader []byte
	var result []byte
	var resultExpected []byte
	var tst *tester.Test

	tst = tester.New(t)

	// Test #1. Normal Test.
	rawHeader = append(
		[]byte("ABC"),
		'\r', '\n',
	)
	resultExpected = []byte("ABC")
	result, err = removeCRLF(rawHeader)
	tst.MustBeNoError(err)
	tst.MustBeEqual(result, resultExpected)

	// Test #2. The last Symbol is not LF.
	rawHeader = append(
		[]byte("ABC"),
		'\r', '\r',
	)
	resultExpected = []byte{}
	result, err = removeCRLF(rawHeader)
	tst.MustBeAnError(err)
	tst.MustBeEqual(err.Error(), ErrHeaderEnding)
	tst.MustBeEqual(result, resultExpected)

	// Test #3. The pre-last Symbol is not CR.
	rawHeader = append(
		[]byte("ABC"),
		'\n', '\n',
	)
	resultExpected = []byte{}
	result, err = removeCRLF(rawHeader)
	tst.MustBeAnError(err)
	tst.MustBeEqual(err.Error(), ErrHeaderEnding)
	tst.MustBeEqual(result, resultExpected)

	// Test #4. The Header is too short.
	rawHeader = []byte{'\r'}
	resultExpected = []byte{}
	result, err = removeCRLF(rawHeader)
	tst.MustBeAnError(err)
	tst.MustBeEqual(err.Error(), ErrHeaderSize)
	tst.MustBeEqual(result, resultExpected)

	// Test #5. Empty Header.
	rawHeader = []byte{'\r', '\n'}
	resultExpected = []byte{}
	result, err = removeCRLF(rawHeader)
	tst.MustBeNoError(err)
	tst.MustBeEqual(result, resultExpected)
}
