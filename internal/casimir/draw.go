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
	metal := color.NRGBA{196, 202, 210, 255}
	hilite := color.NRGBA{232, 236, 240, 255}
	shade := color.NRGBA{148, 154, 162, 255}
	edge := color.NRGBA{120, 128, 138, 255}
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
			col := oxide
			if g.at(r, c) {
				fillRect(img, x0, y0, x0+cw-1, y0+ch-1, metal)
				if cw > 3 && ch > 3 {
					fillRect(img, x0, y0, x0+cw-2, y0, hilite)
					fillRect(img, x0, y0, x0, y0+ch-2, hilite)
					fillRect(img, x0+1, y0+ch-1, x0+cw-1, y0+ch-1, shade)
					fillRect(img, x0+cw-1, y0+1, x0+cw-1, y0+ch-1, shade)
				}
				continue
			}
			fillRect(img, x0, y0, x0+cw-1, y0+ch-1, col)
		}
	}
	// plate frame
	frame(img, margin-4, margin-4, px-margin+3, px-margin+3, edge)
	return png.Encode(w, img)
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

// ASCII is a text plan view of the metal (█) on empty oxide (·).
func (g Grid) ASCII() string {
	buf := make([]byte, 0, g.Rows*(g.Cols+1))
	for r := 0; r < g.Rows; r++ {
		for c := 0; c < g.Cols; c++ {
			if g.at(r, c) {
				buf = append(buf, '#')
			} else {
				buf = append(buf, '.')
			}
		}
		buf = append(buf, '\n')
	}
	return string(buf)
}
