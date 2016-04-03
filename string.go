package fastprinter

import (
	"io"
	"sync"
	"unsafe"
)

const defaultBufSize = 4096

type printStringSafeBytes struct {
	bytes [defaultBufSize]byte
}

var printStringSafeBytesPool = sync.Pool{
	New: func() interface{} {
		return new(printStringSafeBytes)
	},
}

func PrintStringSafe(ww io.Writer, st string) (n int, e error) {
	a := printStringSafeBytesPool.Get().(*printStringSafeBytes)
	numI := len(st) / defaultBufSize
	for i := 0; i < numI; i++ {
		copy(a.bytes[:], st[i*defaultBufSize:i*defaultBufSize+defaultBufSize])
		ww.Write(a.bytes[:])
	}
	if len(st)%defaultBufSize > 0 {
		copy(a.bytes[:], st[numI*defaultBufSize:])
		ww.Write(a.bytes[:len(st)%defaultBufSize])
	}
	printStringSafeBytesPool.Put(a)
	return
}

func PrintString(ww io.Writer, st string) (int, error) {
	if st == "" {
		return 0, nil
	}
	return ww.Write(*(*[]byte)(unsafe.Pointer(&st)))
}
