package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pointlander/hc/internal/casimir"
	"github.com/pointlander/hc/internal/hackrf"
	"github.com/pointlander/hc/internal/spectrum"
)

const (
	defaultRate   = 8e6
	defaultDur    = 250 * time.Millisecond
	defaultSkip   = 50 * time.Millisecond
	defaultFFT    = 4096
	defaultLNA    = 16
	defaultVGA    = 16
	defaultNSigma = 6.0
	defaultMinDB  = 3.0
	defaultOut    = "spectrum.md"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "sim":
			os.Exit(runSim(os.Args[2:]))
		case "help", "-h", "--help":
			usage()
			os.Exit(0)
		}
	}
	os.Exit(run(os.Args[1:]))
}

func usage() {
	fmt.Fprintf(os.Stderr, `hc — compare two HackRF One radios, or simulate a Casimir E-sandwich

Usage:
  hc [flags]       capture both radios and write spectral difference
  hc sim [flags]   simulate radio output of the E-shaped Al sandwich
  hc help

`)
}

func run(args []string) int {
	fs := flag.NewFlagSet("hc", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	freqStr := fs.String("freq", "100e6,433e6,915e6,2.45e9", "center frequency in Hz, comma-separated")
	rate := fs.Float64("rate", defaultRate, "sample rate in Hz")
	dur := fs.Duration("dur", defaultDur, "capture duration per band")
	fftSize := fs.Int("fft", defaultFFT, "FFT size (power of two)")
	overlapFrac := fs.Float64("overlap", 0.5, "Welch overlap fraction")
	lna := fs.Int("lna", defaultLNA, "LNA gain dB (0-40, 8 dB steps)")
	vga := fs.Int("vga", defaultVGA, "VGA gain dB (0-62, 2 dB steps)")
	amp := fs.Bool("amp", false, "enable RF amplifier")
	antenna := fs.Bool("antenna", false, "enable antenna port power")
	nSigma := fs.Float64("sigma", defaultNSigma, "significance in robust σ (MAD)")
	minDB := fs.Float64("mindb", defaultMinDB, "minimum practical residual in dB")
	outPath := fs.String("out", defaultOut, "markdown output path")
	aWant := fs.String("a", "", "radio A serial suffix, USB path, or index")
	bWant := fs.String("b", "", "radio B serial suffix, USB path, or index")
	withSim := fs.Bool("sim", false, "append Casimir E-sandwich simulation to the report")
	if err := fs.Parse(args); err != nil {
		return 2
	}

	freqs, err := parseFreqs(*freqStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "-freq: %v\n", err)
		return 2
	}
	if *rate < 2e6 || *rate > 20e6 {
		fmt.Fprintf(os.Stderr, "-rate: %.0f Hz out of HackRF range (2e6–20e6)\n", *rate)
		return 2
	}
	if *fftSize < 64 || *fftSize&(*fftSize-1) != 0 {
		fmt.Fprintf(os.Stderr, "-fft: %d must be a power of two >= 64\n", *fftSize)
		return 2
	}
	if *overlapFrac < 0 || *overlapFrac >= 1 {
		fmt.Fprintf(os.Stderr, "-overlap: %g must be in [0, 1)\n", *overlapFrac)
		return 2
	}
	if *dur < 20*time.Millisecond {
		fmt.Fprintf(os.Stderr, "-dur: %s too short\n", *dur)
		return 2
	}

	if err := hackrf.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "hackrf init: %v\n", err)
		return 1
	}
	defer hackrf.Exit()

	devs, err := hackrf.List()
	if err != nil {
		fmt.Fprintf(os.Stderr, "hackrf list: %v\n", err)
		return 1
	}
	infoA, infoB, err := hackrf.PickTwo(devs, *aWant, *bWant)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 1
	}

	devA, err := hackrf.OpenDevice(infoA)
	if err != nil {
		fmt.Fprintf(os.Stderr, "open A %s: %v\n", infoA.Label(), err)
		return 1
	}
	defer devA.Close()
	devB, err := hackrf.OpenDevice(infoB)
	if err != nil {
		fmt.Fprintf(os.Stderr, "open B %s: %v\n", infoB.Label(), err)
		return 1
	}
	defer devB.Close()

	if err := configure(devA, *rate, *lna, *vga, *amp, *antenna); err != nil {
		fmt.Fprintf(os.Stderr, "configure A: %v\n", err)
		return 1
	}
	if err := configure(devB, *rate, *lna, *vga, *amp, *antenna); err != nil {
		fmt.Fprintf(os.Stderr, "configure B: %v\n", err)
		return 1
	}

	fwA, _ := devA.Version()
	fwB, _ := devB.Version()
	fmt.Fprintf(os.Stderr, "libhackrf %s (%s)\n", hackrf.LibraryVersion(), hackrf.LibraryRelease())
	fmt.Fprintf(os.Stderr, "A  %s  serial=%s  usb=%s  fw=%s\n", infoA.USBBoardID, infoA.Serial, infoA.USBPath, fwA)
	fmt.Fprintf(os.Stderr, "B  %s  serial=%s  usb=%s  fw=%s\n", infoB.USBBoardID, infoB.Serial, infoB.USBPath, fwB)

	overlap := int(*overlapFrac * float64(*fftSize))
	nbytes := int(*rate*(*dur+defaultSkip).Seconds()) * 2
	if nbytes < (*fftSize)*4 {
		nbytes = (*fftSize) * 4
	}

	rep := spectrum.Report{
		Time:     time.Now().UTC(),
		RateHz:   *rate,
		Duration: *dur,
		FFTSize:  *fftSize,
		Overlap:  overlap,
		LNAGain:  *lna,
		VGAGain:  *vga,
		Amp:      *amp,
		NSigma:   *nSigma,
		MinDB:    *minDB,
		A: spectrum.Radio{
			Name:     infoA.Label(),
			Serial:   infoA.Serial,
			USBPath:  infoA.USBPath,
			Firmware: fwA,
		},
		B: spectrum.Radio{
			Name:     infoB.Label(),
			Serial:   infoB.Serial,
			USBPath:  infoB.USBPath,
			Firmware: fwB,
		},
	}

	opt := spectrum.CompareOptions{NSigma: *nSigma, MinDB: *minDB}
	var sumA, sumB captureStats
	var nStats int

	for _, freq := range freqs {
		fmt.Fprintf(os.Stderr, "capturing %.6f MHz …\n", float64(freq)/1e6)
		if err := devA.SetFreq(freq); err != nil {
			fmt.Fprintf(os.Stderr, "set freq A: %v\n", err)
			return 1
		}
		if err := devB.SetFreq(freq); err != nil {
			fmt.Fprintf(os.Stderr, "set freq B: %v\n", err)
			return 1
		}
		time.Sleep(80 * time.Millisecond)

		rawA, rawB, err := capturePair(devA, devB, nbytes, *dur+defaultSkip+2*time.Second)
		if err != nil {
			fmt.Fprintf(os.Stderr, "capture: %v\n", err)
			return 1
		}
		skipBytes := int(*rate*defaultSkip.Seconds()) * 2
		rawA = dropPrefix(rawA, skipBytes)
		rawB = dropPrefix(rawB, skipBytes)
		if len(rawA) < *fftSize*2 || len(rawB) < *fftSize*2 {
			fmt.Fprintf(os.Stderr, "capture too short: A=%d B=%d bytes\n", len(rawA), len(rawB))
			return 1
		}

		iqA := spectrum.BytesToIQ(rawA)
		iqB := spectrum.BytesToIQ(rawB)
		psdA := spectrum.Welch(iqA, *rate, *fftSize, overlap)
		psdB := spectrum.Welch(iqB, *rate, *fftSize, overlap)
		diff := spectrum.Compare(psdA, psdB, opt)
		rep.Bands = append(rep.Bands, spectrum.Band{
			CenterHz: freq,
			RateHz:   *rate,
			A:        psdA,
			B:        psdB,
			Diff:     diff,
		})

		stA := statsOf(rawA, iqA)
		stB := statsOf(rawB, iqB)
		sumA.samples = stA.samples
		sumB.samples = stB.samples
		sumA.clip += stA.clip
		sumB.clip += stB.clip
		sumA.meanI += stA.meanI
		sumA.meanQ += stA.meanQ
		sumB.meanI += stB.meanI
		sumB.meanQ += stB.meanQ
		sumA.power += stA.power
		sumB.power += stB.power
		nStats++
		fmt.Fprintf(os.Stderr, "  offset %+.2f dB  σ=%.3f dB  significant=%d\n",
			diff.OffsetDB, diff.SigmaDB, len(diff.Peaks))
	}

	if nStats > 0 {
		n := float64(nStats)
		rep.A.Samples = sumA.samples
		rep.A.ClipFrac = sumA.clip / n
		rep.A.MeanI = sumA.meanI / n
		rep.A.MeanQ = sumA.meanQ / n
		rep.A.PowerDBFS = sumA.power / n
		rep.B.Samples = sumB.samples
		rep.B.ClipFrac = sumB.clip / n
		rep.B.MeanI = sumB.meanI / n
		rep.B.MeanQ = sumB.meanQ / n
		rep.B.PowerDBFS = sumB.power / n
	}

	md := rep.Markdown()
	if *withSim {
		sim := casimir.Report{Time: time.Now().UTC(), Device: casimir.DefaultDevice(), Bands: freqFloats(freqs)}
		md += "\n---\n\n" + sim.Markdown()
	}
	if err := os.WriteFile(*outPath, []byte(md), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "write %s: %v\n", *outPath, err)
		return 1
	}
	fmt.Fprintf(os.Stderr, "wrote %s\n", *outPath)
	return 0
}

func runSim(args []string) int {
	fs := flag.NewFlagSet("sim", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	outPath := fs.String("out", "casimir.md", "markdown output path")
	oxide := fs.Float64("oxide", 1e-6, "anodization thickness, m")
	width := fs.Float64("width", 50e-3, "E width (arm direction), m")
	height := fs.Float64("height", 40e-3, "E height (across arms), m")
	spine := fs.Float64("spine", 8e-3, "spine width, m")
	arm := fs.Float64("arm", 8e-3, "arm width, m")
	thick := fs.Float64("thick", 0.4e-3, "E sheet thickness, m")
	eps := fs.Float64("eps", 9.8, "Al2O3 relative permittivity")
	tand := fs.Float64("tand", 0.015, "Al2O3 loss tangent")
	temp := fs.Float64("temp", 293.15, "temperature, K")
	freqStr := fs.String("freq", "100e6,433e6,915e6,2.45e9", "report frequencies, Hz")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	freqs, err := parseFreqs(*freqStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "-freq: %v\n", err)
		return 2
	}
	d := casimir.DefaultDevice()
	d.Oxide = *oxide
	d.Width = *width
	d.Height = *height
	d.Spine = *spine
	d.Arm = *arm
	d.Thick = *thick
	d.EpsR = *eps
	d.TanD = *tand
	d.Temp = *temp
	if err := d.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 2
	}
	rep := casimir.Report{Time: time.Now().UTC(), Device: d, Bands: freqFloats(freqs)}
	cas := d.Casimir()
	fsr, zsr, haveSR := d.SeriesResonance(1e6, 6e9)
	fmt.Fprintf(os.Stderr, "E %.1f×%.1f mm  oxide=%.3g µm  C=%.2f nF\n", d.Width*1e3, d.Height*1e3, d.Oxide*1e6, d.Capacitance()*1e9)
	fmt.Fprintf(os.Stderr, "Casimir P=%.2f mPa  F=%.2e N  mech=%.1f kHz\n", cas.Pressure*1e3, cas.Force, cas.MechHz/1e3)
	if haveSR {
		fmt.Fprintf(os.Stderr, "series |Z| dip %.1f MHz  |Z|=%.2f mΩ  Q=%.2f\n", fsr/1e6, zsr*1e3, d.Quality(fsr))
	} else {
		fmt.Fprintf(os.Stderr, "no series |Z| dip below 6 GHz (RC-like MIM)\n")
	}
	if err := os.WriteFile(*outPath, []byte(rep.Markdown()), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "write %s: %v\n", *outPath, err)
		return 1
	}
	fmt.Fprintf(os.Stderr, "wrote %s\n", *outPath)
	return 0
}

func freqFloats(f []uint64) []float64 {
	out := make([]float64, len(f))
	for i, v := range f {
		out[i] = float64(v)
	}
	return out
}

type captureStats struct {
	samples int
	clip    float64
	meanI   float64
	meanQ   float64
	power   float64
}

func statsOf(raw []byte, iq []complex128) captureStats {
	i, q := spectrum.MeanIQ(iq)
	return captureStats{
		samples: len(iq),
		clip:    spectrum.ClipFraction(raw),
		meanI:   i,
		meanQ:   q,
		power:   spectrum.PowerDBFS(spectrum.MeanPower(iq)),
	}
}

func configure(d *hackrf.Device, rate float64, lna, vga int, amp, antenna bool) error {
	if err := d.SetSampleRate(rate); err != nil {
		return fmt.Errorf("sample rate: %w", err)
	}
	bw := hackrf.ComputeBasebandFilterBW(uint32(rate * 0.75))
	if err := d.SetBasebandFilterBandwidth(bw); err != nil {
		return fmt.Errorf("filter: %w", err)
	}
	if err := d.SetAmpEnable(amp); err != nil {
		return fmt.Errorf("amp: %w", err)
	}
	if err := d.SetAntennaEnable(antenna); err != nil {
		return fmt.Errorf("antenna: %w", err)
	}
	if err := d.SetLNAGain(lna); err != nil {
		return fmt.Errorf("lna: %w", err)
	}
	if err := d.SetVGAGain(vga); err != nil {
		return fmt.Errorf("vga: %w", err)
	}
	return nil
}

func capturePair(a, b *hackrf.Device, nbytes int, timeout time.Duration) ([]byte, []byte, error) {
	recA := newRecorder(nbytes)
	recB := newRecorder(nbytes)
	if err := a.StartRX(recA.push); err != nil {
		return nil, nil, fmt.Errorf("start RX A: %w", err)
	}
	if err := b.StartRX(recB.push); err != nil {
		_ = a.StopRX()
		return nil, nil, fmt.Errorf("start RX B: %w", err)
	}

	timer := time.NewTimer(timeout)
	defer timer.Stop()
	var errA, errB error
	select {
	case <-both(recA.done, recB.done):
	case <-timer.C:
		errA = fmt.Errorf("timeout after %s (A %d/%d B %d/%d bytes)", timeout, recA.n(), nbytes, recB.n(), nbytes)
	}
	if e := a.StopRX(); e != nil && errA == nil {
		errA = e
	}
	if e := b.StopRX(); e != nil && errB == nil {
		errB = e
	}
	if errA != nil {
		return recA.bytes(), recB.bytes(), errA
	}
	if errB != nil {
		return recA.bytes(), recB.bytes(), errB
	}
	return recA.bytes(), recB.bytes(), nil
}

func both(a, b <-chan struct{}) <-chan struct{} {
	out := make(chan struct{})
	go func() {
		<-a
		<-b
		close(out)
	}()
	return out
}

type recorder struct {
	mu   sync.Mutex
	buf  []byte
	off  int
	once sync.Once
	done chan struct{}
}

func newRecorder(n int) *recorder {
	return &recorder{buf: make([]byte, n), done: make(chan struct{})}
}

func (r *recorder) push(b []byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.off >= len(r.buf) {
		return nil
	}
	r.off += copy(r.buf[r.off:], b)
	if r.off >= len(r.buf) {
		r.once.Do(func() { close(r.done) })
	}
	return nil
}

func (r *recorder) n() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.off
}

func (r *recorder) bytes() []byte {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]byte, r.off)
	copy(out, r.buf[:r.off])
	return out
}

func dropPrefix(b []byte, n int) []byte {
	if n <= 0 || n >= len(b) {
		return b
	}
	return b[n:]
}

func parseFreqs(s string) ([]uint64, error) {
	parts := strings.Split(s, ",")
	out := make([]uint64, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		v, err := strconv.ParseFloat(p, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid frequency %q", p)
		}
		if v < 1e6 || v > 6e9 {
			return nil, fmt.Errorf("frequency %g Hz out of HackRF range (1e6–6e9)", v)
		}
		out = append(out, uint64(v+0.5))
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no frequencies")
	}
	return out, nil
}
