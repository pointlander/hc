package spectrum

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
)

// Radio describes one capture path for the markdown report.
type Radio struct {
	Name      string
	Serial    string
	USBPath   string
	Firmware  string
	Samples   int
	ClipFrac  float64
	MeanI     float64
	MeanQ     float64
	PowerDBFS float64
}

// Band is one center-frequency capture from both radios.
type Band struct {
	CenterHz uint64
	RateHz   float64
	A        PSD
	B        PSD
	Diff     Diff
}

// Report is the full comparison written to markdown.
type Report struct {
	Time     time.Time
	RateHz   float64
	Duration time.Duration
	FFTSize  int
	Overlap  int
	LNAGain  int
	VGAGain  int
	Amp      bool
	NSigma   float64
	MinDB    float64
	A        Radio
	B        Radio
	Bands    []Band
}

// Markdown renders the comparison as a standalone markdown document.
func (r Report) Markdown() string {
	var b strings.Builder
	fmt.Fprintf(&b, "# HackRF spectral difference\n\n")
	fmt.Fprintf(&b, "Captured %s.\n\n", r.Time.UTC().Format("2006-01-02 15:04:05 UTC"))
	fmt.Fprintf(&b, "Two HackRF One radios received the same band with identical gain, sample rate, and baseband filter. ")
	fmt.Fprintf(&b, "Each capture is a Hann-windowed Welch PSD. A constant median offset (gain / noise-figure difference) is subtracted; ")
	fmt.Fprintf(&b, "bins whose residual exceeds **max(%.0f·σ<sub>MAD</sub>, %.1f dB)** are reported as significant.\n\n", r.NSigma, r.MinDB)

	b.WriteString("## Radios\n\n")
	b.WriteString("| | Label | Serial | USB | Firmware | Samples | Clip | DC (I, Q) | Power |\n")
	b.WriteString("| --- | --- | --- | --- | --- | ---: | ---: | --- | ---: |\n")
	writeRadioRow(&b, "A", r.A)
	writeRadioRow(&b, "B", r.B)
	b.WriteString("\n")

	b.WriteString("## Setup\n\n")
	fmt.Fprintf(&b, "- Sample rate: %.3f MS/s\n", r.RateHz/1e6)
	fmt.Fprintf(&b, "- Capture: %s per band (first 50 ms discarded)\n", r.Duration)
	fmt.Fprintf(&b, "- FFT: %d, Hann, overlap %d samples (%.0f%%)\n", r.FFTSize, r.Overlap, 100*float64(r.Overlap)/float64(r.FFTSize))
	fmt.Fprintf(&b, "- RX gain: LNA %d dB, VGA %d dB, RF amp %v\n", r.LNAGain, r.VGAGain, r.Amp)
	fmt.Fprintf(&b, "- Baseband filter: nearest supported ≤ 75%% of sample rate\n")
	fmt.Fprintf(&b, "- Significance: residual ≥ max(%.0f σ, %.1f dB) after median offset; DC and analog-filter edges excluded from σ\n\n", r.NSigma, r.MinDB)

	b.WriteString("## Summary\n\n")
	writeSummary(&b, r)

	for _, band := range r.Bands {
		writeBand(&b, band)
	}
	return b.String()
}

func writeRadioRow(b *strings.Builder, tag string, r Radio) {
	usb := r.USBPath
	if usb == "" {
		usb = "—"
	}
	fw := r.Firmware
	if fw == "" {
		fw = "—"
	}
	fmt.Fprintf(b, "| %s | %s | `%s` | %s | %s | %d | %.3f%% | %+.4f, %+.4f | %.1f dBFS |\n",
		tag, r.Name, r.Serial, usb, fw, r.Samples, 100*r.ClipFrac, r.MeanI, r.MeanQ, r.PowerDBFS)
}

func writeSummary(b *strings.Builder, r Report) {
	if len(r.Bands) == 0 {
		b.WriteString("No bands captured.\n\n")
		return
	}
	var meanOff, meanNF, meanDC float64
	nPeaks := 0
	type hit struct {
		center uint64
		p      Peak
	}
	var top []hit
	for _, band := range r.Bands {
		meanOff += band.Diff.OffsetDB
		meanNF += band.Diff.NoiseADB - band.Diff.NoiseBDB
		meanDC += band.Diff.DCADB - band.Diff.DCBDB
		nPeaks += len(band.Diff.Peaks)
		for _, p := range band.Diff.Peaks {
			top = append(top, hit{center: band.CenterHz, p: p})
		}
	}
	n := float64(len(r.Bands))
	meanOff /= n
	meanNF /= n
	meanDC /= n
	hotter := "A"
	if meanOff < 0 {
		hotter = "B"
		meanOff = -meanOff
	}
	fmt.Fprintf(b, "- Across **%d** band(s), radio **%s** is **%.2f dB** hotter (median PSD offset).\n", len(r.Bands), hotter, meanOff)
	fmt.Fprintf(b, "- Mean noise-floor Δ(A−B): **%+.2f dB**.\n", meanNF)
	fmt.Fprintf(b, "- Mean LO leakage Δ(A−B): **%+.2f dB** (DC bin).\n", meanDC)
	fmt.Fprintf(b, "- **%d** significant bins in total.\n\n", nPeaks)

	if len(top) == 0 {
		b.WriteString("No residual exceeded the significance threshold. The radios match within the noise of the estimate, aside from the median offset above.\n\n")
		return
	}
	sort.Slice(top, func(i, j int) bool {
		return math.Abs(top[i].p.ResidualDB) > math.Abs(top[j].p.ResidualDB)
	})
	if len(top) > 8 {
		top = top[:8]
	}
	b.WriteString("Largest residuals:\n\n")
	b.WriteString("| Center | Offset | A | B | Δ | residual | σ |\n")
	b.WriteString("| ---: | ---: | ---: | ---: | ---: | ---: | ---: |\n")
	for _, h := range top {
		fmt.Fprintf(b, "| %s | %+s | %.1f dBFS | %.1f dBFS | %+.2f dB | %+.2f dB | %+.1f |\n",
			fmtHz(float64(h.center)), fmtHz(h.p.FreqHz), h.p.ADB, h.p.BDB, h.p.DeltaDB, h.p.ResidualDB, h.p.Sigma)
	}
	b.WriteString("\n")
}

func writeBand(b *strings.Builder, band Band) {
	d := band.Diff
	fmt.Fprintf(b, "## Band %s\n\n", fmtHz(float64(band.CenterHz)))
	fmt.Fprintf(b, "- Welch windows: A %d, B %d\n", band.A.Windows, band.B.Windows)
	fmt.Fprintf(b, "- Median offset A−B: **%+.2f dB**\n", d.OffsetDB)
	fmt.Fprintf(b, "- Residual σ (MAD): **%.3f dB**; threshold **%.2f dB**\n", d.SigmaDB, d.ThresholdDB)
	fmt.Fprintf(b, "- Noise floor: A **%.1f dBFS**, B **%.1f dBFS** (Δ %+.2f dB)\n", d.NoiseADB, d.NoiseBDB, d.NoiseADB-d.NoiseBDB)
	fmt.Fprintf(b, "- LO leakage (DC): A **%.1f dBFS**, B **%.1f dBFS** (Δ %+.2f dB, residual %+0.2f dB)\n",
		d.DCADB, d.DCBDB, d.DCADB-d.DCBDB, d.DCADB-d.DCBDB-d.OffsetDB)
	fmt.Fprintf(b, "- RMS residual: **%.2f dB**\n", d.RMSResidual)
	fmt.Fprintf(b, "- Significant bins: **%d**\n\n", len(d.Peaks))

	if len(d.Peaks) > 0 {
		limit := len(d.Peaks)
		if limit > 24 {
			limit = 24
		}
		b.WriteString("| Offset | A | B | Δ | residual | σ |\n")
		b.WriteString("| ---: | ---: | ---: | ---: | ---: | ---: |\n")
		for _, p := range d.Peaks[:limit] {
			fmt.Fprintf(b, "| %+s | %.1f dBFS | %.1f dBFS | %+.2f dB | %+.2f dB | %+.1f |\n",
				fmtHz(p.FreqHz), p.ADB, p.BDB, p.DeltaDB, p.ResidualDB, p.Sigma)
		}
		if len(d.Peaks) > limit {
			fmt.Fprintf(b, "\n_%d additional significant bins omitted._\n", len(d.Peaks)-limit)
		}
		b.WriteString("\n")
	}

	b.WriteString("Residual spectrum (dB, median offset removed). `A` is hotter above the axis.\n\n")
	b.WriteString("```\n")
	b.WriteString(ASCIIPlot(band.A, d.ResidualDB, 72, 17, d.ThresholdDB))
	b.WriteString("```\n\n")
}

// ASCIIPlot draws a downsampled residual spectrum. DC is the center column.
func ASCIIPlot(p PSD, residual []float64, width, height int, thr float64) string {
	if width < 16 {
		width = 16
	}
	if height < 5 {
		height = 5
	}
	n := min(p.FFTSize, len(residual))
	if n == 0 {
		return ""
	}
	// Order bins from -Fs/2 .. +Fs/2 and drop analog-filter edges so
	// baseband roll-off does not set the plot scale.
	ordered := make([]float64, n)
	half := n / 2
	copy(ordered[0:half], residual[n-half:])
	copy(ordered[half:], residual[0:n-half])
	trim := n / 10
	if trim < 4 {
		trim = 4
	}
	loHz := -p.SampleRate / 2
	hiHz := p.SampleRate / 2
	if n > 2*trim+8 {
		ordered = ordered[trim : n-trim]
		loHz = -p.SampleRate/2 + float64(trim)*p.SampleRate/float64(n)
		hiHz = p.SampleRate/2 - float64(trim)*p.SampleRate/float64(n)
	}

	cols := make([]float64, width)
	nOrd := len(ordered)
	for i, v := range ordered {
		c := i * width / nOrd
		if c >= width {
			c = width - 1
		}
		if math.Abs(v) > math.Abs(cols[c]) {
			cols[c] = v
		}
	}

	ymax := thr * 2
	for _, v := range cols {
		if math.Abs(v) > ymax {
			ymax = math.Abs(v)
		}
	}
	if ymax < 1 {
		ymax = 1
	}

	grid := make([][]byte, height)
	for y := 0; y < height; y++ {
		row := make([]byte, width)
		for x := range row {
			row[x] = ' '
		}
		grid[y] = row
	}
	zero := height / 2
	for x := 0; x < width; x++ {
		grid[zero][x] = '-'
	}
	// threshold lines
	thrRow := int(math.Round(float64(zero) - thr/ymax*float64(zero)))
	nthrRow := int(math.Round(float64(zero) + thr/ymax*float64(zero)))
	for _, y := range []int{thrRow, nthrRow} {
		if y >= 0 && y < height && y != zero {
			for x := 0; x < width; x++ {
				if grid[y][x] == ' ' {
					grid[y][x] = '.'
				}
			}
		}
	}
	for x, v := range cols {
		y := int(math.Round(float64(zero) - v/ymax*float64(zero)))
		if y < 0 {
			y = 0
		}
		if y > height-1 {
			y = height - 1
		}
		mark := byte('*')
		if math.Abs(v) >= thr {
			mark = '#'
		}
		grid[y][x] = mark
		if y == zero {
			continue
		}
		step := 1
		if y > zero {
			step = -1
		}
		for yy := y + step; yy != zero && yy >= 0 && yy < height; yy += step {
			if grid[yy][x] == ' ' || grid[yy][x] == '.' {
				grid[yy][x] = '|'
			}
		}
	}

	var b strings.Builder
	for y := 0; y < height; y++ {
		label := "      "
		if y == 0 {
			label = fmt.Sprintf("%+5.1f ", ymax)
		} else if y == zero {
			label = "  0.0 "
		} else if y == height-1 {
			label = fmt.Sprintf("%+5.1f ", -ymax)
		}
		b.WriteString(label)
		b.Write(grid[y])
		b.WriteByte('\n')
	}
	lo := fmtHz(loHz)
	mid := "0"
	hi := fmtHz(hiHz)
	pad := width - len(lo) - len(mid) - len(hi)
	if pad < 2 {
		pad = 2
	}
	left := pad / 2
	right := pad - left
	fmt.Fprintf(&b, "      %s%s%s%s%s\n", lo, strings.Repeat(" ", left), mid, strings.Repeat(" ", right), hi)
	return b.String()
}

func fmtHz(hz float64) string {
	sign := ""
	if hz < 0 {
		sign = "-"
		hz = -hz
	} else if hz > 0 {
		sign = "+"
	}
	switch {
	case hz >= 1e9:
		return fmt.Sprintf("%s%.6g GHz", sign, hz/1e9)
	case hz >= 1e6:
		return fmt.Sprintf("%s%.6g MHz", sign, hz/1e6)
	case hz >= 1e3:
		return fmt.Sprintf("%s%.4g kHz", sign, hz/1e3)
	case hz == 0:
		return "0 Hz"
	default:
		return fmt.Sprintf("%s%.4g Hz", sign, hz)
	}
}
