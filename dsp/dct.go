package dsp

import (
	"math"
	"multimedia/fftw"
)

// 32-point DCT
func DCT32(samples []float64) []float64 {
	bins := make([]float64, 32)
	for k := 0; k < 32; k++ {
		var bin float64
		for n := 0; n < 32; n++ {
			sample := samples[n]
			phase := math.Pi * float64(k) * (float64(2*n) + 1) / 64
			filter := math.Cos(phase)
			bin += sample * filter
		}
		bin *= 2
		bins[k] = bin
	}

	return bins
}

func DFT64(samples []float64) []complex128 {
	bins := make([]complex128, 64)
	for k := 0; k < 64; k++ {
		var bin complex128
		for n := 0; n < 64; n++ {
			sample := samples[n]
			phase := 2 * math.Pi * float64(k) * float64(n) / 64
			bin += complex(sample*math.Cos(phase), sample*math.Sin(phase))
		}
		bins[k] = bin
	}

	return bins
}

func DCT32ByDFT(samples []float64) []float64 {
	logicalSamples := make([]float64, 64)
	for i, sample := range samples {
		logicalSamples[i] = sample
		logicalSamples[63-i] = sample
	}
	dftBins := DFT64(logicalSamples)

	bins := make([]float64, 32)
	for k := 0; k < 32; k++ {
		phase := math.Pi * float64(k) / 64
		shift := complex(math.Cos(phase), math.Sin(phase))
		bin := real(dftBins[k] * shift)
		bins[k] = bin
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

func DCTIV18(samples []float64) []float64 {
	bins := make([]float64, 18)
	for k := 0; k < 18; k++ {
		var bin float64
		for n := 0; n < 18; n++ {
			sample := samples[n]
			phase := math.Pi * (float64(2*k) + 1) * (float64(2*n) + 1) / 72
			filter := math.Cos(phase)
			bin += sample * filter
		}
		bin *= 2
		bins[k] = bin
	}

	return bins
}

func ZeroPad36(samples []float64) []float64 {
	zeros := make([]float64, 36)

	return append(samples, zeros...)
}

func HalfShift72(samples []float64) []float64 {
	shiftedSamples := make([]float64, 72)
	copy(shiftedSamples[0:18], samples[54:72])
	copy(shiftedSamples[54:72], samples[0:54])

	return shiftedSamples
}

func Overlap72(samples []float64) []float64 {
	overlappedSamples := make([]float64, 18)
	for i := 0; i < 18; i++ {
		overlappedSamples[i] = samples[0] - samples[35-i] - samples[36+i] + samples[71-i]
	}

	return overlappedSamples
}
