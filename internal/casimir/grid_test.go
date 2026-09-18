package casimir

import (
	"bytes"
	"math"
	"math/rand"
	"strings"
	"testing"
)

func TestFullGridCapacitance(t *testing.T) {
	d := DefaultDevice()
	n := 8
	g := Grid{Rows: n, Cols: n, Span: 0.04, Metal: make([]bool, n*n), Mat: d}
	for i := range g.Metal {
		g.Metal[i] = true
	}
	z, ok := g.Impedance(1e5)
	if !ok {
		t.Fatal("impedance")
	}
	cEst := -1 / (2 * math.Pi * 1e5 * imag(z))
	cWant := 2 * eps0 * d.EpsR * g.Area() / d.Oxide
	rel := math.Abs(cEst-cWant) / cWant
	if rel > 0.2 {
		t.Fatalf("C_est=%g C=%g rel=%g Z=%v", cEst, cWant, rel, z)
	}
}

func TestEmptyGridNoZ(t *testing.T) {
	d := DefaultDevice()
	g := Grid{Rows: 8, Cols: 8, Span: 0.04, Metal: make([]bool, 64), Mat: d}
	if _, ok := g.Impedance(915e6); ok {
		t.Fatal("empty sheet should not have Z")
	}
}

func TestFillEHasTinesAndSpine(t *testing.T) {
	d := DefaultDevice()
	g := Grid{Rows: 16, Cols: 16, Metal: make([]bool, 256), Mat: d}
	g.FillE(d)
	if g.nMetal() < 20 {
		t.Fatalf("metal cells=%d", g.nMetal())
	}
	// spine is the bottom rows: last few rows should have more metal
	bottom := 0
	top := 0
	for c := 0; c < 16; c++ {
		if g.at(15, c) {
			bottom++
		}
		if g.at(0, c) {
			top++
		}
	}
	if bottom < top {
		t.Fatalf("spine should sit at the bottom of the plan view; bottom=%d top=%d", bottom, top)
	}
}

func TestPNGHeader(t *testing.T) {
	d := DefaultDevice()
	g := Grid{Rows: 8, Cols: 8, Span: 0.04, Metal: make([]bool, 64), Mat: d}
	g.Metal[0] = true
	var buf bytes.Buffer
	if err := g.WritePNG(&buf, 128); err != nil {
		t.Fatal(err)
	}
	b := buf.Bytes()
	if len(b) < 8 || string(b[:8]) != "\x89PNG\r\n\x1a\n" {
		t.Fatal("not a PNG")
	}
}

func TestEvolveBeatsEmpty(t *testing.T) {
	cfg := DefaultEvolve()
	cfg.Rows, cfg.Cols = 8, 8
	cfg.Pop, cfg.Gen = 8, 4
	cfg.Freqs = []float64{915e6}
	best := Evolve(cfg, rand.New(rand.NewSource(1)))
	if best.Power <= 0 || best.Fit <= 0 {
		t.Fatalf("power=%g fit=%g", best.Power, best.Fit)
	}
	if best.Fill <= 0 {
		t.Fatal("no metal")
	}
}

func TestPruneDropsIslands(t *testing.T) {
	g := Grid{Rows: 4, Cols: 4, Span: 0.04, Metal: make([]bool, 16), Mat: DefaultDevice()}
	g.Metal[0] = true                                     // island
	g.Metal[5], g.Metal[6], g.Metal[9] = true, true, true // blob near center
	g.Prune()
	if g.at(0, 0) {
		t.Fatal("corner island should be removed")
	}
	if g.nMetal() < 2 {
		t.Fatal("feed blob removed")
	}
}

func TestOverhangCapacitanceIsPlateArea(t *testing.T) {
	d := DefaultDevice()
	n := 8
	plate := 0.04
	full := Grid{Rows: n, Cols: n, Span: plate, Metal: make([]bool, n*n), Mat: d}
	big := Grid{Rows: n, Cols: n, Span: 2 * plate, PlateSpan: plate, Metal: make([]bool, n*n), Mat: d}
	for i := range full.Metal {
		full.Metal[i] = true
		big.Metal[i] = true
	}
	zFull, ok := full.Impedance(1e5)
	if !ok {
		t.Fatal("full")
	}
	zBig, ok := big.Impedance(1e5)
	if !ok {
		t.Fatal("big")
	}
	cFull := -1 / (2 * math.Pi * 1e5 * imag(zFull))
	cBig := -1 / (2 * math.Pi * 1e5 * imag(zBig))
	cPlate := 2 * eps0 * d.EpsR * plate * plate / d.Oxide
	if math.Abs(cFull-cPlate)/cPlate > 0.25 {
		t.Fatalf("flush C=%g want ~%g", cFull, cPlate)
	}
	if math.Abs(cBig-cPlate)/cPlate > 0.25 {
		t.Fatalf("overhang C=%g want ~plate %g (got sheet area would be 4×)", cBig, cPlate)
	}
	if big.OverlapArea() > big.Area()*0.4 {
		t.Fatalf("overlap %g should be ~1/4 of metal %g", big.OverlapArea(), big.Area())
	}
}

func TestOverhangASCII(t *testing.T) {
	g := Grid{Rows: 4, Cols: 4, Span: 0.08, PlateSpan: 0.04, Metal: make([]bool, 16), Mat: DefaultDevice()}
	for i := range g.Metal {
		g.Metal[i] = true
	}
	s := g.ASCII()
	if !strings.Contains(s, "+") || !strings.Contains(s, "#") {
		t.Fatalf("want both overlap and overhang marks:\n%s", s)
	}
}

func TestASCII(t *testing.T) {
	g := Grid{Rows: 2, Cols: 2, Metal: []bool{true, false, false, true}}
	s := g.ASCII()
	if !strings.Contains(s, "#") || !strings.Contains(s, ".") {
		t.Fatalf("%q", s)
	}
}
