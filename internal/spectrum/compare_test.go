package spectrum

import (
	"math"
	"math/cmplx"
	"strings"
	"testing"
	"time"
)

func toneIQ(n int, fs, hz, amp float64) []complex128 {
	iq := make([]complex128, n)
	for i := range iq {
		iq[i] = complex(amp, 0) * cmplx.Exp(complex(0, 2*math.Pi*hz*float64(i)/fs))
	}
	return iq
}

func add(a, b []complex128) []complex128 {
	out := make([]complex128, len(a))
	for i := range a {
		out[i] = a[i] + b[i]
	}
	return out
}

func TestCompareIdentical(t *testing.T) {
	const n, fs = 8192, 1e6
	iq := toneIQ(n, fs, 50e3, 0.1)
	for i := range iq {
		iq[i] += complex(0.002, -0.001)
	}
	a := Welch(iq, fs, 1024, 512)
	d := Compare(a, a, CompareOptions{NSigma: 6, MinDB: 3})
	if math.Abs(d.OffsetDB) > 1e-9 {
		t.Fatalf("offset %g", d.OffsetDB)
	}
	if len(d.Peaks) != 0 {
		t.Fatalf("peaks=%d, want 0 on identical spectra", len(d.Peaks))
	}
}

func TestCompareConstantOffset(t *testing.T) {
	const n, fs = 8192, 1e6
	base := toneIQ(n, fs, 80e3, 0.05)
	hot := make([]complex128, n)
	copy(hot, base)
	// +6 dB is amplitude * 2
	for i := range hot {
		hot[i] *= 2
	}
	a := Welch(hot, fs, 1024, 512)
	b := Welch(base, fs, 1024, 512)
	d := Compare(a, b, CompareOptions{NSigma: 6, MinDB: 3})
	if math.Abs(d.OffsetDB-6) > 0.5 {
		t.Fatalf("offset %g, want ~6 dB", d.OffsetDB)
	}
	for _, p := range d.Peaks {
		if math.Abs(p.FreqHz-80e3) < fs/1024*2 {
			t.Fatalf("common tone treated as significant: %+v", p)
		}
	}
}

func TestCompareFindsUniqueSpur(t *testing.T) {
	const n, fs = 16384, 1e6
	common := add(toneIQ(n, fs, 50e3, 0.1), toneIQ(n, fs, 0, 0.02)) // DC + tone
	aIQ := add(common, toneIQ(n, fs, 0, 0.05))                      // extra DC on A
	bIQ := add(common, toneIQ(n, fs, -200e3, 0.08))                 // unique spur on B
	a := Welch(aIQ, fs, 1024, 512)
	b := Welch(bIQ, fs, 1024, 512)
	d := Compare(a, b, CompareOptions{NSigma: 6, MinDB: 3})
	if len(d.Peaks) == 0 {
		t.Fatal("expected significant peaks")
	}
	var sawDC, sawSpur bool
	for _, p := range d.Peaks {
		if math.Abs(p.FreqHz) < fs/1024 {
			sawDC = true
			if p.ResidualDB <= 0 {
				t.Fatalf("DC residual should be A hotter: %+v", p)
			}
		}
		if math.Abs(p.FreqHz+200e3) < fs/1024*2 {
			sawSpur = true
			if p.ResidualDB >= 0 {
				t.Fatalf("spur residual should be B hotter: %+v", p)
			}
		}
	}
	if !sawDC {
		t.Fatal("did not flag DC difference")
	}
	if !sawSpur {
		t.Fatalf("did not flag -200 kHz spur; peaks=%v", d.Peaks)
	}
}

func TestMarkdownContainsSignificantBins(t *testing.T) {
	const n, fs = 16384, 1e6
	common := add(toneIQ(n, fs, 50e3, 0.1), toneIQ(n, fs, 0, 0.02))
	a := Welch(add(common, toneIQ(n, fs, 0, 0.05)), fs, 1024, 512)
	b := Welch(add(common, toneIQ(n, fs, -200e3, 0.08)), fs, 1024, 512)
	d := Compare(a, b, CompareOptions{NSigma: 6, MinDB: 3})
	r := Report{
		Time:     time.Date(2026, 9, 12, 0, 0, 0, 0, time.UTC),
		RateHz:   fs,
		Duration: 250 * time.Millisecond,
		FFTSize:  1024,
		Overlap:  512,
		LNAGain:  16,
		VGAGain:  16,
		NSigma:   6,
		MinDB:    3,
		A:        Radio{Name: "288e2dc3", Serial: "aaa", Samples: n},
		B:        Radio{Name: "2c8f32c3", Serial: "bbb", Samples: n},
		Bands:    []Band{{CenterHz: 915e6, RateHz: fs, A: a, B: b, Diff: d}},
	}
	md := r.Markdown()
	for _, want := range []string{"HackRF spectral difference", "Significant bins", "915", "288e2dc3"} {
		if !strings.Contains(md, want) {
			t.Fatalf("markdown missing %q", want)
		}
	}
	if len(d.Peaks) == 0 {
		t.Fatal("fixture should have peaks")
	}
}

func TestASCIIPlotZeroNoStem(t *testing.T) {
	n := 64
	p := PSD{FFTSize: n, SampleRate: 1e6, Power: make([]float64, n)}
	resid := make([]float64, n)
	resid[8] = 6
	s := ASCIIPlot(p, resid, 32, 9, 3)
	if strings.Count(s, "|") > 40 {
		t.Fatalf("stems filled the plot:\n%s", s)
	}
	if !strings.Contains(s, "#") {
		t.Fatalf("missing significant mark:\n%s", s)
	}
}

func TestFmtHz(t *testing.T) {
	if g := fmtHz(915e6); g != "+915 MHz" {
		t.Fatalf("%q", g)
	}
	if g := fmtHz(-250e3); g != "-250 kHz" {
		t.Fatalf("%q", g)
	}
	if g := fmtHz(0); g != "0 Hz" {
		t.Fatalf("%q", g)
	}
}
