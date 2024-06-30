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

// type LameContext struct {
// 	cContext *C.lame_global_flags
// }

// func (ctx *LameContext) Output() []float64 {
// 	xr := ctx.cContext.internal_flags.l3_side.tt[0][0].xr
// 	output := make([]float64, 576)
// 	for i := range xr {
// 		output[i] = float64(xr[i])
// 	}
// 	return output
// }

// func AddEncodeDelay(data []float32, delay int) []float32 {
// 	data = append(make([]float32, delay), data...)
// 	return data
// }

// func NewLameContext() *LameContext {
// 	cContext := C.lame_init()
// 	_ = C.lame_init_params(cContext)
// 	cContext.internal_flags.cfg.channels_in = 1
// 	cContext.internal_flags.cfg.channels_out = 1
// 	return &LameContext{cContext: cContext}
// }

// func ExecuteMDCTWithSubband(ctx *LameContext, inLeft, inRight *Array) {
// 	C.mdct_sub48(ctx.cContext.internal_flags, inLeft.cArray, inRight.cArray)
// }

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
