package casimir

import (
	"image"
	"image/color"
	"image/png"
	"io"
)

// WritePNG paints a top-down view of the center sheet on the anodized plate.
func (g Grid) WritePNG(w io.Writer, px int) error {
	if px < 64 {
		px = 512
	}
	margin := px / 12
	inner := px - 2*margin
	img := image.NewNRGBA(image.Rect(0, 0, px, px))
	plate := color.NRGBA{36, 40, 46, 255}
	oxide := color.NRGBA{58, 64, 72, 255}
	outside := color.NRGBA{28, 30, 34, 255}
	metal := color.NRGBA{196, 202, 210, 255}
	overhang := color.NRGBA{210, 186, 140, 255} // metal past the plates
	hilite := color.NRGBA{232, 236, 240, 255}
	ohilite := color.NRGBA{236, 214, 170, 255}
	shade := color.NRGBA{148, 154, 162, 255}
	oshade := color.NRGBA{160, 132, 90, 255}
	edge := color.NRGBA{120, 128, 138, 255}
	plateEdge := color.NRGBA{212, 168, 72, 255}
	bg := color.NRGBA{18, 20, 24, 255}
	fillRect(img, 0, 0, px, px, bg)
	fillRect(img, margin-4, margin-4, px-margin+4, px-margin+4, plate)
	cw := inner / g.Cols
	ch := inner / g.Rows
	if cw < 1 {
		cw = 1
	}
	if ch < 1 {
		ch = 1
	}
	ox := margin + (inner-cw*g.Cols)/2
	oy := margin + (inner-ch*g.Rows)/2
	for r := 0; r < g.Rows; r++ {
		for c := 0; c < g.Cols; c++ {
			x0 := ox + c*cw
			y0 := oy + r*ch
			onPlate := g.covered(r, c)
			if g.at(r, c) {
				fill, hi, sh := metal, hilite, shade
				if !onPlate {
					fill, hi, sh = overhang, ohilite, oshade
				}
				fillRect(img, x0, y0, x0+cw-1, y0+ch-1, fill)
				if cw > 3 && ch > 3 {
					fillRect(img, x0, y0, x0+cw-2, y0, hi)
					fillRect(img, x0, y0, x0, y0+ch-2, hi)
					fillRect(img, x0+1, y0+ch-1, x0+cw-1, y0+ch-1, sh)
					fillRect(img, x0+cw-1, y0+1, x0+cw-1, y0+ch-1, sh)
				}
				continue
			}
			col := oxide
			if !onPlate {
				col = outside
			}
			fillRect(img, x0, y0, x0+cw-1, y0+ch-1, col)
		}
	}
	frame(img, margin-4, margin-4, px-margin+3, px-margin+3, edge)
	px0, py0, px1, py1 := platePixels(g, ox, oy, cw, ch)
	frame(img, px0, py0, px1, py1, plateEdge)
	return png.Encode(w, img)
}

func platePixels(g Grid, ox, oy, cw, ch int) (x0, y0, x1, y1 int) {
	sw, sh := g.sheetSize()
	ps := g.effectivePlateSpan()
	if sw <= 0 || sh <= 0 {
		return ox, oy, ox, oy
	}
	if ps > sw {
		ps = sw
	}
	sheetW := g.Cols * cw
	sheetH := g.Rows * ch
	pw := int(ps / sw * float64(sheetW))
	ph := int(ps / sh * float64(sheetH))
	if pw < 1 {
		pw = 1
	}
	if ph < 1 {
		ph = 1
	}
	x0 = ox + (sheetW-pw)/2
	y0 = oy + (sheetH-ph)/2
	return x0, y0, x0 + pw - 1, y0 + ph - 1
}

func fillRect(img *image.NRGBA, x0, y0, x1, y1 int, col color.NRGBA) {
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	b := img.Bounds()
	for y := y0; y <= y1; y++ {
		if y < b.Min.Y || y >= b.Max.Y {
			continue
		}
		for x := x0; x <= x1; x++ {
			if x < b.Min.X || x >= b.Max.X {
				continue
			}
			img.SetNRGBA(x, y, col)
		}
	}
}

func frame(img *image.NRGBA, x0, y0, x1, y1 int, col color.NRGBA) {
	fillRect(img, x0, y0, x1, y0, col)
	fillRect(img, x0, y1, x1, y1, col)
	fillRect(img, x0, y0, x0, y1, col)
	fillRect(img, x1, y0, x1, y1, col)
}

// ASCII is a text plan view: '#' metal under the plates, '+' overhang, '.' empty.
func (g Grid) ASCII() string {
	buf := make([]byte, 0, g.Rows*(g.Cols+1))
	for r := 0; r < g.Rows; r++ {
		for c := 0; c < g.Cols; c++ {
			switch {
			case g.at(r, c) && g.covered(r, c):
				buf = append(buf, '#')
			case g.at(r, c):
				buf = append(buf, '+')
			default:
				buf = append(buf, '.')
			}
		}
		buf = append(buf, '\n')
	}
	return string(buf)
}
