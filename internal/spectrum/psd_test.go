package spectrum

import (
	"math"
	"math/cmplx"
	"testing"
)

func TestWelchTonePeak(t *testing.T) {
	const n, fs, tone = 8192, 2e6, 50e3
	iq := make([]complex128, n)
	for i := range iq {
		iq[i] = cmplx.Exp(complex(0, 2*math.Pi*tone*float64(i)/fs))
	}
	p := Welch(iq, fs, 1024, 512)
	best, bestP := 0, 0.0
	for k, v := range p.Power {
		if v > bestP {
			best, bestP = k, v
		}
	}
	got := p.BinHz(best)
	if math.Abs(got-tone) > fs/1024*1.5 {
		t.Fatalf("peak %g Hz, want %g", got, tone)
	}
}

func TestWelchNegativeFreq(t *testing.T) {
	const n, fs, tone = 8192, 2e6, -250e3
	iq := make([]complex128, n)
	for i := range iq {
		iq[i] = cmplx.Exp(complex(0, 2*math.Pi*tone*float64(i)/fs))
	}
	p := Welch(iq, fs, 1024, 512)
	best, bestP := 0, 0.0
	for k, v := range p.Power {
		if v > bestP {
			best, bestP = k, v
		}
	}
	got := p.BinHz(best)
	if math.Abs(got-tone) > fs/1024*1.5 {
		t.Fatalf("peak %g Hz, want %g", got, tone)
	}
}

func TestNoiseFloorBelowTone(t *testing.T) {
	const n, fs, tone = 16384, 1e6, 100e3
	iq := make([]complex128, n)
	for i := range iq {
		iq[i] = 0.01*cmplx.Exp(complex(0, 2*math.Pi*tone*float64(i)/fs)) + 0.0001
	}
	p := Welch(iq, fs, 1024, 512)
	nf := p.NoiseFloorDB()
	best := -1e9
	for k := range p.Power {
		if db := p.DB(k); db > best {
			best = db
		}
	}
	if best-nf < 20 {
		t.Fatalf("tone %.1f dB, floor %.1f dB", best, nf)
	}
}
