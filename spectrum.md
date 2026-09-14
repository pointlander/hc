# HackRF spectral difference

Captured 2026-09-14 15:10:30 UTC.

Two HackRF One radios received the same band with identical gain, sample rate, and baseband filter. Each capture is a Hann-windowed Welch PSD. A constant median offset (gain / noise-figure difference) is subtracted; bins whose residual exceeds **max(6·σ<sub>MAD</sub>, 3.0 dB)** are reported as significant.

## Radios

| | Label | Serial | USB | Firmware | Samples | Clip | DC (I, Q) | Power |
| --- | --- | --- | --- | --- | ---: | ---: | --- | ---: |
| A | 288e2dc3 | `0000000000000000f75461dc288e2dc3` | 5-2 | 2023.01.1 | 2000000 | 0.000% | -0.0037, +0.0124 | -30.0 dBFS |
| B | 2c8f32c3 | `0000000000000000f75461dc2c8f32c3` | 5-1 | 2023.01.1 | 2000000 | 0.000% | +0.0113, +0.0054 | -33.1 dBFS |

## Setup

- Sample rate: 8.000 MS/s
- Capture: 250ms per band (first 50 ms discarded)
- FFT: 4096, Hann, overlap 2048 samples (50%)
- RX gain: LNA 16 dB, VGA 16 dB, RF amp false
- Baseband filter: nearest supported ≤ 75% of sample rate
- Significance: residual ≥ max(6 σ, 3.0 dB) after median offset; DC and analog-filter edges excluded from σ

## Summary

- Across **4** band(s), radio **B** is **0.10 dB** hotter (median PSD offset).
- Mean noise-floor Δ(A−B): **-0.12 dB**.
- Mean LO leakage Δ(A−B): **+3.65 dB** (DC bin).
- **14** significant bins in total.

Largest residuals:

| Center | Offset | A | B | Δ | residual | σ |
| ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| +915 MHz | -3 MHz | -21.8 dBFS | -30.9 dBFS | +9.09 dB | +8.52 dB | +38.2 |
| +915 MHz | -1.50391 MHz | -37.0 dBFS | -29.5 dBFS | -7.51 dB | -8.08 dB | -36.3 |
| +915 MHz | -1.45312 MHz | -40.3 dBFS | -32.9 dBFS | -7.33 dB | -7.90 dB | -35.4 |
| +433 MHz | +1 MHz | -44.9 dBFS | -38.1 dBFS | -6.79 dB | -7.42 dB | -34.5 |
| +915 MHz | +2 MHz | -41.7 dBFS | -36.2 dBFS | -5.42 dB | -6.00 dB | -26.9 |
| +2.45 GHz | -1.953 kHz | -0.4 dBFS | -3.6 dBFS | +3.18 dB | +5.05 dB | +6.2 |
| +915 MHz | +3 MHz | -42.7 dBFS | -38.3 dBFS | -4.42 dB | -4.99 dB | -22.4 |
| +100 MHz | +701.2 kHz | -28.1 dBFS | -23.8 dBFS | -4.25 dB | -4.51 dB | -17.0 |

## Band +100 MHz

- Welch windows: A 975, B 975
- Median offset A−B: **+0.26 dB**
- Residual σ (MAD): **0.264 dB**; threshold **3.00 dB**
- Noise floor: A **-45.6 dBFS**, B **-45.9 dBFS** (Δ +0.27 dB)
- LO leakage (DC): A **3.6 dBFS**, B **-0.3 dBFS** (Δ +3.84 dB, residual +3.59 dB)
- RMS residual: **0.74 dB**
- Significant bins: **5**

| Offset | A | B | Δ | residual | σ |
| ---: | ---: | ---: | ---: | ---: | ---: |
| +701.2 kHz | -28.1 dBFS | -23.8 dBFS | -4.25 dB | -4.51 dB | -17.0 |
| +669.9 kHz | -38.0 dBFS | -34.1 dBFS | -3.98 dB | -4.23 dB | -16.0 |
| +750 kHz | -39.6 dBFS | -35.8 dBFS | -3.79 dB | -4.05 dB | -15.3 |
| +636.7 kHz | -39.7 dBFS | -36.1 dBFS | -3.65 dB | -3.90 dB | -14.8 |
| -1.953 kHz | -2.5 dBFS | -6.3 dBFS | +3.84 dB | +3.59 dB | +13.6 |

Residual spectrum (dB, median offset removed). `A` is hotter above the axis.

```
 +6.0                                                                         
                                                                              
                                                                              
                                         ##                                   
      ...................................||...................................
                                         ||                                   
                    *                    ||                             *     
         * ** **   *|* * **    ** *****  ||    *     * ***  * *****   **|  ** 
  0.0 ------------------------------------------------------------------------
      *** *  *  ***   * *  ****  *     **  **** |||** *   |* *     *|*   **  *
                                                *||       *         *         
                                                 ||                           
      ...........................................||...........................
                                                 ||                           
                                                 ##                           
                                                                              
 -6.0                                                                         
      -3.20117 MHz                       0                        +3.20117 MHz
```

## Band +433 MHz

- Welch windows: A 975, B 975
- Median offset A−B: **+0.64 dB**
- Residual σ (MAD): **0.215 dB**; threshold **3.00 dB**
- Noise floor: A **-45.3 dBFS**, B **-46.0 dBFS** (Δ +0.65 dB)
- LO leakage (DC): A **3.6 dBFS**, B **-0.2 dBFS** (Δ +3.78 dB, residual +3.14 dB)
- RMS residual: **0.30 dB**
- Significant bins: **2**

| Offset | A | B | Δ | residual | σ |
| ---: | ---: | ---: | ---: | ---: | ---: |
| +1 MHz | -44.9 dBFS | -38.1 dBFS | -6.79 dB | -7.42 dB | -34.5 |
| -1.953 kHz | -2.4 dBFS | -6.2 dBFS | +3.78 dB | +3.15 dB | +14.6 |

Residual spectrum (dB, median offset removed). `A` is hotter above the axis.

```
 +7.4                                                                         
                                                                              
                                                                              
                                                                              
                                                                              
      ...................................##...................................
                                         ||                                   
        *****      ****** * *     * ** **||     *  *       *  ** **      ** * 
  0.0 **-------**-*--------*-------*--*------*-*-**----*----**-----*---*------
             **  *       *   *|***         ** *     *|* ***     *   *** *  * *
                              *                      |                        
      ...............................................|........................
                                                     |                        
                                                     |                        
                                                     |                        
                                                     |                        
 -7.4                                                #                        
      -3.20117 MHz                       0                        +3.20117 MHz
```

## Band +915 MHz

- Welch windows: A 975, B 975
- Median offset A−B: **+0.57 dB**
- Residual σ (MAD): **0.223 dB**; threshold **3.00 dB**
- Noise floor: A **-45.4 dBFS**, B **-46.0 dBFS** (Δ +0.59 dB)
- LO leakage (DC): A **3.6 dBFS**, B **-0.2 dBFS** (Δ +3.80 dB, residual +3.23 dB)
- RMS residual: **0.55 dB**
- Significant bins: **6**

| Offset | A | B | Δ | residual | σ |
| ---: | ---: | ---: | ---: | ---: | ---: |
| -3 MHz | -21.8 dBFS | -30.9 dBFS | +9.09 dB | +8.52 dB | +38.2 |
| -1.50391 MHz | -37.0 dBFS | -29.5 dBFS | -7.51 dB | -8.08 dB | -36.3 |
| -1.45312 MHz | -40.3 dBFS | -32.9 dBFS | -7.33 dB | -7.90 dB | -35.4 |
| +2 MHz | -41.7 dBFS | -36.2 dBFS | -5.42 dB | -6.00 dB | -26.9 |
| +3 MHz | -42.7 dBFS | -38.3 dBFS | -4.42 dB | -4.99 dB | -22.4 |
| -1.953 kHz | -2.4 dBFS | -6.2 dBFS | +3.80 dB | +3.23 dB | +14.5 |

Residual spectrum (dB, median offset removed). `A` is hotter above the axis.

```
 +8.5   #                                                                     
        |                                                                     
        |                                                                     
        |                                                                     
        |                                                                     
      ..|................................##...................................
        |                                ||                                   
        |*  **             *        **  *|| *  ** *          **     **  *     
  0.0 **--**--**-**-*-**-----*-*****--**-----*---*-**-*----**--*--**--*--**-*-
                *  * *  ||* * |            *  *      * ****     |*     *   | *
                        ||    *                                 |          |  
      ..................#|......................................|..........|..
                         |                                      |          |  
                         |                                      |          #  
                         |                                      #             
                         |                                                    
 -8.5                    #                                                    
      -3.20117 MHz                       0                        +3.20117 MHz
```

## Band +2.45 GHz

- Welch windows: A 975, B 975
- Median offset A−B: **-1.88 dB**
- Residual σ (MAD): **0.820 dB**; threshold **4.92 dB**
- Noise floor: A **-41.3 dBFS**, B **-39.3 dBFS** (Δ -1.99 dB)
- LO leakage (DC): A **5.6 dBFS**, B **2.4 dBFS** (Δ +3.18 dB, residual +5.05 dB)
- RMS residual: **0.85 dB**
- Significant bins: **1**

| Offset | A | B | Δ | residual | σ |
| ---: | ---: | ---: | ---: | ---: | ---: |
| -1.953 kHz | -0.4 dBFS | -3.6 dBFS | +3.18 dB | +5.05 dB | +6.2 |

Residual spectrum (dB, median offset removed). `A` is hotter above the axis.

```
 +9.8                                                                         
                                                                              
                                                                              
                                                                              
      ...................................##...................................
                                         ||                                   
                                         ||                                   
             **                       ***||**************************         
  0.0 -------------------------------*----------------------------------------
      *******  **********|***********                                *********
                         *                                                    
                                                                              
      ........................................................................
                                                                              
                                                                              
                                                                              
 -9.8                                                                         
      -3.20117 MHz                       0                        +3.20117 MHz
```


---

# Casimir E-sandwich radio simulation

Simulated 2026-09-14 15:10:32 UTC.

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
| E outline (spine × tines) | 40.0 mm × 50.0 mm |
| Spine thickness / tine width | 8.0 mm / 8.0 mm |
| Slot between tines | 8.0 mm |
| E sheet thickness | 0.40 mm |
| Anodization (each plate) | **1 µm** |
| Al₂O₃ ε<sub>r</sub> / tanδ | 9.8 / 0.015 |
| Al conductivity | 3.56e+07 S/m |
| Temperature | 293.1 K |
| Probe Z | 50 Ω |
| E metal area | 13.28 cm² |

## Electromagnetics

Lying flat, the E is a **very low-impedance stripline** between the two plates (Z<sub>0</sub> milliohms). Skin-effect loss in the aluminum dominates, so the geometric half-wave modes (TM<sub>10</sub> along the tines ~ 958 MHz, slot path ~ 725 MHz) are **overdamped**. The structure behaves as a ~230 nF MIM capacitor: |Z| falls with frequency and the radio output is a smooth thermal continuum, not a comb of spurs.

| | |
| --- | ---: |
| MIM capacitance (both gaps) | **230.46 nF** |
| Phase velocity c/√ε<sub>r</sub> | 0.319 c |
| Ideal TM<sub>10</sub> (along tines) | 957.7 MHz |
| Ideal TM<sub>01</sub> (along spine) | 1197.1 MHz |
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
