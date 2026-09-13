package casimir

import (
	"math"
	"math/cmplx"
)

// Impedance is the Thevenin impedance of the flat E as seen from a probe
// at the junction of the spine and the center tine.
func (d Device) Impedance(freq float64) complex128 {
	if freq < 1 {
		freq = 1
	}
	armLen := d.armLen()
	leg := d.Arm + d.slot() // center-line of center tine to center-line of an outer tine

	zcArm, gArm := d.line(freq, d.Arm)
	zcSp, gSp := d.line(freq, d.Spine)

	zTop := parallel(openStub(zcArm, gArm, armLen), openStub(zcSp, gSp, d.Arm/2))
	zBot := parallel(openStub(zcArm, gArm, armLen), openStub(zcSp, gSp, d.Arm/2))
	zCtr := openStub(zcArm, gArm, armLen)
	zUp := loadedLine(zcSp, gSp, leg, zTop)
	zDn := loadedLine(zcSp, gSp, leg, zBot)
	return parallel(zCtr, parallel(zUp, zDn))
}

// Capacitance is the low-frequency MIM capacitance of both gaps (F).
func (d Device) Capacitance() float64 {
	// Two plates, each C = ε A / d, in parallel from the E's point of view.
	return 2 * eps0 * d.EpsR * d.Area() / d.Oxide
}

func (d Device) line(freq, width float64) (zc, gamma complex128) {
	w := width
	h := d.Oxide
	omega := 2 * math.Pi * freq
	// Symmetric stripline: E sheet lying flat between two grounds, gap h on each side.
	cp := 2 * eps0 * d.EpsR * w / h
	lp := mu0 * h / (2 * w)
	delta := math.Sqrt(2 / (omega * mu0 * d.Sigma))
	rs := 1 / (d.Sigma * delta)
	rp := 4 * rs / w
	gp := omega * cp * d.TanD
	zser := complex(rp, omega*lp)
	ysh := complex(gp, omega*cp)
	gamma = cmplx.Sqrt(zser * ysh)
	if real(gamma) < 0 {
		gamma = -gamma
	}
	zc = cmplx.Sqrt(zser / ysh)
	if real(zc) < 0 {
		zc = -zc
	}
	return zc, gamma
}

func openStub(zc, gamma complex128, length float64) complex128 {
	th := gamma * complex(length, 0)
	t := cmplx.Tanh(th)
	if cmplx.Abs(t) < 1e-18 {
		return complex(1e18, 0)
	}
	return zc / t
}

func loadedLine(zc, gamma complex128, length float64, zL complex128) complex128 {
	th := gamma * complex(length, 0)
	t := cmplx.Tanh(th)
	return zc * (zL + zc*t) / (zc + zL*t)
}

func parallel(a, b complex128) complex128 {
	s := a + b
	if s == 0 {
		return 0
	}
	return a * b / s
}

// SeriesResonance is the interior frequency of minimum |Z| in (lo, hi).
// ok is false when |Z| is monotonic on the interval (RC-like).
func (d Device) SeriesResonance(lo, hi float64) (freq, absZ float64, ok bool) {
	const n = 512
	bestF, bestZ := lo, math.Inf(1)
	logLo, logHi := math.Log(lo), math.Log(hi)
	for i := 0; i <= n; i++ {
		f := math.Exp(logLo + (logHi-logLo)*float64(i)/float64(n))
		az := cmplx.Abs(d.Impedance(f))
		if az < bestZ {
			bestZ = az
			bestF = f
		}
	}
	ok = bestF > lo*1.05 && bestF < hi*0.95
	return bestF, bestZ, ok
}

// Quality is a crude Q ≈ f / Δf_3dB around the series-resonance dip of |Z|.
func (d Device) Quality(f0 float64) float64 {
	if f0 <= 0 {
		return 0
	}
	z0 := cmplx.Abs(d.Impedance(f0))
	target := z0 * math.Sqrt2
	lo, hi := f0, f0
	for f := f0; f > f0/20 && f > 1e5; f *= 0.98 {
		if cmplx.Abs(d.Impedance(f)) >= target {
			lo = f
			break
		}
		lo = f
	}
	for f := f0; f < f0*20 && f < 6e9; f *= 1.02 {
		if cmplx.Abs(d.Impedance(f)) >= target {
			hi = f
			break
		}
		hi = f
	}
	bw := hi - lo
	if bw <= 0 {
		return 0
	}
	return f0 / bw
}
