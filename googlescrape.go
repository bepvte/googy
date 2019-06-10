// +build ignore
package main

// #cgo LDFLAGS: -L${SRCDIR}/libs -Wl,-rpath,$ORIGIN/libs -Wl,-rpath,$ORIGIN -lgooglescrape
// #include "libs/googlescrape.h"
import "C"
import (
    "unsafe"
	"errors"
	"bytes"
	"github.com/vmihailenco/msgpack"
)
//not used

type result struct {
	Link string `msgpack:"link"`
	Description string `msgpack:"description"`
}
type resultCollection []result

func google(s string) ([]result, error) {
	cs := C.CString(s)
	r := C.google(cs)
	C.free(unsafe.Pointer(cs))
	defer C.freeGResults(r)
	eboy := C.GoBytes(unsafe.Pointer(r.val), r.len)
	decoder := msgpack.NewDecoder(bytes.NewReader(eboy))
	if x, err := decoder.DecodeString(); err == nil {
		return nil, errors.New(x)
	}
	decoder.Reset(bytes.NewReader(eboy))
	var ret resultCollection
	if err := decoder.Decode(&ret); err != nil {
		return nil, err
	}
	return ret, nil
}