package casimir

import (
	"fmt"
	"math"
	"strings"
	"time"
)

// Report is the markdown write-up of the simulated radio output.
type Report struct {
	Time   time.Time
	Device Device
	Bands  []float64
}

// Markdown renders the Casimir-device radio simulation.
func (r Report) Markdown() string {
	d := r.Device
	var b strings.Builder
	fmt.Fprintf(&b, "# Casimir E-sandwich radio simulation\n\n")
	fmt.Fprintf(&b, "Simulated %s.\n\n", r.Time.UTC().Format("2006-01-02 15:04:05 UTC"))
	b.WriteString("An **E-shaped aluminum** sheet is clamped between two **anodized aluminum plates**. ")
	b.WriteString("Each inner face carries **1 µm** of anodic Al₂O₃, so the E sees two metal–insulator–metal gaps. ")
	b.WriteString("The Casimir pressure lives in those gaps. The radio output is the Johnson–Nyquist field of the same structure, ")
	b.WriteString("shaped by lossy stripline modes of the E. At RF, ħω ≪ kT, so zero-point energy does not radiate; ")
	b.WriteString("what a HackRF can in principle couple to is thermal, with a spectral shape set by Re(Z(f)).\n\n")

	b.WriteString("## Geometry\n\n")
	b.WriteString("```\n")
	b.WriteString(schematic())
	b.WriteString("```\n\n")
	fmt.Fprintf(&b, "| | |\n| --- | ---: |\n")
	fmt.Fprintf(&b, "| E outline | %.1f mm × %.1f mm |\n", d.Width*1e3, d.Height*1e3)
	fmt.Fprintf(&b, "| Spine / arm | %.1f mm / %.1f mm |\n", d.Spine*1e3, d.Arm*1e3)
	fmt.Fprintf(&b, "| Slot height | %.1f mm |\n", d.slotH()*1e3)
	fmt.Fprintf(&b, "| E thickness | %.2f mm |\n", d.Thick*1e3)
	fmt.Fprintf(&b, "| Anodization (each plate) | **%.3g µm** |\n", d.Oxide*1e6)
	fmt.Fprintf(&b, "| Al₂O₃ ε<sub>r</sub> / tanδ | %.1f / %.3f |\n", d.EpsR, d.TanD)
	fmt.Fprintf(&b, "| Al conductivity | %.2e S/m |\n", d.Sigma)
	fmt.Fprintf(&b, "| Temperature | %.1f K |\n", d.Temp)
	fmt.Fprintf(&b, "| Probe Z | %.0f Ω |\n", d.Zload)
	fmt.Fprintf(&b, "| E metal area | %.2f cm² |\n\n", d.Area()*1e4)

	c := d.Capacitance()
	v := d.phaseVelocity()
	f10 := v / (2 * d.Width)
	f01 := v / (2 * d.Height)
	fslot := v / (2 * (d.Width + 2*d.slotH()))
	fsr, zsr, haveSR := d.SeriesResonance(1e6, 6e9)
	q := 0.0
	if haveSR {
		q = d.Quality(fsr)
	}
	fmt.Fprintf(&b, "## Electromagnetics\n\n")
	fmt.Fprintf(&b, "The 1 µm gap makes this a **very low-impedance stripline** (Z<sub>0</sub> milliohms). ")
	fmt.Fprintf(&b, "Skin-effect loss in the aluminum dominates, so the geometric half-wave modes (TM<sub>10</sub> ~ %.0f MHz, slot path ~ %.0f MHz) are **overdamped**. ", f10/1e6, fslot/1e6)
	fmt.Fprintf(&b, "The structure behaves as a ~%.0f nF MIM capacitor: |Z| falls with frequency and the radio output is a smooth thermal continuum, not a comb of spurs.\n\n", c*1e9)
	fmt.Fprintf(&b, "| | |\n| --- | ---: |\n")
	fmt.Fprintf(&b, "| MIM capacitance (both gaps) | **%.2f nF** |\n", c*1e9)
	fmt.Fprintf(&b, "| Phase velocity c/√ε<sub>r</sub> | %.3f c |\n", v/c0)
	fmt.Fprintf(&b, "| Ideal TM<sub>10</sub> (along arms) | %.1f MHz |\n", f10/1e6)
	fmt.Fprintf(&b, "| Ideal TM<sub>01</sub> (across arms) | %.1f MHz |\n", f01/1e6)
	fmt.Fprintf(&b, "| Ideal slot-lengthened path | %.1f MHz |\n", fslot/1e6)
	if haveSR {
		fmt.Fprintf(&b, "| Series |Z| dip | **%s**, \\|Z\\|=%.2f mΩ, Q≈%.2f |\n\n", fmtHz(fsr), zsr*1e3, q)
	} else {
		fmt.Fprintf(&b, "| Series |Z| dip | none below 6 GHz (RC-like) |\n\n")
	}

	cas := d.Casimir()
	vpatch := d.PatchTransduced(cas.Xthermal)
	fmt.Fprintf(&b, "## Casimir force\n\n")
	fmt.Fprintf(&b, "Leading Lifshitz term for perfect reflectors filled with Al₂O₃, plus a plasma-wavelength correction for Al (λ<sub>p</sub>≈107 nm ≪ 1 µm):\n\n")
	fmt.Fprintf(&b, "$$P = \\frac{\\pi^2 \\hbar c}{240\\, n\\, d^4}\\,\\eta_{\\mathrm{Al}},\\quad n=\\sqrt{\\varepsilon_r}$$\n\n")
	fmt.Fprintf(&b, "| | |\n| --- | ---: |\n")
	fmt.Fprintf(&b, "| Pressure per gap | **%.2f mPa** |\n", cas.Pressure*1e3)
	fmt.Fprintf(&b, "| Force on E (both faces) | **%.2e N** (%.2f nN) |\n", cas.Force, cas.Force*1e9)
	fmt.Fprintf(&b, "| Energy both gaps | %.2e J |\n", cas.Energy)
	fmt.Fprintf(&b, "| Plate bending fundamental | **%.1f kHz** |\n", cas.MechHz/1e3)
	fmt.Fprintf(&b, "| Thermal x<sub>rms</sub> | %.2e m (%.2f pm) |\n", cas.Xthermal, cas.Xthermal*1e12)
	fmt.Fprintf(&b, "| Patch-potential δV<sub>rms</sub> (%.0f mV) | **%.2e V** |\n\n", d.PatchV*1e3, vpatch)
	b.WriteString("The mechanical channel is audio/ultrasonic, not a HackRF band. ")
	b.WriteString("It would only appear at UHF/microwave if an RF pump mixed with the motion (not assumed here).\n\n")

	bands := r.Bands
	if len(bands) == 0 {
		bands = []float64{100e6, 433e6, 915e6, 2.45e9}
	}
	b.WriteString("## Radio output (Nyquist into 50 Ω)\n\n")
	b.WriteString("Open-circuit voltage PSD is 4kT Re(Z). Available power accounts for mismatch to 50 Ω. ")
	b.WriteString("HackRF noise is roughly −170 dBm/Hz; this source is many tens of dB below that unless a near-field probe is pressed to the E.\n\n")
	b.WriteString("| f | Re(Z) | Im(Z) | S<sub>v</sub> (open) | P(50 Ω) |\n")
	b.WriteString("| ---: | ---: | ---: | ---: | ---: |\n")
	for _, f := range bands {
		s := d.Nyquist(f)
		fmt.Fprintf(&b, "| %s | %.2f mΩ | %.2f mΩ | %.2e V²/Hz | **%.1f dBm/Hz** |\n",
			fmtHz(f), real(s.Z)*1e3, imag(s.Z)*1e3, s.VocHz, s.P50dBmHz)
	}
	b.WriteString("\n")

	sweep := d.Sweep(1e6, 6e9, 96)
	b.WriteString("Available power vs frequency (dBm/Hz):\n\n")
	b.WriteString("```\n")
	b.WriteString(plotSweep(sweep, 72, 14))
	b.WriteString("```\n\n")

	b.WriteString("## What to look for on the HackRFs\n\n")
	b.WriteString("1. **No sharp Casimir spurs.** The 1 µm MIM is overdamped; thermal output is a smooth continuum tens of dB below HackRF noise.\n")
	b.WriteString("2. **Near-field only.** A probe on the E would see milliohm-source Johnson noise; free-space coupling between the two radios will not show this sandwich.\n")
	b.WriteString("3. **Mechanical / patch-potential noise around the plate fundamental**, which is kHz, not a spectral-difference peak in the 100 MHz–2.45 GHz captures.\n")
	b.WriteString("4. Differences already seen between the two HackRFs (LO leakage, +1 MHz at 433 MHz, ±2/3 MHz at 915 MHz) are **receiver fingerprints**, not this sandwich.\n")
	return b.String()
}

func schematic() string {
	return strings.Join([]string{
		"   anodized Al plate",
		"  +------------------------------+",
		"  | 1 µm Al2O3                   |",
		"  |   ########################   |",
		"  |   ####                       |",
		"  |   ########################   |  E-shaped Al",
		"  |   ####                       |",
		"  |   ########################   |",
		"  | 1 µm Al2O3                   |",
		"  +------------------------------+",
		"   anodized Al plate",
	}, "\n") + "\n"
}

func fmtHz(hz float64) string {
	a := math.Abs(hz)
	switch {
	case a >= 0.999e9:
		return fmt.Sprintf("%.2f GHz", hz/1e9)
	case a >= 0.999e6:
		return fmt.Sprintf("%.2f MHz", hz/1e6)
	case a >= 0.999e3:
		return fmt.Sprintf("%.2f kHz", hz/1e3)
	default:
		return fmt.Sprintf("%.2f Hz", hz)
	}
}

func plotSweep(s []Spectrum, width, height int) string {
	if width < 16 {
		width = 16
	}
	if height < 5 {
		height = 5
	}
	vals := make([]float64, len(s))
	minv, maxv := 1e9, -1e9
	for i, p := range s {
		v := p.P50dBmHz
		vals[i] = v
		if v < minv {
			minv = v
		}
		if v > maxv {
			maxv = v
		}
	}
	if maxv-minv < 1 {
		maxv = minv + 1
	}
	cols := make([]float64, width)
	for i := range cols {
		cols[i] = minv
	}
	for i, v := range vals {
		c := i * width / len(vals)
		if c >= width {
			c = width - 1
		}
		if v > cols[c] {
			cols[c] = v
		}
	}
	grid := make([][]byte, height)
	for y := 0; y < height; y++ {
		row := make([]byte, width)
		for x := range row {
			row[x] = ' '
		}
		grid[y] = row
	}
	span := maxv - minv
	zero := height - 1
	for x, v := range cols {
		y := int(math.Round(float64(zero) * (maxv - v) / span))
		if y < 0 {
			y = 0
		}
		if y > zero {
			y = zero
		}
		grid[y][x] = '*'
		for yy := y + 1; yy <= zero; yy++ {
			if grid[yy][x] == ' ' {
				grid[yy][x] = '|'
			}
		}
	}
	var b strings.Builder
	for y := 0; y < height; y++ {
		label := "        "
		if y == 0 {
			label = fmt.Sprintf("%+6.0f ", maxv)
		} else if y == height-1 {
			label = fmt.Sprintf("%+6.0f ", minv)
		}
		b.WriteString(label)
		b.Write(grid[y])
		b.WriteByte('\n')
	}
	lo := fmtHz(s[0].Freq)
	hi := fmtHz(s[len(s)-1].Freq)
	pad := width - len(lo) - len(hi)
	if pad < 1 {
		pad = 1
	}
	fmt.Fprintf(&b, "        %s%s%s\n", lo, strings.Repeat(" ", pad), hi)
	b.WriteString("        dBm/Hz into 50 Ω, log-frequency\n")
	return b.String()
}
