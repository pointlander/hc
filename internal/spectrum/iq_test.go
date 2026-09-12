package spectrum

import (
	"math"
	"testing"
)

func TestBytesToIQ(t *testing.T) {
	b := []byte{64, 192, 127, 0} // int8 64, -64, 127, 0
	z := BytesToIQ(b)
	if len(z) != 2 {
		t.Fatalf("len=%d", len(z))
	}
	if math.Abs(real(z[0])-0.5) > 1e-9 || math.Abs(imag(z[0])+0.5) > 1e-9 {
		t.Fatalf("z0=%v", z[0])
	}
	if math.Abs(real(z[1])-127.0/128.0) > 1e-9 || imag(z[1]) != 0 {
		t.Fatalf("z1=%v", z[1])
	}
}

func TestClipFraction(t *testing.T) {
	b := []byte{0, 1, 127, 128} // 127 and -128 (0x80) are rails
	if g := ClipFraction(b); math.Abs(g-0.5) > 1e-12 {
		t.Fatalf("clip=%g", g)
	}
}

func TestPowerDBFS(t *testing.T) {
	if g := PowerDBFS(1); math.Abs(g) > 1e-12 {
		t.Fatalf("full scale %g", g)
	}
	if g := PowerDBFS(0.01); math.Abs(g+20) > 1e-12 {
		t.Fatalf("0.01 -> %g", g)
	}
}

func TestMeanIQ(t *testing.T) {
	z := []complex128{1 + 2i, 3 + 4i}
	i, q := MeanIQ(z)
	if math.Abs(i-2) > 1e-12 || math.Abs(q-3) > 1e-12 {
		t.Fatalf("mean=%g,%g", i, q)
	}
}
