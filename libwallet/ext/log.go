// Copyright (c) 2013-2021 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package ext

import "github.com/decred/slog"

var log = slog.Disabled

// UseLogger uses a specified Logger to output package logging info.
func UseLogger(logger slog.Logger) {
	log = logger
}
