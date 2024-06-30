package dsp

import (
	"math"
	"multimedia/fftw"
	"multimedia/lame"
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

func ZeroPadMatrix36(samples []float64) []float64 {
	zeros := make([]float64, 36)

	return append(samples, zeros...)
}

func ShiftMatrix72(samples []float64) []float64 {
	shiftedSamples := make([]float64, 72)
	copy(shiftedSamples[0:9], samples[63:72])
	copy(shiftedSamples[9:72], samples[0:63])

	return shiftedSamples
}

func RepeatMatrix72(samples []float64) []float64 {
	repeatedSamples := make([]float64, 18)
	for i := 0; i < 18; i++ {
		repeatedSamples[i] = samples[0] - samples[35-i] - samples[36+i] + samples[71-i]
	}

	return repeatedSamples
}

func MDCTLong(samples []float64) []float64 {
	mdctSamples := make([]float64, 36)
	copy(mdctSamples, samples)
	mdctSamples = ZeroPadMatrix36(mdctSamples)
	mdctSamples = ShiftMatrix72(mdctSamples)
	mdctSamples = RepeatMatrix72(mdctSamples)
	mdctSamples = DCTIV18(mdctSamples)

	return mdctSamples
}

func MDCTLongByLAME(samples []float64) []float64 {
	input32 := make([]float32, 36)
	for i := range samples {
		input32[i] = float32(samples[i])
	}
	output32 := make([]float32, 18)

	in := lame.NewArray(input32)
	out := lame.NewArray(output32)
	lame.MDCTLong(in, out)

	output64 := make([]float64, 18)
	for i := range output32 {
		output64[i] = float64(output32[i])
	}

	return output64
}
