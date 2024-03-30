package dsp

import (
	"math"
	"math/big"
	"streaming/fftw"
)

// 32-point DCT
func DCT32(samples []float64) []float64 {
	var bins []float64
	for j := 0; j < 32; j++ {
		binWithPrec := big.NewFloat(0).SetPrec(100)
		for k := 0; k < 32; k++ {
			sample := samples[k]
			phase := math.Pi * float64(j) * (float64(2*k) + 1) / 64
			filter := math.Cos(phase)
			binWithPrec.Add(
				binWithPrec,
				big.NewFloat(0).SetPrec(100).
					Mul(
						big.NewFloat(sample).SetPrec(100),
						big.NewFloat(filter).SetPrec(100),
					),
			)
		}

		bin, _ := binWithPrec.Mul(
			binWithPrec,
			big.NewFloat(2).SetPrec(100),
		).Float64()
		bins = append(bins, bin)
	}

	return bins
}

func DCT32ByFFTW(samples []float64) []float64 {
	output := make([]float64, 32)
	_ = copy(output, samples)

	in := fftw.NewArray(output)
	plan := fftw.NewPlan(32, in, in, fftw.FFTWREDFT10, fftw.FFTWESTIMATE)
	fftw.ExecutePlan(plan)

	return output
}
