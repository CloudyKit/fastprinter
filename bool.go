// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fastprinter

import "io"

var (
	_true  = ([]byte)("true")
	_false = ([]byte)("false")
)

func PrintBool(w io.Writer, b bool) (int, error) {
	if b {
		return w.Write(_true)
	}
	return w.Write(_false)
}
