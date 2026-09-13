# HackRF spectral difference

Captured 2026-09-12 21:00:36 UTC.

Two HackRF One radios received the same band with identical gain, sample rate, and baseband filter. Each capture is a Hann-windowed Welch PSD. A constant median offset (gain / noise-figure difference) is subtracted; bins whose residual exceeds **max(6·σ<sub>MAD</sub>, 3.0 dB)** are reported as significant.

## Radios

| | Label | Serial | USB | Firmware | Samples | Clip | DC (I, Q) | Power |
| --- | --- | --- | --- | --- | ---: | ---: | --- | ---: |
| A | 288e2dc3 | `0000000000000000f75461dc288e2dc3` | 5-2 | 2023.01.1 | 2000000 | 0.000% | -0.0039, +0.0124 | -30.0 dBFS |
| B | 2c8f32c3 | `0000000000000000f75461dc2c8f32c3` | 5-1 | 2023.01.1 | 2000000 | 0.000% | +0.0112, +0.0054 | -33.0 dBFS |

## Setup

- Sample rate: 8.000 MS/s
- Capture: 250ms per band (first 50 ms discarded)
- FFT: 4096, Hann, overlap 2048 samples (50%)
- RX gain: LNA 16 dB, VGA 16 dB, RF amp false
- Baseband filter: nearest supported ≤ 75% of sample rate
- Significance: residual ≥ max(6 σ, 3.0 dB) after median offset; DC and analog-filter edges excluded from σ

## Summary

- Across **4** band(s), radio **B** is **0.36 dB** hotter (median PSD offset).
- Mean noise-floor Δ(A−B): **-0.39 dB**.
- Mean LO leakage Δ(A−B): **+3.72 dB** (DC bin).
- **13** significant bins in total.

Largest residuals:

| Center | Offset | A | B | Δ | residual | σ |
| ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| +915 MHz | +3 MHz | -44.0 dBFS | -36.9 dBFS | -7.11 dB | -7.79 dB | -36.7 |
| +433 MHz | +1 MHz | -45.0 dBFS | -38.7 dBFS | -6.26 dB | -6.89 dB | -31.2 |
| +915 MHz | -3 MHz | -22.3 dBFS | -17.8 dBFS | -4.46 dB | -5.14 dB | -24.2 |
| +433 MHz | +2 MHz | -44.3 dBFS | -39.8 dBFS | -4.46 dB | -5.09 dB | -23.1 |
| +100 MHz | +763.7 kHz | -38.1 dBFS | -33.8 dBFS | -4.22 dB | -4.48 dB | -17.9 |
| +100 MHz | +681.6 kHz | -36.0 dBFS | -32.1 dBFS | -3.92 dB | -4.18 dB | -16.7 |
| +100 MHz | +634.8 kHz | -36.3 dBFS | -32.4 dBFS | -3.89 dB | -4.15 dB | -16.6 |
| +100 MHz | +714.8 kHz | -35.7 dBFS | -31.9 dBFS | -3.79 dB | -4.05 dB | -16.2 |

## Band +100 MHz

- Welch windows: A 975, B 975
- Median offset A−B: **+0.26 dB**
- Residual σ (MAD): **0.251 dB**; threshold **3.00 dB**
- Noise floor: A **-45.6 dBFS**, B **-45.8 dBFS** (Δ +0.27 dB)
- LO leakage (DC): A **3.6 dBFS**, B **-0.3 dBFS** (Δ +3.86 dB, residual +3.60 dB)
- RMS residual: **0.73 dB**
- Significant bins: **5**

| Offset | A | B | Δ | residual | σ |
| ---: | ---: | ---: | ---: | ---: | ---: |
| +763.7 kHz | -38.1 dBFS | -33.8 dBFS | -4.22 dB | -4.48 dB | -17.9 |
| +681.6 kHz | -36.0 dBFS | -32.1 dBFS | -3.92 dB | -4.18 dB | -16.7 |
| +634.8 kHz | -36.3 dBFS | -32.4 dBFS | -3.89 dB | -4.15 dB | -16.6 |
| +714.8 kHz | -35.7 dBFS | -31.9 dBFS | -3.79 dB | -4.05 dB | -16.2 |
| 0 Hz | 3.6 dBFS | -0.3 dBFS | +3.86 dB | +3.60 dB | +14.4 |

Residual spectrum (dB, median offset removed). `A` is hotter above the axis.

```
 +6.0                                                                         
                                                                              
                                                                              
                                         ##                                   
      ...................................||...................................
                    *                    ||                                   
                   *|*                   ||                            **     
        ** *  **   ||| **         *****  || ***      ** **  * ****   **||   * 
  0.0 ------*--------------------*----------------------------------------**--
      **  *  *  ***   *  ********      **  |   *|||**  *  ** *    ***    *   *
                                           *    |||                           
                                                *||                           
      ...........................................||...........................
                                                 ||                           
                                                 ##                           
                                                                              
 -6.0                                                                         
      -3.20117 MHz                       0                        +3.20117 MHz
```

## Band +433 MHz

- Welch windows: A 975, B 975
- Median offset A−B: **+0.63 dB**
- Residual σ (MAD): **0.220 dB**; threshold **3.00 dB**
- Noise floor: A **-45.3 dBFS**, B **-45.9 dBFS** (Δ +0.63 dB)
- LO leakage (DC): A **3.6 dBFS**, B **-0.3 dBFS** (Δ +3.86 dB, residual +3.23 dB)
- RMS residual: **0.31 dB**
- Significant bins: **3**

| Offset | A | B | Δ | residual | σ |
| ---: | ---: | ---: | ---: | ---: | ---: |
| +1 MHz | -45.0 dBFS | -38.7 dBFS | -6.26 dB | -6.89 dB | -31.2 |
| +2 MHz | -44.3 dBFS | -39.8 dBFS | -4.46 dB | -5.09 dB | -23.1 |
| -1.953 kHz | -2.4 dBFS | -6.3 dBFS | +3.87 dB | +3.24 dB | +14.7 |

Residual spectrum (dB, median offset removed). `A` is hotter above the axis.

```
 +6.9                                                                         
                                                                              
                                                                              
                                                                              
                                         ##                                   
      .............*.....................||...................................
        *          |                     ||                                   
       *|*** * ****|* ****   *    * ***  ||     *****        **       *  **   
  0.0 ---------------*---------**--*---*---------------**--*------***---*-----
      *     * *           *** |  *      *  *****     |*  ** *  *|*   * *   ***
                              |                      |          |             
      ........................*......................|..........|.............
                                                     |          |             
                                                     |          |             
                                                     |          #             
                                                     |                        
 -6.9                                                #                        
      -3.20117 MHz                       0                        +3.20117 MHz
```

## Band +915 MHz

- Welch windows: A 975, B 975
- Median offset A−B: **+0.68 dB**
- Residual σ (MAD): **0.212 dB**; threshold **3.00 dB**
- Noise floor: A **-45.4 dBFS**, B **-46.1 dBFS** (Δ +0.69 dB)
- LO leakage (DC): A **3.6 dBFS**, B **-0.3 dBFS** (Δ +3.90 dB, residual +3.22 dB)
- RMS residual: **0.34 dB**
- Significant bins: **5**

| Offset | A | B | Δ | residual | σ |
| ---: | ---: | ---: | ---: | ---: | ---: |
| +3 MHz | -44.0 dBFS | -36.9 dBFS | -7.11 dB | -7.79 dB | -36.7 |
| -3 MHz | -22.3 dBFS | -17.8 dBFS | -4.46 dB | -5.14 dB | -24.2 |
| +4 MHz | -43.2 dBFS | -40.4 dBFS | -2.80 dB | -3.48 dB | -16.4 |
| +1.953 kHz | -2.4 dBFS | -6.3 dBFS | +3.90 dB | +3.22 dB | +15.2 |
| +2 MHz | -35.8 dBFS | -33.3 dBFS | -2.47 dB | -3.15 dB | -14.8 |

Residual spectrum (dB, median offset removed). `A` is hotter above the axis.

```
 +7.8                                                                         
                                                                              
                                                                              
                                                                              
                                                                              
      ...................................##...................................
                                         ||                                   
            *     *     *** ** **     ** ||*       **  * * *   *   * *        
  0.0 -*---*-*---*--****---*-----*-***--*---***-*-----*-*-*--**--**-*-*****-**
      * |**   ***  *          *   *            * **  *      *   |          |  
        |                                                       |          |  
      ..|.......................................................#..........|..
        |                                                                  |  
        #                                                                  |  
                                                                           |  
                                                                           |  
 -7.8                                                                      #  
      -3.20117 MHz                       0                        +3.20117 MHz
```

## Band +2.45 GHz

- Welch windows: A 975, B 975
- Median offset A−B: **-3.02 dB**
- Residual σ (MAD): **1.115 dB**; threshold **6.69 dB**
- Noise floor: A **-40.5 dBFS**, B **-37.4 dBFS** (Δ -3.14 dB)
- LO leakage (DC): A **5.6 dBFS**, B **2.4 dBFS** (Δ +3.24 dB, residual +6.27 dB)
- RMS residual: **0.99 dB**
- Significant bins: **0**

Residual spectrum (dB, median offset removed). `A` is hotter above the axis.

```
+13.4                                                                         
                                                                              
                                                                              
                                                                              
      ...................................**...................................
                                         ||                     *             
                                         ||                     |             
             *                        ***||*********************|******       
  0.0 ------------------------------**---------------------------------**-----
      ******* **********************                                     *****
                                                                              
                                                                              
      ........................................................................
                                                                              
                                                                              
                                                                              
-13.4                                                                         
      -3.20117 MHz                       0                        +3.20117 MHz
```

