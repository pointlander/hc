package spectrum

import (
	"math"
	"sort"
)

const (
	madToSigma = 1.4826
	floorCutDB = -120.0 // ignore numerical zeros when estimating offset/σ
)

// Peak is one frequency bin where the two radios differ significantly.
type Peak struct {
	Bin        int
	FreqHz     float64
	ADB        float64
	BDB        float64
	DeltaDB    float64
	ResidualDB float64
	Sigma      float64
}

// Diff is the significant spectral difference between two PSDs captured
// with identical settings.
type Diff struct {
	OffsetDB    float64 // median(A_dB - B_dB) over measured passband bins
	SigmaDB     float64 // robust σ of the residual (1.4826·MAD)
	ThresholdDB float64
	NoiseADB    float64
	NoiseBDB    float64
	DCADB       float64
	DCBDB       float64
	RMSResidual float64
	Peaks       []Peak
	ResidualDB  []float64
	DeltaDB     []float64
}

// CompareOptions controls significance testing.
type CompareOptions struct {
	NSigma float64 // residual significance in robust σ; default 6
	MinDB  float64 // practical floor in dB; default 3
}

// Compare finds bins where two radios' spectra differ after removing a
// constant gain/noise-figure offset. Significance is |residual| ≥
// max(NSigma·σ_MAD, MinDB). Offset and σ are estimated from passband bins
// both radios actually measured (above the numerical floor). DC, analog
// filter edges, and unique-signal bins are excluded from that estimate so
// a spur cannot hide other differences. Reported peaks are local maxima.
func Compare(a, b PSD, opt CompareOptions) Diff {
	if opt.NSigma <= 0 {
		opt.NSigma = 6
	}
	if opt.MinDB <= 0 {
		opt.MinDB = 3
	}
	n := min(len(a.Power), len(b.Power))
	delta := make([]float64, n)
	for k := 0; k < n; k++ {
		delta[k] = a.DB(k) - b.DB(k)
	}

	pass := passbandBins(n)
	measured := make([]int, 0, len(pass))
	for _, k := range pass {
		if a.DB(k) > floorCutDB && b.DB(k) > floorCutDB {
			measured = append(measured, k)
		}
	}
	statBins := measured
	if len(statBins) < 16 {
		statBins = pass
	}

	offSamples := make([]float64, len(statBins))
	for i, k := range statBins {
		offSamples[i] = delta[k]
	}
	offset := median(offSamples)

	resid := make([]float64, n)
	var sumSq float64
	for k := 0; k < n; k++ {
		resid[k] = delta[k] - offset
		sumSq += resid[k] * resid[k]
	}
	residKeep := make([]float64, len(statBins))
	for i, k := range statBins {
		residKeep[i] = resid[k]
	}
	sigma := madToSigma * mad(residKeep, median(residKeep))
	if sigma < 1e-6 {
		sigma = 1e-6
	}
	thr := max(opt.NSigma*sigma, opt.MinDB)

	fs := a.SampleRate
	if fs == 0 {
		fs = b.SampleRate
	}
	nfft := a.FFTSize
	if nfft == 0 {
		nfft = b.FFTSize
	}
	tmp := PSD{FFTSize: nfft, SampleRate: fs}

	peaks := make([]Peak, 0)
	for k := 0; k < n; k++ {
		if math.Abs(resid[k]) < thr {
			continue
		}
		if a.DB(k) <= floorCutDB && b.DB(k) <= floorCutDB {
			continue
		}
		if !localMaxAbs(resid, k) {
			continue
		}
		peaks = append(peaks, Peak{
			Bin:        k,
			FreqHz:     tmp.BinHz(k),
			ADB:        a.DB(k),
			BDB:        b.DB(k),
			DeltaDB:    delta[k],
			ResidualDB: resid[k],
			Sigma:      resid[k] / sigma,
		})
	}
	sort.Slice(peaks, func(i, j int) bool {
		ai, aj := math.Abs(peaks[i].ResidualDB), math.Abs(peaks[j].ResidualDB)
		if ai != aj {
			return ai > aj
		}
		return math.Abs(peaks[i].FreqHz) < math.Abs(peaks[j].FreqHz)
	})
	if nfft > 0 && fs > 0 {
		// Keep the strongest peak in each ~16-bin window so a wide FM
		// occupier does not emit dozens of Hann sidelobe maxima.
		peaks = thinPeaks(peaks, 16*fs/float64(nfft))
	}

	rms := 0.0
	if n > 0 {
		rms = math.Sqrt(sumSq / float64(n))
	}

	dcA, dcB := -180.0, -180.0
	if len(a.Power) > 0 {
		dcA = a.DB(0)
	}
	if len(b.Power) > 0 {
		dcB = b.DB(0)
	}

	return Diff{
		OffsetDB:    offset,
		SigmaDB:     sigma,
		ThresholdDB: thr,
		NoiseADB:    a.NoiseFloorDB(),
		NoiseBDB:    b.NoiseFloorDB(),
		DCADB:       dcA,
		DCBDB:       dcB,
		RMSResidual: rms,
		Peaks:       peaks,
		ResidualDB:  resid,
		DeltaDB:     delta,
	}
}

func thinPeaks(peaks []Peak, minSepHz float64) []Peak {
	if minSepHz <= 0 || len(peaks) < 2 {
		return peaks
	}
	keep := make([]Peak, 0, len(peaks))
	for _, p := range peaks {
		if p.Bin == 0 {
			filtered := make([]Peak, 0, len(keep)+1)
			for _, q := range keep {
				if math.Abs(q.FreqHz) >= minSepHz {
					filtered = append(filtered, q)
				}
			}
			keep = append(filtered, p)
			continue
		}
		close := false
		for _, q := range keep {
			if math.Abs(p.FreqHz-q.FreqHz) < minSepHz {
				close = true
				break
			}
		}
		if !close {
			keep = append(keep, p)
		}
	}
	sort.Slice(keep, func(i, j int) bool {
		ai, aj := math.Abs(keep[i].ResidualDB), math.Abs(keep[j].ResidualDB)
		if ai != aj {
			return ai > aj
		}
		return math.Abs(keep[i].FreqHz) < math.Abs(keep[j].FreqHz)
	})
	return keep
}

func localMaxAbs(v []float64, k int) bool {
	n := len(v)
	if n == 0 {
		return false
	}
	ak := math.Abs(v[k])
	prev := math.Abs(v[(k-1+n)%n])
	next := math.Abs(v[(k+1)%n])
	return ak >= prev && ak >= next
}

func passbandBins(n int) []int {
	if n < 16 {
		out := make([]int, 0, n)
		for k := 1; k < n; k++ {
			out = append(out, k)
		}
		return out
	}
	edge := n / 10
	if edge < 4 {
		edge = 4
	}
	dc := 3
	out := make([]int, 0, n)
	for k := 0; k < n; k++ {
		if k <= dc || k >= n-dc {
			continue
		}
		if k > n/2-edge && k < n/2+edge {
			continue
		}
		out = append(out, k)
	}
	return out
}
