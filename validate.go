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
