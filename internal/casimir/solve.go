package casimir

import "math/cmplx"

func solveComplex(a [][]complex128, b []complex128) ([]complex128, bool) {
	n := len(b)
	if n == 0 {
		return nil, false
	}
	m := make([][]complex128, n)
	rhs := make([]complex128, n)
	for i := 0; i < n; i++ {
		m[i] = append([]complex128(nil), a[i]...)
		rhs[i] = b[i]
	}
	for i := 0; i < n; i++ {
		pivot := i
		best := cmplx.Abs(m[i][i])
		for k := i + 1; k < n; k++ {
			if c := cmplx.Abs(m[k][i]); c > best {
				best, pivot = c, k
			}
		}
		if best < 1e-30 {
			return nil, false
		}
		m[i], m[pivot] = m[pivot], m[i]
		rhs[i], rhs[pivot] = rhs[pivot], rhs[i]
		piv := m[i][i]
		for k := i + 1; k < n; k++ {
			f := m[k][i] / piv
			for j := i; j < n; j++ {
				m[k][j] -= f * m[i][j]
			}
			rhs[k] -= f * rhs[i]
		}
	}
	x := make([]complex128, n)
	for i := n - 1; i >= 0; i-- {
		s := rhs[i]
		for j := i + 1; j < n; j++ {
			s -= m[i][j] * x[j]
		}
		x[i] = s / m[i][i]
	}
	return x, true
}
