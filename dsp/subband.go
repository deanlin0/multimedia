package dsp

import (
	"math"
	"multimedia/fftw"
)

// 32-point DCT
func DCT32(samples []float64) []float64 {
	var bins []float64
	for k := 0; k < 32; k++ {
		var bin float64
		for n := 0; n < 32; n++ {
			sample := samples[n]
			phase := math.Pi * float64(k) * (float64(2*n) + 1) / 64
			filter := math.Cos(phase)
			bin += sample * filter
		}
		bin *= 2
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
