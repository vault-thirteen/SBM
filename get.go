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
	"github.com/vault-thirteen/bit"
)

// Returns the Format.
func (sbm Sbm) GetFormat() SbmFormat {
	return sbm.format
}

func (sbm Sbm) GetArrayBytes() []byte {
	return sbm.pixelArray.data.bytes
}

func (sbm Sbm) GetArrayBits() []bit.Bit {
	return sbm.pixelArray.data.bits
}

func (sbm Sbm) GetArrayWidth() uint {
	return sbm.pixelArray.metaData.width
}

func (sbm Sbm) GetArrayHeight() uint {
	return sbm.pixelArray.metaData.height
}

func (sbm Sbm) GetArrayArea() uint {
	return sbm.pixelArray.metaData.area
}
