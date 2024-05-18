// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package action

import (
	"context"
	"io"
)

// CloseIt is an Action that can close any closable [Data].
func CloseIt[T io.Closer](_ context.Context, c T) error {
	return c.Close()
}
