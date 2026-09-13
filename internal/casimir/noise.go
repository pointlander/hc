package casimir

import "math"

// Spectrum is the Nyquist radio output of the device at one frequency.
type Spectrum struct {
	Freq     float64
	Z        complex128
	VocHz    float64 // open-circuit voltage PSD, V²/Hz
	P50      float64 // available power into Zload, W/Hz
	P50dBmHz float64
}

// Nyquist is the Johnson–Nyquist spectrum of the E at freq, including
// mismatch into Zload. At radio frequencies ħω ≪ kT, so the Casimir
// zero-point term does not radiate; the distinctive shape is Re(Z(f)).
func (d Device) Nyquist(freq float64) Spectrum {
	z := d.Impedance(freq)
	re := real(z)
	if re < 0 {
		re = 0
	}
	voc := 4 * kB * d.Temp * re
	zl := d.Zload
	h := complex(zl, 0) / (z + complex(zl, 0))
	vload2 := voc * abs2(h)
	p := 0.0
	if zl > 0 {
		p = vload2 / zl
	}
	return Spectrum{
		Freq:     freq,
		Z:        z,
		VocHz:    voc,
		P50:      p,
		P50dBmHz: dbmHz(p),
	}
}

func abs2(z complex128) float64 {
	return real(z)*real(z) + imag(z)*imag(z)
}

func dbmHz(w float64) float64 {
	if w <= 1e-40 {
		return -400
	}
	return 10*math.Log10(w) + 30
}

// Sweep evaluates Nyquist on log-spaced frequencies.
func (d Device) Sweep(lo, hi float64, n int) []Spectrum {
	if n < 2 {
		n = 2
	}
	out := make([]Spectrum, n)
	logLo, logHi := math.Log(lo), math.Log(hi)
	for i := 0; i < n; i++ {
		f := math.Exp(logLo + (logHi-logLo)*float64(i)/float64(n-1))
		out[i] = d.Nyquist(f)
	}
	return out
}

// PatchTransduced is the rms voltage from thermal plate motion acting on
// the patch-potential biased MIM capacitor: δV = V_patch · x / d.
// This channel lives at the mechanical frequency, far below HackRF.
func (d Device) PatchTransduced(xrms float64) float64 {
	if d.Oxide <= 0 {
		return 0
	}
	return d.PatchV * xrms / d.Oxide
}
