package casimir

import (
	"math"
	"math/cmplx"
	"math/rand"
)

// IQ synthesizes complex baseband samples of the voltage across Zload
// at centerHz. Values are in volts. The spectrum is colored by Nyquist(f).
func (d Device) IQ(centerHz, rate float64, n int, rng *rand.Rand) []complex128 {
	if n < 16 || n&(n-1) != 0 {
		panic("casimir: IQ length must be a power of two ≥ 16")
	}
	if rng == nil {
		rng = rand.New(rand.NewSource(1))
	}
	df := rate / float64(n)
	spec := make([]complex128, n)
	for k := 0; k < n; k++ {
		fbb := binHz(k, n, rate)
		s := d.Nyquist(centerHz + fbb)
		h := complex(d.Zload, 0) / (s.Z + complex(d.Zload, 0))
		// RMS volts in this bin, circularly symmetric complex Gaussian.
		sigma := math.Sqrt(s.VocHz * abs2(h) * df)
		g := complex(rng.NormFloat64(), rng.NormFloat64()) * complex(sigma/math.Sqrt2, 0)
		if k == 0 || k == n/2 {
			g = complex(rng.NormFloat64()*sigma, 0)
		}
		spec[k] = g * complex(float64(n), 0)
	}
	return ifft(spec)
}

func binHz(k, n int, rate float64) float64 {
	if k <= n/2 {
		return float64(k) * rate / float64(n)
	}
	return float64(k-n) * rate / float64(n)
}

func ifft(x []complex128) []complex128 {
	n := len(x)
	y := append([]complex128(nil), x...)
	for i, v := range y {
		y[i] = cmplx.Conj(v)
	}
	fft(y)
	inv := 1 / float64(n)
	for i, v := range y {
		y[i] = cmplx.Conj(v) * complex(inv, 0)
	}
	return y
}

func fft(x []complex128) {
	n := len(x)
	for i, j := 0, 0; i < n; i++ {
		if i < j {
			x[i], x[j] = x[j], x[i]
		}
		m := n >> 1
		for m >= 1 && j >= m {
			j -= m
			m >>= 1
		}
		j += m
	}
	for size := 2; size <= n; size <<= 1 {
		half := size / 2
		ang := -2 * math.Pi / float64(size)
		wStep := complex(math.Cos(ang), math.Sin(ang))
		for start := 0; start < n; start += size {
			w := complex(1, 0)
			for k := 0; k < half; k++ {
				even := x[start+k]
				odd := x[start+k+half] * w
				x[start+k] = even + odd
				x[start+k+half] = even - odd
				w *= wStep
			}
		}
	}
}
