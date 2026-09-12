package spectrum

import "sort"

// PSD is a two-sided Welch power spectral density of complex I/Q.
// Bin 0 is DC; bins 1..n/2-1 are positive frequencies; n/2 is Nyquist;
// bins n/2+1..n-1 are negative frequencies.
type PSD struct {
	Power      []float64 // linear mean |X[k]|^2 / sum(w^2), full-scale = 1
	Windows    int
	FFTSize    int
	SampleRate float64
}

// BinHz is the baseband frequency of bin k (negative for the upper half).
func (p PSD) BinHz(k int) float64 {
	n := p.FFTSize
	if n == 0 {
		return 0
	}
	if k <= n/2 {
		return float64(k) * p.SampleRate / float64(n)
	}
	return float64(k-n) * p.SampleRate / float64(n)
}

// DB is the bin power in dBFS.
func (p PSD) DB(k int) float64 {
	if k < 0 || k >= len(p.Power) {
		return -180
	}
	return db(p.Power[k])
}

// Welch estimates the two-sided PSD of complex I/Q with a Hann window.
func Welch(iq []complex128, sampleRate float64, nfft, noverlap int) PSD {
	if nfft < 16 || nfft&(nfft-1) != 0 {
		panic("spectrum: Welch nfft must be a power of two >= 16")
	}
	if noverlap < 0 || noverlap >= nfft {
		noverlap = nfft / 2
	}
	hop := nfft - noverlap
	w := hann(nfft)
	var wss float64
	for _, v := range w {
		wss += v * v
	}
	if wss == 0 {
		wss = 1
	}

	acc := make([]float64, nfft)
	x := make([]complex128, nfft)
	windows := 0
	for start := 0; start+nfft <= len(iq); start += hop {
		for i := 0; i < nfft; i++ {
			x[i] = iq[start+i] * complex(w[i], 0)
		}
		FFT(x)
		for k := 0; k < nfft; k++ {
			re := real(x[k])
			im := imag(x[k])
			acc[k] += (re*re + im*im) / wss
		}
		windows++
	}
	if windows == 0 {
		return PSD{Power: acc, FFTSize: nfft, SampleRate: sampleRate}
	}
	inv := 1 / float64(windows)
	for k := range acc {
		acc[k] *= inv
	}
	return PSD{
		Power:      acc,
		Windows:    windows,
		FFTSize:    nfft,
		SampleRate: sampleRate,
	}
}

// NoiseFloorDB is the median bin power in dBFS, excluding DC and filter edges.
func (p PSD) NoiseFloorDB() float64 {
	n := len(p.Power)
	if n < 16 {
		return -180
	}
	lo := 4
	hi := n - 4
	edge := n / 10
	if edge < 4 {
		edge = 4
	}
	vals := make([]float64, 0, n)
	for k := lo; k < hi; k++ {
		if k > n/2-edge && k < n/2+edge {
			continue
		}
		if k == 0 {
			continue
		}
		vals = append(vals, p.Power[k])
	}
	if len(vals) == 0 {
		return -180
	}
	return db(median(vals))
}

func median(v []float64) float64 {
	if len(v) == 0 {
		return 0
	}
	s := append([]float64(nil), v...)
	sort.Float64s(s)
	if len(s)%2 == 1 {
		return s[len(s)/2]
	}
	return 0.5 * (s[len(s)/2-1] + s[len(s)/2])
}

func mad(v []float64, med float64) float64 {
	dev := make([]float64, len(v))
	for i, x := range v {
		d := x - med
		if d < 0 {
			d = -d
		}
		dev[i] = d
	}
	return median(dev)
}
