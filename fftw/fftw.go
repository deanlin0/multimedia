package fftw

// #cgo LDFLAGS: -lfftw3
// #include <fftw3.h>
import "C"
import "unsafe"

var (
	FFTWREDFT10 = FFTType{
		cKind: C.FFTW_REDFT10,
	}
)

var (
	FFTWESTIMATE = Flag{
		cFlag: C.FFTW_ESTIMATE,
	}
)

type FFTType struct {
	cKind C.fftw_r2r_kind
}

type Flag struct {
	cFlag C.uint
}

type Plan struct {
	cPlan C.fftw_plan
}

type Array struct {
	cArray *C.double
}

func NewArray(data []float64) *Array {
	return &Array{
		cArray: (*C.double)(unsafe.Pointer(&data[0])),
	}
}

func NewPlan(n int, in, out *Array, fftType FFTType, flag Flag) *Plan {
	cPlan := C.fftw_plan_r2r_1d(C.int(n), in.cArray, out.cArray, fftType.cKind, flag.cFlag)
	return &Plan{cPlan}
}

func ExecutePlan(plan *Plan) {
	C.fftw_execute(plan.cPlan)
}
