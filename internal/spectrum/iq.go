package spectrum

import "math"

const fullScale = 128.0

// BytesToIQ converts interleaved HackRF int8 I/Q bytes into complex128
// samples scaled to roughly ±1.
func BytesToIQ(b []byte) []complex128 {
	n := len(b) / 2
	out := make([]complex128, n)
	for i := 0; i < n; i++ {
		i8 := int8(b[2*i])
		q8 := int8(b[2*i+1])
		out[i] = complex(float64(i8)/fullScale, float64(q8)/fullScale)
	}
	return out
}

// ClipFraction is the fraction of int8 samples at the rail (±127 or -128).
func ClipFraction(b []byte) float64 {
	if len(b) == 0 {
		return 0
	}
	n := 0
	for _, v := range b {
		s := int8(v)
		if s == 127 || s == -128 {
			n++
		}
	}
	return float64(n) / float64(len(b))
}

// MeanIQ returns the mean of the I and Q channels.
func MeanIQ(z []complex128) (i, q float64) {
	if len(z) == 0 {
		return 0, 0
	}
	for _, v := range z {
		i += real(v)
		q += imag(v)
	}
	n := float64(len(z))
	return i / n, q / n
}

// MeanPower returns mean |z|^2 of the samples.
func MeanPower(z []complex128) float64 {
	if len(z) == 0 {
		return 0
	}
	var s float64
	for _, v := range z {
		re := real(v)
		im := imag(v)
		s += re*re + im*im
	}
	return s / float64(len(z))
}

// PowerDBFS converts mean power (relative to full scale) to dBFS.
func PowerDBFS(p float64) float64 {
	if p <= 1e-18 {
		return -180
	}
	return 10 * math.Log10(p)
}

func db(p float64) float64 {
	if p <= 1e-18 {
		return -180
	}
	return 10 * math.Log10(p)
}
