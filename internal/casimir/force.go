package casimir

import "math"

// Force is the Casimir attraction on the E from both anodized plates.
type Force struct {
	Pressure float64 // Pa, per gap
	Force    float64 // N, both faces
	Energy   float64 // J, both gaps
	Area     float64 // m²
	MechHz   float64 // plate bending fundamental, Hz
	Xthermal float64 // rms thermal displacement, m
}

// Casimir evaluates the Lifshitz leading term for two Al plates filled
// with Al2O3 of thickness Oxide. Perfect-conductor vacuum pressure is
// reduced by n=√εr (mode speed in the dielectric) and a small plasma
// correction for aluminum at d ~ 1 µm (λp ≈ 107 nm ≪ d).
func (d Device) Casimir() Force {
	area := d.Area()
	n := math.Sqrt(d.EpsR)
	p0 := math.Pi * math.Pi * hbar * c0 / (240 * math.Pow(d.Oxide, 4))
	p := p0 / n * plasmaReduction(d.Oxide)
	// Both faces of the E see a gap.
	f := 2 * p * area
	// U = −P A d / 3 per gap (since P ∝ 1/d^4 ⇒ U ∝ 1/d^3).
	u := -2 * p * area * d.Oxide / 3

	mech, xth := d.mechanics()
	return Force{
		Pressure: p,
		Force:    f,
		Energy:   u,
		Area:     area,
		MechHz:   mech,
		Xthermal: xth,
	}
}

func plasmaReduction(d float64) float64 {
	// Aluminum plasma wavelength ~ 107 nm. At 1 µm the finite-conductivity
	// reduction of the T=0 Casimir pressure is only a few percent.
	const lambdaP = 107e-9
	x := lambdaP / d
	if x <= 0 {
		return 1
	}
	eta := 1 - 4*x/3
	if eta < 0.2 {
		return 0.2
	}
	return eta
}

func (d Device) mechanics() (freq, xrms float64) {
	// Kirchhoff plate, simply supported, mass density μ = ρ t.
	const (
		young = 70e9
		nu    = 0.33
		rho   = 2700
	)
	t := d.Thick
	flex := young * t * t * t / (12 * (1 - nu*nu))
	mu := rho * t
	if mu <= 0 || flex <= 0 {
		return 0, 0
	}
	lx, ly := d.Width, d.Height
	omega := math.Pi * math.Pi * (1/(lx*lx) + 1/(ly*ly)) * math.Sqrt(flex/mu)
	freq = omega / (2 * math.Pi)
	mass := rho * d.Area() * t
	if omega <= 0 || mass <= 0 {
		return freq, 0
	}
	k := mass * omega * omega
	xrms = math.Sqrt(kB * d.Temp / k)
	return freq, xrms
}
