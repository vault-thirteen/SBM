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

func Test_validateFormat(t *testing.T) {

	const VersionNone = 0

	var err error
	var headerFormat HeaderDataVersion
	var tst *tester.Test

	tst = tester.New(t)

	// Test #1. Positive.
	headerFormat = HeaderDataVersion{
		version: SbmFormatVersion1,
	}
	err = validateFormat(headerFormat)
	tst.MustBeNoError(err)

	// Test #1. Negative.
	headerFormat = HeaderDataVersion{
		version: VersionNone,
	}
	err = validateFormat(headerFormat)
	tst.MustBeAnError(err)
	tst.MustBeEqual(err.Error(), ErrVersion)
}
