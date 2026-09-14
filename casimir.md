# Casimir E-sandwich radio simulation

Simulated 2026-09-14 15:10:15 UTC.

An **E-shaped aluminum** sheet **lies flat** (rotated 90°, spine along the sandwich, three arms as tines) between two **anodized aluminum plates**. Each inner face carries **1 µm** of anodic Al₂O₃, so both faces of the E see a metal–insulator–metal gap. The Casimir pressure lives in those gaps. The radio output is the Johnson–Nyquist field of the same structure, shaped by lossy stripline modes of the E. At RF, ħω ≪ kT, so zero-point energy does not radiate; what a HackRF can in principle couple to is thermal, with a spectral shape set by Re(Z(f)).

## Geometry

```
side (gap exaggerated):
  ==============================  top anodized Al plate
  ------------------------------  1 µm Al2O3
  ##############################  E sheet, lying flat
  ------------------------------  1 µm Al2O3
  ==============================  bottom anodized Al plate

plan (through the top plate), E rotated 90°:
     ####      ####      ####
     ####      ####      ####
     ####      ####      ####
     ########################
          tines up, spine along the plates
```

| | |
| --- | ---: |
| E outline (spine × tines) | 100.0 mm × 150.0 mm |
| Spine thickness / tine width | 20.0 mm / 20.0 mm |
| Slot between tines | 20.0 mm |
| E sheet thickness | 0.40 mm |
| Anodization (each plate) | **1 µm** |
| Al₂O₃ ε<sub>r</sub> / tanδ | 9.8 / 0.015 |
| Al conductivity | 3.56e+07 S/m |
| Temperature | 293.1 K |
| Probe Z | 50 Ω |
| E metal area | 98.00 cm² |

## Electromagnetics

Lying flat, the E is a **very low-impedance stripline** between the two plates (Z<sub>0</sub> milliohms). Skin-effect loss in the aluminum dominates, so the geometric half-wave modes (TM<sub>10</sub> along the tines ~ 319 MHz, slot path ~ 252 MHz) are **overdamped**. The structure behaves as a ~1701 nF MIM capacitor: |Z| falls with frequency and the radio output is a smooth thermal continuum, not a comb of spurs.

| | |
| --- | ---: |
| MIM capacitance (both gaps) | **1700.71 nF** |
| Phase velocity c/√ε<sub>r</sub> | 0.319 c |
| Ideal TM<sub>10</sub> (along tines) | 319.2 MHz |
| Ideal TM<sub>01</sub> (along spine) | 478.8 MHz |
| Ideal slot-lengthened path | 252.0 MHz |
| Series |Z| dip | none below 6 GHz (RC-like) |

## Casimir force

Leading Lifshitz term for perfect reflectors filled with Al₂O₃, plus a plasma-wavelength correction for Al (λ<sub>p</sub>≈107 nm ≪ 1 µm):

$$P = \frac{\pi^2 \hbar c}{240\, n\, d^4}\,\eta_{\mathrm{Al}},\quad n=\sqrt{\varepsilon_r}$$

| | |
| --- | ---: |
| Pressure per gap | **0.36 mPa** |
| Force on E (both faces) | **6.98e-06 N** (6978.76 nN) |
| Energy both gaps | -2.33e-12 J |
| Plate bending fundamental | **0.1 kHz** |
| Thermal x<sub>rms</sub> | 6.96e-13 m (0.70 pm) |
| Patch-potential δV<sub>rms</sub> (50 mV) | **3.48e-08 V** |

The mechanical channel is audio/ultrasonic, not a HackRF band. It would only appear at UHF/microwave if an RF pump mixed with the motion (not assumed here).

## Radio output (Nyquist into 50 Ω)

Open-circuit voltage PSD is 4kT Re(Z). Available power accounts for mismatch to 50 Ω. HackRF noise is roughly −170 dBm/Hz; this source is many tens of dB below that unless a near-field probe is pressed to the E.

| f | Re(Z) | Im(Z) | S<sub>v</sub> (open) | P(50 Ω) |
| ---: | ---: | ---: | ---: | ---: |
| 100.00 MHz | 4.31 mΩ | -3.92 mΩ | 6.98e-23 V²/Hz | **-208.6 dBm/Hz** |
| 433.00 MHz | 2.96 mΩ | -2.74 mΩ | 4.80e-23 V²/Hz | **-210.2 dBm/Hz** |
| 915.00 MHz | 2.49 mΩ | -2.25 mΩ | 4.04e-23 V²/Hz | **-210.9 dBm/Hz** |
| 2.45 GHz | 2.00 mΩ | -1.71 mΩ | 3.25e-23 V²/Hz | **-211.9 dBm/Hz** |

Available power vs frequency (dBm/Hz):

```
  -208                           ******                                        
                               ***||||||****                                    
                             **|||||||||||||***                                 
                           **||||||||||||||||||***                              
                         **|||||||||||||||||||||||**                            
                       **|||||||||||||||||||||||||||***                         
                     **||||||||||||||||||||||||||||||||***                      
                   **|||||||||||||||||||||||||||||||||||||***                   
        *        **||||||||||||||||||||||||||||||||||||||||||***                
        |********|||||||||||||||||||||||||||||||||||||||||||||||****            
        ||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||***         
        |||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||***      
        ||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||****  
  -213 ||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||**
        1.00 MHz                                                        6.00 GHz
        dBm/Hz into 50 Ω, log-frequency
```

## What to look for on the HackRFs

1. **No sharp Casimir spurs.** The 1 µm MIM is overdamped; thermal output is a smooth continuum tens of dB below HackRF noise.
2. **Near-field only.** A probe on the E would see milliohm-source Johnson noise; free-space coupling between the two radios will not show this sandwich.
3. **Mechanical / patch-potential noise around the plate fundamental**, which is kHz, not a spectral-difference peak in the 100 MHz–2.45 GHz captures.
4. Differences already seen between the two HackRFs (LO leakage, +1 MHz at 433 MHz, ±2/3 MHz at 915 MHz) are **receiver fingerprints**, not this sandwich.
