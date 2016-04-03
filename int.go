// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fastprinter

import (
	"fmt"
	"io"
	"sync"
)

func PrintUintBase(w io.Writer, i uint64, base int) (int, error) {
	return formatBits(w, i, base, false)
}

func PrintUint(w io.Writer, i uint64) (int, error) {
	return formatBits(w, i, 10, false)
}

func PrintIntBase(w io.Writer, i int64, base int) (int, error) {
	return formatBits(w, uint64(i), base, i < 0)
}

func PrintInt(w io.Writer, i int64) (int, error) {
	return formatBits(w, uint64(i), 10, i < 0)
}

const (
	digits = "0123456789abcdefghijklmnopqrstuvwxyz"
)

var shifts = [len(digits) + 1]uint{
	1 << 1: 1,
	1 << 2: 2,
	1 << 3: 3,
	1 << 4: 4,
	1 << 5: 5,
}

const formatbitsArrayLen = 64 + 1 // +1 for sign of 64bit value in base 2
type formatbitsBytes struct {
	bytes []byte
}

var formatbitsBytesPool = sync.Pool{
	New: func() interface{} {
		return &formatbitsBytes{make([]byte, formatbitsArrayLen, formatbitsArrayLen)}
	},
}

// formatBits computes the string representation of u in the given base.
// If neg is set, u is treated as negative int64 value. If append_ is
// set, the string is appended to dst and the resulting byte slice is
// returned as the first result value; otherwise the string is returned
// as the second result value.
//
func formatBits(dst io.Writer, u uint64, base int, neg bool) (int, error) {

	if base < 2 || base > len(digits) {
		return 0, fmt.Errorf("fastprinter: illegal base")
	}
	// 2 <= base && base <= len(digits)
	var a = formatbitsBytesPool.Get().(*formatbitsBytes)

	i := formatbitsArrayLen

	if neg {
		u = -u
	}

	// convert bits
	if base == 10 {
		// common case: use constants for / because
		// the compiler can optimize it into a multiply+shift

		if ^uintptr(0)>>32 == 0 {
			for u > uint64(^uintptr(0)) {
				q := u / 1e9
				us := uintptr(u - q*1e9) // us % 1e9 fits into a uintptr
				for j := 9; j > 0; j-- {
					i--
					qs := us / 10
					a.bytes[i] = byte(us - qs*10 + '0')
					us = qs
				}
				u = q
			}
		}

		// u guaranteed to fit into a uintptr
		us := uintptr(u)
		for us >= 10 {
			i--
			q := us / 10
			a.bytes[i] = byte(us - q*10 + '0')
			us = q
		}
		// u < 10
		i--
		a.bytes[i] = byte(us + '0')

	} else if s := shifts[base]; s > 0 {
		// base is power of 2: use shifts and masks instead of / and %
		b := uint64(base)
		m := uintptr(b) - 1 // == 1<<s - 1
		for u >= b {
			i--
			a.bytes[i] = digits[uintptr(u)&m]
			u >>= s
		}
		// u < base
		i--
		a.bytes[i] = digits[uintptr(u)]

	} else {
		// general case
		b := uint64(base)
		for u >= b {
			i--
			q := u / b
			a.bytes[i] = digits[uintptr(u-q*b)]
			u = q
		}
		// u < base
		i--
		a.bytes[i] = digits[uintptr(u)]
	}

	// add sign, if any
	if neg {
		i--
		a.bytes[i] = '-'
	}

	counter, err := dst.Write(a.bytes[i:])
	formatbitsBytesPool.Put(a)

	return counter, err
}
