package casimir

import (
	"fmt"
	"math"
)

const (
	c0   = 299792458.0
	eps0 = 8.8541878128e-12
	mu0  = 1.25663706212e-6
	hbar = 1.054571817e-34
	kB   = 1.380649e-23
	eta0 = 376.730313461
)

// Device is an E-shaped aluminum conductor sandwiched between two
// anodized aluminum plates. The anodic Al2O3 on each inner face is the
// Casimir / MIM gap.
type Device struct {
	Width  float64 // arm-length direction, m
	Height float64 // across the three arms, m
	Spine  float64 // vertical bar width, m
	Arm    float64 // each of the three arm widths, m
	Thick  float64 // E-sheet thickness, m
	Oxide  float64 // anodization thickness per plate, m
	EpsR   float64 // Al2O3 relative permittivity
	TanD   float64 // dielectric loss tangent
	Sigma  float64 // aluminum conductivity, S/m
	Temp   float64 // kelvin
	PatchV float64 // typical patch potential, V
	Zload  float64 // probe impedance, Ω
}

// DefaultDevice is a 50×40 mm E in 0.4 mm Al sheet between plates
// anodized to 1 µm. Lateral size is chosen so the lowest stripline
// self-resonance sits in the VHF/UHF range the HackRFs already scan.
func DefaultDevice() Device {
	return Device{
		Width:  50e-3,
		Height: 40e-3,
		Spine:  8e-3,
		Arm:    8e-3,
		Thick:  0.4e-3,
		Oxide:  1e-6,
		EpsR:   9.8,
		TanD:   0.015,
		Sigma:  3.56e7,
		Temp:   293.15,
		PatchV: 0.05,
		Zload:  50,
	}
}

func (d Device) slotH() float64 {
	return (d.Height - 3*d.Arm) / 2
}

func (d Device) armLen() float64 {
	return d.Width - d.Spine
}

// Area is the in-plane metal area of the E (bounding rectangle minus the two slots).
func (d Device) Area() float64 {
	s := d.slotH()
	if s < 0 {
		s = 0
	}
	return d.Width*d.Height - 2*d.armLen()*s
}

func (d Device) Validate() error {
	if d.Oxide <= 0 || d.Width <= 0 || d.Height <= 0 || d.Spine <= 0 || d.Arm <= 0 || d.Thick <= 0 {
		return fmt.Errorf("casimir: dimensions must be positive")
	}
	if d.Spine >= d.Width {
		return fmt.Errorf("casimir: spine wider than the E")
	}
	if 3*d.Arm >= d.Height {
		return fmt.Errorf("casimir: three arms do not fit in height")
	}
	if d.EpsR < 1 {
		return fmt.Errorf("casimir: εr must be ≥ 1")
	}
	if d.Sigma <= 0 || d.Temp <= 0 || d.Zload <= 0 {
		return fmt.Errorf("casimir: sigma, temp, and Zload must be positive")
	}
	return nil
}

func (d Device) phaseVelocity() float64 {
	return c0 / math.Sqrt(d.EpsR)
}
