package casimir

import "math"

// Grid is a coplanar metal occupancy on a rectangular lattice lying
// flat between the two anodized plates. Empty cells are the 1 µm oxide
// facing itself; metal cells are the center sheet.
type Grid struct {
	Rows, Cols int
	Span       float64 // bounding-box width (= height if square cells), m
	Metal      []bool  // row-major, length Rows*Cols
	Mat        Device
}

func (g Grid) cellW() float64 {
	if g.Cols == 0 {
		return 0
	}
	return g.Span / float64(g.Cols)
}

func (g Grid) cellH() float64 {
	if g.Rows == 0 {
		return 0
	}
	hspan := g.Span * float64(g.Rows) / float64(g.Cols)
	return hspan / float64(g.Rows)
}

func (g Grid) at(r, c int) bool {
	if r < 0 || c < 0 || r >= g.Rows || c >= g.Cols {
		return false
	}
	return g.Metal[r*g.Cols+c]
}

func (g Grid) nMetal() int {
	n := 0
	for _, m := range g.Metal {
		if m {
			n++
		}
	}
	return n
}

// Area is the metal area in m².
func (g Grid) Area() float64 {
	return float64(g.nMetal()) * g.cellW() * g.cellH()
}

// Prune removes metal that is not 4-connected to the centroid feed.
func (g *Grid) Prune() {
	idx, list := g.nodes()
	feed := g.feedIndex(list)
	keep := g.connected(list, idx, feed)
	for i := range g.Metal {
		g.Metal[i] = false
	}
	for i, on := range keep {
		if on {
			n := list[i]
			g.Metal[n.r*g.Cols+n.c] = true
		}
	}
}

// FillE paints the Device's E (tines up, spine along the bottom) onto g.
func (g *Grid) FillE(d Device) {
	for i := range g.Metal {
		g.Metal[i] = false
	}
	if d.Width <= 0 || d.Height <= 0 {
		return
	}
	for r := 0; r < g.Rows; r++ {
		for c := 0; c < g.Cols; c++ {
			x := (float64(c) + 0.5) / float64(g.Cols) * d.Width
			y := (float64(g.Rows-1-r) + 0.5) / float64(g.Rows) * d.Height // row 0 is top
			g.Metal[r*g.Cols+c] = eOccupied(d, x, y)
		}
	}
	g.Span = math.Max(d.Width, d.Height)
}

func eOccupied(d Device, x, y float64) bool {
	// Spine along y=0..Spine, x=0..Width; tines of width Arm going up.
	if y >= 0 && y <= d.Spine && x >= 0 && x <= d.Width {
		return true
	}
	slot := d.slot()
	armLen := d.armLen()
	if y < d.Spine || y > d.Spine+armLen {
		return false
	}
	for i := 0; i < 3; i++ {
		x0 := float64(i) * (d.Arm + slot)
		if x >= x0 && x <= x0+d.Arm {
			return true
		}
	}
	return false
}

type node struct{ r, c int }

func (g Grid) nodes() (idx map[int]int, list []node) {
	idx = make(map[int]int)
	for r := 0; r < g.Rows; r++ {
		for c := 0; c < g.Cols; c++ {
			if !g.at(r, c) {
				continue
			}
			k := r*g.Cols + c
			idx[k] = len(list)
			list = append(list, node{r, c})
		}
	}
	return idx, list
}

func (g Grid) feedIndex(list []node) int {
	if len(list) == 0 {
		return -1
	}
	cr := float64(g.Rows-1) / 2
	cc := float64(g.Cols-1) / 2
	best, bestD := 0, 1e18
	for i, n := range list {
		dr := float64(n.r) - cr
		dc := float64(n.c) - cc
		d := dr*dr + dc*dc
		if d < bestD {
			best, bestD = i, d
		}
	}
	return best
}

func (g Grid) connected(list []node, idx map[int]int, feed int) []bool {
	keep := make([]bool, len(list))
	if feed < 0 || feed >= len(list) {
		return keep
	}
	q := []int{feed}
	keep[feed] = true
	dirs := [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	for len(q) > 0 {
		i := q[0]
		q = q[1:]
		n := list[i]
		for _, d := range dirs {
			rr, cc := n.r+d[0], n.c+d[1]
			if !g.at(rr, cc) {
				continue
			}
			j, ok := idx[rr*g.Cols+cc]
			if !ok || keep[j] {
				continue
			}
			keep[j] = true
			q = append(q, j)
		}
	}
	return keep
}

// Impedance is the Thevenin Z of the metal island containing the
// centroid feed, from a nodal RLGC mesh (both plates as ground).
func (g Grid) Impedance(freq float64) (complex128, bool) {
	if freq < 1 {
		freq = 1
	}
	idx, list := g.nodes()
	if len(list) == 0 {
		return 0, false
	}
	feed := g.feedIndex(list)
	keep := g.connected(list, idx, feed)
	var used []node
	old := make([]int, len(list))
	for i, n := range list {
		old[i] = -1
		if keep[i] {
			old[i] = len(used)
			used = append(used, n)
		}
	}
	n := len(used)
	if n == 0 {
		return 0, false
	}
	feed = old[feed]
	omega := 2 * math.Pi * freq
	cw, ch := g.cellW(), g.cellH()
	mat := g.Mat
	cCell := 2 * eps0 * mat.EpsR * cw * ch / mat.Oxide
	gCell := omega * cCell * mat.TanD
	delta := math.Sqrt(2 / (omega * mu0 * mat.Sigma))
	skin := delta
	if mat.Thick > 0 && mat.Thick < skin {
		skin = mat.Thick
	}
	rsq := 1 / (mat.Sigma * skin)

	y := make([][]complex128, n)
	for i := 0; i < n; i++ {
		y[i] = make([]complex128, n)
		y[i][i] = complex(gCell, omega*cCell)
	}
	pos := make(map[int]int, n)
	for i, nd := range used {
		pos[nd.r*g.Cols+nd.c] = i
	}
	addEdge := func(i, j int, length, width float64) {
		if i > j {
			return
		}
		r := rsq * (length / width)
		l := mu0 * mat.Oxide / 2 * (length / width)
		z := complex(r, omega*l)
		if z == 0 {
			return
		}
		yy := 1 / z
		y[i][i] += yy
		y[j][j] += yy
		y[i][j] -= yy
		y[j][i] -= yy
	}
	for i, nd := range used {
		if j, ok := pos[(nd.r)*g.Cols+(nd.c+1)]; ok {
			addEdge(i, j, cw, ch)
		}
		if j, ok := pos[(nd.r+1)*g.Cols+nd.c]; ok {
			addEdge(i, j, ch, cw)
		}
	}
	rhs := make([]complex128, n)
	rhs[feed] = 1
	v, ok := solveComplex(y, rhs)
	if !ok {
		return 0, false
	}
	return v[feed], true
}

// Nyquist is available thermal power of the grid sheet into Mat.Zload.
func (g Grid) Nyquist(freq float64) (Spectrum, bool) {
	z, ok := g.Impedance(freq)
	if !ok {
		return Spectrum{Freq: freq, P50dBmHz: -400}, false
	}
	re := real(z)
	if re < 0 {
		re = 0
	}
	voc := 4 * kB * g.Mat.Temp * re
	zl := g.Mat.Zload
	h := complex(zl, 0) / (z + complex(zl, 0))
	p := voc * abs2(h) / zl
	return Spectrum{Freq: freq, Z: z, VocHz: voc, P50: p, P50dBmHz: dbmHz(p)}, true
}

// MeanPower is the mean available W/Hz over freqs. Zero if the sheet is open.
func (g Grid) MeanPower(freqs []float64) float64 {
	if len(freqs) == 0 {
		return 0
	}
	var s float64
	n := 0
	for _, f := range freqs {
		sp, ok := g.Nyquist(f)
		if !ok {
			return 0
		}
		s += sp.P50
		n++
	}
	if n == 0 {
		return 0
	}
	return s / float64(n)
}
