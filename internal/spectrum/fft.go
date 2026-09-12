package spectrum

import "math"

// FFT computes an in-place radix-2 Cooley–Tukey FFT on x.
// len(x) must be a power of two.
func FFT(x []complex128) {
	n := len(x)
	if n == 0 || n&(n-1) != 0 {
		panic("spectrum: FFT length must be a power of two")
	}
	bitReverse(x)
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

func bitReverse(x []complex128) {
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
}

func hann(n int) []float64 {
	w := make([]float64, n)
	if n <= 1 {
		if n == 1 {
			w[0] = 1
		}
		return w
	}
	for i := 0; i < n; i++ {
		w[i] = 0.5 * (1 - math.Cos(2*math.Pi*float64(i)/float64(n-1)))
	}
	return w
}
