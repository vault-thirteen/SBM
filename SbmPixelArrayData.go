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

// SBM Internal Data Model: Pixel Array Data.
type SbmPixelArrayData struct {
	bytes []byte
	bits  []bit.Bit
}
