package main

// #cgo LDFLAGS: -L${SRCDIR}/libs -lgooglescrape
// #include "googlescrape.h"
import "C"
import "unsafe"
import "errors"

type result struct {
	url, desc string
}

func google(s string) (ret []result, err error) {
	cs := C.CString(s)
	r := C.google(cs)
	C.free(unsafe.Pointer(cs))
	defer func() {
		if r.err == nil {
			C.freeGResults(r)
		}
	}()
	if r.err != nil { return []result{}, errors.New(C.GoString(r.err))}
	res := (*[3]C.GResult)(unsafe.Pointer(r.ret))[:3:3]
	ret = make([]result, 3)
	for x:=0; x<=2; x++ {
		ret[x] = result{
			url: C.GoString(res[x].link),
			desc: C.GoString(res[x].description),
		}
	}
	return
}