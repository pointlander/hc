package casimir

import (
	"math"
	"math/rand"
	"sort"
)

// Individual is one candidate center-sheet: metal occupancy plus overall size.
type Individual struct {
	Metal  []bool
	Span   float64
	Power  float64 // mean available W/Hz
	Fit    float64 // Power·Area; ranking objective
	P50dBm float64
	Fill   float64
	Area   float64
}

// EvolveConfig controls the genetic search.
type EvolveConfig struct {
	Rows, Cols int
	Pop, Gen   int
	MutP       float64
	Elite      int
	MinSpan    float64
	MaxSpan    float64
	PlateSpan  float64 // fixed anodized plates; sheet Span may exceed this
	Freqs      []float64
	Mat        Device
	Log        func(gen int, best Individual)
}

func DefaultEvolve() EvolveConfig {
	return EvolveConfig{
		Rows:      16,
		Cols:      16,
		Pop:       40,
		Gen:       60,
		MutP:      0.04,
		Elite:     2,
		MinSpan:   2e-3,
		MaxSpan:   0.25,
		PlateSpan: 50e-3, // DefaultDevice sandwich; sheet may overhang
		Freqs:     []float64{100e6, 433e6, 915e6, 2.45e9},
		Mat:       DefaultDevice(),
	}
}

func (c EvolveConfig) grid(ind Individual) Grid {
	m := append([]bool(nil), ind.Metal...)
	return Grid{Rows: c.Rows, Cols: c.Cols, Span: ind.Span, PlateSpan: c.PlateSpan, Metal: m, Mat: c.Mat}
}

func (c EvolveConfig) Evaluate(ind *Individual) {
	g := c.grid(*ind)
	g.Prune()
	ind.Metal = g.Metal
	ind.Power = g.MeanPower(c.Freqs)
	ind.Area = g.Area()
	// Rank by thermal power times metal area so the search does not
	// collapse to a single cell (P50 rises as C falls, F_Casimir ∝ A).
	ind.Fit = ind.Power * ind.Area
	ind.P50dBm = dbmHz(ind.Power)
	n := 0
	for _, v := range ind.Metal {
		if v {
			n++
		}
	}
	ind.Fill = float64(n) / float64(len(ind.Metal))
}

// Evolve runs a generational GA and returns the best sheet.
func Evolve(cfg EvolveConfig, rng *rand.Rand) Individual {
	if rng == nil {
		rng = rand.New(rand.NewSource(1))
	}
	if cfg.Rows < 4 {
		cfg.Rows = 16
	}
	if cfg.Cols < 4 {
		cfg.Cols = 16
	}
	if cfg.Pop < 8 {
		cfg.Pop = 40
	}
	if cfg.Gen < 1 {
		cfg.Gen = 60
	}
	if cfg.Elite < 1 {
		cfg.Elite = 2
	}
	if cfg.MutP <= 0 {
		cfg.MutP = 0.04
	}
	if len(cfg.Freqs) == 0 {
		cfg.Freqs = DefaultEvolve().Freqs
	}
	if cfg.MinSpan <= 0 {
		cfg.MinSpan = 2e-3
	}
	if cfg.MaxSpan <= cfg.MinSpan {
		cfg.MaxSpan = 0.25
	}
	pop := make([]Individual, cfg.Pop)
	for i := range pop {
		pop[i] = cfg.seed(rng, i)
		cfg.Evaluate(&pop[i])
	}
	sortPop(pop)
	if cfg.Log != nil {
		cfg.Log(0, pop[0])
	}
	for g := 1; g <= cfg.Gen; g++ {
		next := make([]Individual, 0, cfg.Pop)
		for i := 0; i < cfg.Elite && i < len(pop); i++ {
			next = append(next, cloneInd(pop[i]))
		}
		for len(next) < cfg.Pop {
			a := tournament(pop, rng)
			b := tournament(pop, rng)
			child := cfg.crossover(a, b, rng)
			cfg.mutate(&child, rng)
			cfg.Evaluate(&child)
			next = append(next, child)
		}
		pop = next
		sortPop(pop)
		if cfg.Log != nil {
			cfg.Log(g, pop[0])
		}
	}
	return pop[0]
}

func sortPop(pop []Individual) {
	sort.Slice(pop, func(i, j int) bool { return pop[i].Fit > pop[j].Fit })
}

func cloneInd(in Individual) Individual {
	out := in
	out.Metal = append([]bool(nil), in.Metal...)
	return out
}

func tournament(pop []Individual, rng *rand.Rand) Individual {
	best := pop[rng.Intn(len(pop))]
	for k := 0; k < 2; k++ {
		c := pop[rng.Intn(len(pop))]
		if c.Fit > best.Fit {
			best = c
		}
	}
	return best
}

func (c EvolveConfig) seed(rng *rand.Rand, i int) Individual {
	n := c.Rows * c.Cols
	ind := Individual{Metal: make([]bool, n), Span: c.randSpan(rng)}
	switch i % 5 {
	case 0:
		g := Grid{Rows: c.Rows, Cols: c.Cols, Metal: ind.Metal, Mat: c.Mat}
		g.FillE(c.Mat)
		ind.Metal = g.Metal
		ind.Span = math.Max(c.Mat.Width, c.Mat.Height)
		if ind.Span < c.MinSpan {
			ind.Span = c.MinSpan
		}
	case 1:
		c.randomWalk(ind.Metal, rng, 12+rng.Intn(40))
	case 2:
		c.filledRect(ind.Metal, rng)
	case 3:
		c.meander(ind.Metal, rng)
	default:
		c.randomWalk(ind.Metal, rng, 8+rng.Intn(24))
		c.sparsify(ind.Metal, rng, 0.5)
	}
	return ind
}

func (c EvolveConfig) randSpan(rng *rand.Rand) float64 {
	u := rng.Float64()
	log0, log1 := math.Log(c.MinSpan), math.Log(c.MaxSpan)
	return math.Exp(log0 + u*(log1-log0))
}

func (c EvolveConfig) randomWalk(m []bool, rng *rand.Rand, steps int) {
	r, col := c.Rows/2, c.Cols/2
	dirs := [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	for i := 0; i < steps; i++ {
		m[r*c.Cols+col] = true
		d := dirs[rng.Intn(4)]
		nr, nc := r+d[0], col+d[1]
		if nr >= 0 && nr < c.Rows && nc >= 0 && nc < c.Cols {
			r, col = nr, nc
		}
	}
}

func (c EvolveConfig) filledRect(m []bool, rng *rand.Rand) {
	h := 1 + rng.Intn(c.Rows/2+1)
	w := 1 + rng.Intn(c.Cols/2+1)
	r0 := rng.Intn(c.Rows - h + 1)
	c0 := rng.Intn(c.Cols - w + 1)
	for r := r0; r < r0+h; r++ {
		for col := c0; col < c0+w; col++ {
			m[r*c.Cols+col] = true
		}
	}
}

func (c EvolveConfig) meander(m []bool, rng *rand.Rand) {
	r := 1 + rng.Intn(c.Rows-2)
	left := true
	for col := 0; col < c.Cols; col++ {
		m[r*c.Cols+col] = true
		if col%2 == 1 && r+1 < c.Rows && r-1 >= 0 {
			if left {
				m[(r+1)*c.Cols+col] = true
				r++
			} else {
				m[(r-1)*c.Cols+col] = true
				r--
			}
			left = !left
		}
	}
}

func (c EvolveConfig) sparsify(m []bool, rng *rand.Rand, p float64) {
	for i := range m {
		if m[i] && rng.Float64() < p {
			m[i] = false
		}
	}
}

func (c EvolveConfig) crossover(a, b Individual, rng *rand.Rand) Individual {
	child := Individual{Metal: make([]bool, len(a.Metal)), Span: a.Span}
	if rng.Float64() < 0.5 {
		child.Span = b.Span
	}
	if rng.Float64() < 0.5 {
		child.Span = math.Sqrt(a.Span * b.Span)
	}
	r0 := rng.Intn(c.Rows)
	c0 := rng.Intn(c.Cols)
	r1 := r0 + 1 + rng.Intn(c.Rows-r0)
	c1 := c0 + 1 + rng.Intn(c.Cols-c0)
	for r := 0; r < c.Rows; r++ {
		for col := 0; col < c.Cols; col++ {
			src := a.Metal
			if r >= r0 && r < r1 && col >= c0 && col < c1 {
				src = b.Metal
			}
			child.Metal[r*c.Cols+col] = src[r*c.Cols+col]
		}
	}
	return child
}

func (c EvolveConfig) mutate(ind *Individual, rng *rand.Rand) {
	for i := range ind.Metal {
		if rng.Float64() < c.MutP {
			ind.Metal[i] = !ind.Metal[i]
		}
	}
	ind.Span *= math.Exp(rng.NormFloat64() * 0.12)
	if ind.Span < c.MinSpan {
		ind.Span = c.MinSpan
	}
	if ind.Span > c.MaxSpan {
		ind.Span = c.MaxSpan
	}
}
