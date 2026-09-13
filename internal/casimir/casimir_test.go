package casimir

import (
	"math"
	"math/cmplx"
	"math/rand"
	"strings"
	"testing"
	"time"
)

func TestDefaultOxideIs1um(t *testing.T) {
	d := DefaultDevice()
	if err := d.Validate(); err != nil {
		t.Fatal(err)
	}
	if math.Abs(d.Oxide-1e-6) > 0 {
		t.Fatalf("oxide=%g, want 1e-6", d.Oxide)
	}
}

func TestCapacitanceScales(t *testing.T) {
	d := DefaultDevice()
	c := d.Capacitance()
	if c <= 0 {
		t.Fatal("C")
	}
	d2 := d
	d2.Oxide *= 2
	if math.Abs(d2.Capacitance()-c/2) > 1e-18 {
		t.Fatalf("C vs d: %g vs %g", d2.Capacitance(), c/2)
	}
	d3 := d
	d3.EpsR *= 2
	if math.Abs(d3.Capacitance()-2*c) > 1e-18 {
		t.Fatalf("C vs εr: %g vs %g", d3.Capacitance(), 2*c)
	}
}

func TestLowFreqLooksCapacitive(t *testing.T) {
	d := DefaultDevice()
	f := 1e5
	z := d.Impedance(f)
	if imag(z) >= 0 {
		t.Fatalf("want capacitive Im(Z)<0, got %v", z)
	}
	cEst := -1 / (2 * math.Pi * f * imag(z))
	c := d.Capacitance()
	rel := math.Abs(cEst-c) / c
	if rel > 0.15 {
		t.Fatalf("C_est=%g C=%g rel=%g", cEst, c, rel)
	}
}

func TestCasimirPressureScalesD4(t *testing.T) {
	d := DefaultDevice()
	p1 := d.Casimir().Pressure
	d.Oxide *= 2
	p2 := d.Casimir().Pressure
	ratio := p1 / p2
	if math.Abs(ratio-16) > 2 {
		t.Fatalf("P(d)/P(2d)=%g, want ~16", ratio)
	}
	if p1 < 1e-5 || p1 > 0.01 {
		t.Fatalf("pressure %g Pa, expected ~0.3 mPa at 1 µm in Al2O3", p1)
	}
}

func TestNyquistMatches4kTReZ(t *testing.T) {
	d := DefaultDevice()
	s := d.Nyquist(915e6)
	want := 4 * kB * d.Temp * real(s.Z)
	if math.Abs(s.VocHz-want)/want > 1e-12 {
		t.Fatalf("Sv=%g want %g", s.VocHz, want)
	}
	if s.P50 <= 0 || s.P50dBmHz > -150 {
		t.Fatalf("available power too large: %g W/Hz (%g dBm/Hz)", s.P50, s.P50dBmHz)
	}
}

func TestIQLengthAndFinite(t *testing.T) {
	d := DefaultDevice()
	iq := d.IQ(915e6, 8e6, 1024, rand.New(rand.NewSource(1)))
	if len(iq) != 1024 {
		t.Fatalf("len=%d", len(iq))
	}
	var p float64
	for _, z := range iq {
		if cmplx.IsNaN(z) || cmplx.IsInf(z) {
			t.Fatalf("bad sample %v", z)
		}
		p += abs2(z)
	}
	if p == 0 {
		t.Fatal("zero power")
	}
}

func TestValidateRejectsStackedArms(t *testing.T) {
	d := DefaultDevice()
	d.Arm = d.Height
	if err := d.Validate(); err == nil {
		t.Fatal("expected error")
	}
}

func TestMarkdownHasDeviceAndForce(t *testing.T) {
	r := Report{Time: time.Date(2026, 9, 13, 0, 0, 0, 0, time.UTC), Device: DefaultDevice()}
	md := r.Markdown()
	for _, want := range []string{"Casimir", "1 µm", "E-shaped", "Nyquist", "Pa"} {
		if !strings.Contains(md, want) {
			t.Fatalf("missing %q", want)
		}
	}
}

func TestImpedanceAtBands(t *testing.T) {
	d := DefaultDevice()
	for _, f := range []float64{1e6, 10e6, 100e6, 915e6, 2.45e9} {
		z := d.Impedance(f)
		t.Logf("f=%.3g Hz Z=%+.3g %+.3gi |Z|=%.3g", f, real(z), imag(z), cmplx.Abs(z))
	}
}

func TestPhaseVelocity(t *testing.T) {
	d := DefaultDevice()
	v := d.phaseVelocity()
	want := c0 / math.Sqrt(d.EpsR)
	if math.Abs(v-want) > 1e-6 {
		t.Fatalf("%g vs %g", v, want)
	}
}
