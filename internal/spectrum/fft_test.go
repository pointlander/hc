package spectrum

import (
	"math"
	"math/cmplx"
	"testing"
)

func TestFFTSinusoidPeak(t *testing.T) {
	const n, k0 = 256, 17
	x := make([]complex128, n)
	for i := range x {
		x[i] = cmplx.Exp(complex(0, 2*math.Pi*float64(k0)*float64(i)/float64(n)))
	}
	FFT(x)
	best, bestMag := 0, 0.0
	for k, v := range x {
		m := cmplx.Abs(v)
		if m > bestMag {
			best, bestMag = k, m
		}
	}
	if best != k0 {
		t.Fatalf("peak bin %d, want %d", best, k0)
	}
	if math.Abs(bestMag-float64(n)) > 1e-6 {
		t.Fatalf("peak magnitude %g, want %d", bestMag, n)
	}
}

func TestFFTRejectsNonPow2(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic")
		}
	}()
	FFT(make([]complex128, 3))
}
