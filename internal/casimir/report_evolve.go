package casimir

import (
	"fmt"
	"strings"
	"time"
)

// EvolveReport is the markdown write-up of a genetic search.
type EvolveReport struct {
	Time     time.Time
	Config   EvolveConfig
	Best     Individual
	Baseline Individual
}

func (r EvolveReport) Markdown() string {
	g := r.Config.grid(r.Best)
	var b strings.Builder
	fmt.Fprintf(&b, "# Evolved Casimir center sheet\n\n")
	fmt.Fprintf(&b, "Evolved %s.\n\n", r.Time.UTC().Format("2006-01-02 15:04:05 UTC"))
	b.WriteString("A genetic algorithm searched binary metal occupancy on a rectangular lattice ")
	b.WriteString("lying flat between the two 1 µm-anodized plates. Ranking fitness is ")
	b.WriteString("**⟨P₅₀⟩·A** (mean available thermal power into 50 Ω times metal area) so the search ")
	b.WriteString("does not collapse to a single cell: raw P₅₀ rises as C falls, while Casimir force scales with area. ")
	b.WriteString("Power is evaluated at 100 MHz, 433 MHz, 915 MHz, and 2.45 GHz. Span is a gene. ")
	b.WriteString("The E-shaped seed is the baseline.\n\n")
	fmt.Fprintf(&b, "| | |\n| --- | ---: |\n")
	fmt.Fprintf(&b, "| Grid | %d × %d |\n", r.Config.Rows, r.Config.Cols)
	fmt.Fprintf(&b, "| Population / generations | %d / %d |\n", r.Config.Pop, r.Config.Gen)
	fmt.Fprintf(&b, "| Mutation | %.0f%% bits / generation |\n", 100*r.Config.MutP)
	fmt.Fprintf(&b, "| Span range | %.1f–%.0f mm |\n\n", r.Config.MinSpan*1e3, r.Config.MaxSpan*1e3)

	b.WriteString("## Best sheet\n\n")
	fmt.Fprintf(&b, "| | |\n| --- | ---: |\n")
	fmt.Fprintf(&b, "| Span | **%.2f mm** |\n", r.Best.Span*1e3)
	fmt.Fprintf(&b, "| Fill | **%.1f%%** (%d cells) |\n", 100*r.Best.Fill, g.nMetal())
	fmt.Fprintf(&b, "| Metal area | %.3f mm² |\n", g.Area()*1e6)
	fmt.Fprintf(&b, "| Mean P(50 Ω) | **%.2f dBm/Hz** |\n", r.Best.P50dBm)
	fmt.Fprintf(&b, "| Fitness ⟨P⟩·A | **%.3e** W·m²/Hz |\n", r.Best.Fit)
	if r.Baseline.Fit > 0 {
		fmt.Fprintf(&b, "| E-shape baseline P | %.2f dBm/Hz |\n", r.Baseline.P50dBm)
		fmt.Fprintf(&b, "| E-shape baseline ⟨P⟩·A | %.3e W·m²/Hz |\n", r.Baseline.Fit)
		fmt.Fprintf(&b, "| Fitness gain | **×%.2f** |\n", r.Best.Fit/r.Baseline.Fit)
	}
	b.WriteString("\n")
	b.WriteString("Band-by-band available power:\n\n")
	b.WriteString("| f | Re(Z) | Im(Z) | P(50 Ω) |\n")
	b.WriteString("| ---: | ---: | ---: | ---: |\n")
	for _, f := range r.Config.Freqs {
		s, ok := g.Nyquist(f)
		if !ok {
			fmt.Fprintf(&b, "| %s | — | — | — |\n", fmtHz(f))
			continue
		}
		fmt.Fprintf(&b, "| %s | %.3g Ω | %.3g Ω | **%.1f dBm/Hz** |\n",
			fmtHz(f), real(s.Z), imag(s.Z), s.P50dBmHz)
	}
	b.WriteString("\nPlan view (`#` = aluminum, `.` = oxide only):\n\n```\n")
	b.WriteString(g.ASCII())
	b.WriteString("```\n\n")
	b.WriteString("![Evolved center sheet](sheet.png)\n\n")
	b.WriteString("![Evolved sheet, studio view](sheet_photo.jpg)\n")
	return b.String()
}
