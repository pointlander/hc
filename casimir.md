# Casimir E-sandwich radio simulation

Simulated 2026-09-13 16:07:45 UTC.

An **E-shaped aluminum** sheet is clamped between two **anodized aluminum plates**. Each inner face carries **1 µm** of anodic Al₂O₃, so the E sees two metal–insulator–metal gaps. The Casimir pressure lives in those gaps. The radio output is the Johnson–Nyquist field of the same structure, shaped by lossy stripline modes of the E. At RF, ħω ≪ kT, so zero-point energy does not radiate; what a HackRF can in principle couple to is thermal, with a spectral shape set by Re(Z(f)).

## Geometry

```
   anodized Al plate
  +------------------------------+
  | 1 µm Al2O3                   |
  |   ########################   |
  |   ####                       |
  |   ########################   |  E-shaped Al
  |   ####                       |
  |   ########################   |
  | 1 µm Al2O3                   |
  +------------------------------+
   anodized Al plate
```

| | |
| --- | ---: |
| E outline | 50.0 mm × 40.0 mm |
| Spine / arm | 8.0 mm / 8.0 mm |
| Slot height | 8.0 mm |
| E thickness | 0.40 mm |
| Anodization (each plate) | **1 µm** |
| Al₂O₃ ε<sub>r</sub> / tanδ | 9.8 / 0.015 |
| Al conductivity | 3.56e+07 S/m |
| Temperature | 293.1 K |
| Probe Z | 50 Ω |
| E metal area | 13.28 cm² |

## Electromagnetics

The 1 µm gap makes this a **very low-impedance stripline** (Z<sub>0</sub> milliohms). Skin-effect loss in the aluminum dominates, so the geometric half-wave modes (TM<sub>10</sub> ~ 958 MHz, slot path ~ 725 MHz) are **overdamped**. The structure behaves as a ~230 nF MIM capacitor: |Z| falls with frequency and the radio output is a smooth thermal continuum, not a comb of spurs.

| | |
| --- | ---: |
| MIM capacitance (both gaps) | **230.46 nF** |
| Phase velocity c/√ε<sub>r</sub> | 0.319 c |
| Ideal TM<sub>10</sub> (along arms) | 957.7 MHz |
| Ideal TM<sub>01</sub> (across arms) | 1197.1 MHz |
| Ideal slot-lengthened path | 725.5 MHz |
| Series |Z| dip | none below 6 GHz (RC-like) |

## Casimir force

Leading Lifshitz term for perfect reflectors filled with Al₂O₃, plus a plasma-wavelength correction for Al (λ<sub>p</sub>≈107 nm ≪ 1 µm):

$$P = \frac{\pi^2 \hbar c}{240\, n\, d^4}\,\eta_{\mathrm{Al}},\quad n=\sqrt{\varepsilon_r}$$

| | |
| --- | ---: |
| Pressure per gap | **0.36 mPa** |
| Force on E (both faces) | **9.46e-07 N** (945.69 nN) |
| Energy both gaps | -3.15e-13 J |
| Plate bending fundamental | **1.0 kHz** |
| Thermal x<sub>rms</sub> | 2.67e-13 m (0.27 pm) |
| Patch-potential δV<sub>rms</sub> (50 mV) | **1.33e-08 V** |

The mechanical channel is audio/ultrasonic, not a HackRF band. It would only appear at UHF/microwave if an RF pump mixed with the motion (not assumed here).

## Radio output (Nyquist into 50 Ω)

Open-circuit voltage PSD is 4kT Re(Z). Available power accounts for mismatch to 50 Ω. HackRF noise is roughly −170 dBm/Hz; this source is many tens of dB below that unless a near-field probe is pressed to the E.

| f | Re(Z) | Im(Z) | S<sub>v</sub> (open) | P(50 Ω) |
| ---: | ---: | ---: | ---: | ---: |
| 100.00 MHz | 8.81 mΩ | -9.16 mΩ | 1.43e-22 V²/Hz | **-205.5 dBm/Hz** |
| 433.00 MHz | 7.61 mΩ | -6.80 mΩ | 1.23e-22 V²/Hz | **-206.1 dBm/Hz** |
| 915.00 MHz | 6.22 mΩ | -5.65 mΩ | 1.01e-22 V²/Hz | **-207.0 dBm/Hz** |
| 2.45 GHz | 5.01 mΩ | -4.27 mΩ | 8.11e-23 V²/Hz | **-207.9 dBm/Hz** |

Available power vs frequency (dBm/Hz):

```
  -204 *                                                                       
        |                                                                       
        |*                                                                      
        ||*                                  ********                           
        |||*                               **||||||||***                        
        ||||                             **|||||||||||||***                     
        ||||*                          **||||||||||||||||||**                   
        |||||**                       *||||||||||||||||||||||**                 
        |||||||                     **|||||||||||||||||||||||||***              
        |||||||**                  *||||||||||||||||||||||||||||||***           
        |||||||||*              ***||||||||||||||||||||||||||||||||||***        
        ||||||||||*            *||||||||||||||||||||||||||||||||||||||||***     
        |||||||||||**       ***||||||||||||||||||||||||||||||||||||||||||||***  
  -209 |||||||||||||*******||||||||||||||||||||||||||||||||||||||||||||||||||**
        1.00 MHz                                                        6.00 GHz
        dBm/Hz into 50 Ω, log-frequency
```

## What to look for on the HackRFs

1. **No sharp Casimir spurs.** The 1 µm MIM is overdamped; thermal output is a smooth continuum tens of dB below HackRF noise.
2. **Near-field only.** A probe on the E would see milliohm-source Johnson noise; free-space coupling between the two radios will not show this sandwich.
3. **Mechanical / patch-potential noise around the plate fundamental**, which is kHz, not a spectral-difference peak in the 100 MHz–2.45 GHz captures.
4. Differences already seen between the two HackRFs (LO leakage, +1 MHz at 433 MHz, ±2/3 MHz at 915 MHz) are **receiver fingerprints**, not this sandwich.
