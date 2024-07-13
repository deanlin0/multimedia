package lame

// #cgo LDFLAGS: -lnewmdct -lmp3lame
// #include <strings.h>
// #include <stdint.h>
// #include "lame/config.h"
// #include "lame/lame.h"
// #include "lame/lame_global_flags.h"
// #include "lame/machine.h"
// #include "lame/encoder.h"
// #include "lame/l3side.h"
// #include "lame/id3tag.h"
// #include "lame/util.h"
// #include "lame/newmdct.h"
import "C"
import "unsafe"

type Array struct {
	cArray *C.float
}

func NewArray(data []float32) *Array {
	return &Array{
		cArray: (*C.float)(unsafe.Pointer(&data[0])),
	}
}

func MDCTShort(inOut *Array) {
	C.mdct_short(inOut.cArray)
}

func MDCTLong(in, out *Array) {
	C.mdct_long(out.cArray, in.cArray)
}
