# hc

Compare two HackRF One radios, or simulate the radio output of a Casimir E-sandwich.

Two boards receive the same band with identical gain and sample rate. A Hann-windowed Welch PSD is taken on each path; a constant median offset (gain / noise figure) is subtracted; bins whose residual exceeds **max(6 σ<sub>MAD</sub>, 3 dB)** are written to markdown. `hc sim` models an E-shaped aluminum sheet lying flat between two anodized plates and reports the Nyquist spectrum of that structure.

## Requirements

- Go 1.25
- libhackrf (`libhackrf0`; headers optional — a minimal `hackrf.h` is vendored)
- Two HackRF One radios, user in the `plugdev` group (capture only)

```sh
go test ./...
go build -o hc ./cmd/hc
```

## Capture

```sh
./hc -out spectrum.md
./hc -sim -out spectrum.md          # append the Casimir report
./hc -freq 915e6 -dur 500ms -a 2dc3 -b 32c3
```

Default bands are 100 MHz, 433 MHz, 915 MHz, and 2.45 GHz at 8 MS/s, 250 ms per band, LNA 16 dB, VGA 16 dB, RF amp off. Radios are chosen by serial suffix, USB path, or list index (`-a`, `-b`); otherwise the first two enumerated boards.

| Flag | Default | |
| --- | --- | --- |
| `-freq` | `100e6,433e6,915e6,2.45e9` | center frequencies, Hz |
| `-rate` | `8e6` | sample rate, Hz |
| `-dur` | `250ms` | capture per band |
| `-fft` | `4096` | Welch FFT size |
| `-sigma` | `6` | residual significance in robust σ |
| `-mindb` | `3` | minimum residual, dB |
| `-out` | `spectrum.md` | markdown path |
| `-sim` | false | append Casimir simulation |

An example capture is in [`spectrum.md`](spectrum.md).

## Casimir E-sandwich

The E is a coplanar aluminum sheet (spine and three tines in one plane) between two plates, each with **1 µm** of anodic Al₂O₃. Both faces of the E see that MIM gap. At RF, ħω ≪ kT, so zero-point energy does not radiate; the radio output is Johnson–Nyquist noise shaped by Re(Z(f)). The 1 µm gap makes a milliohm stripline: geometric modes are overdamped, and available power into 50 Ω is tens of dB below a HackRF noise floor unless a probe is on the E.

```
side (gap exaggerated):
  ==============================  top anodized Al plate
  ------------------------------  1 µm Al2O3
  ##############################  E sheet, lying flat
  ------------------------------  1 µm Al2O3
  ==============================  bottom anodized Al plate

plan (through the top plate):
     ####      ####      ####
     ####      ####      ####
     ####      ####      ####
     ########################
          tines up, spine along the plates
```

```sh
./hc sim -out casimir.md
./hc sim -width 0.1 -height 0.15 -spine 0.02 -arm 0.02 -oxide 1e-6
```

| Flag | Default | |
| --- | --- | --- |
| `-width` | `0.1` | along the spine, m |
| `-height` | `0.15` | along the tines, m |
| `-spine` | `0.02` | spine bar thickness, m |
| `-arm` | `0.02` | tine width, m |
| `-thick` | `4e-4` | sheet thickness, m |
| `-oxide` | `1e-6` | anodization per plate, m |
| `-eps` | `9.8` | Al₂O₃ ε<sub>r</sub> |
| `-out` | `casimir.md` | markdown path |

An example simulation is in [`casimir.md`](casimir.md).

## License

BSD-3-Clause. See [LICENSE](LICENSE).
