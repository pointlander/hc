# HackRF spectral difference

Captured 2026-09-12 19:56:37 UTC.

Two HackRF One radios received the same band with identical gain, sample rate, and baseband filter. Each capture is a Hann-windowed Welch PSD. A constant median offset (gain / noise-figure difference) is subtracted; bins whose residual exceeds **max(6·σ<sub>MAD</sub>, 3.0 dB)** are reported as significant.

## Radios

| | Label | Serial | USB | Firmware | Samples | Clip | DC (I, Q) | Power |
| --- | --- | --- | --- | --- | ---: | ---: | --- | ---: |
| A | 288e2dc3 | `0000000000000000f75461dc288e2dc3` | 5-2 | 2023.01.1 | 2000000 | 0.000% | -0.0038, +0.0123 | -30.0 dBFS |
| B | 2c8f32c3 | `0000000000000000f75461dc2c8f32c3` | 5-1 | 2023.01.1 | 2000000 | 0.000% | +0.0113, +0.0055 | -33.4 dBFS |

## Setup

- Sample rate: 8.000 MS/s
- Capture: 250ms per band (first 50 ms discarded)
- FFT: 4096, Hann, overlap 2048 samples (50%)
- RX gain: LNA 16 dB, VGA 16 dB, RF amp false
- Baseband filter: nearest supported ≤ 75% of sample rate
- Significance: residual ≥ max(6 σ, 3.0 dB) after median offset; DC and analog-filter edges excluded from σ

## Summary

- Across **4** band(s), radio **A** is **0.67 dB** hotter (median PSD offset).
- Mean noise-floor Δ(A−B): **+0.66 dB**.
- Mean LO leakage Δ(A−B): **+3.65 dB** (DC bin).
- **24** significant bins in total.

Largest residuals:

| Center | Offset | A | B | Δ | residual | σ |
| ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| +100 MHz | +712.9 kHz | -35.0 dBFS | -26.9 dBFS | -8.18 dB | -8.32 dB | -27.0 |
| +100 MHz | +673.8 kHz | -36.6 dBFS | -28.7 dBFS | -7.96 dB | -8.10 dB | -26.3 |
| +100 MHz | +748 kHz | -37.9 dBFS | -30.2 dBFS | -7.70 dB | -7.84 dB | -25.5 |
| +100 MHz | +642.6 kHz | -39.0 dBFS | -31.4 dBFS | -7.59 dB | -7.72 dB | -25.1 |
| +100 MHz | +3.5293 MHz | -42.7 dBFS | -35.3 dBFS | -7.39 dB | -7.52 dB | -24.5 |
| +100 MHz | +3.4707 MHz | -42.5 dBFS | -35.1 dBFS | -7.36 dB | -7.50 dB | -24.4 |
| +433 MHz | +1 MHz | -45.2 dBFS | -39.2 dBFS | -6.09 dB | -6.80 dB | -31.5 |
| +100 MHz | +3.43945 MHz | -44.6 dBFS | -38.1 dBFS | -6.56 dB | -6.70 dB | -21.8 |

## Band +100 MHz

- Welch windows: A 975, B 975
- Median offset A−B: **+0.14 dB**
- Residual σ (MAD): **0.308 dB**; threshold **3.00 dB**
- Noise floor: A **-45.5 dBFS**, B **-45.7 dBFS** (Δ +0.20 dB)
- LO leakage (DC): A **3.5 dBFS**, B **-0.2 dBFS** (Δ +3.74 dB, residual +3.60 dB)
- RMS residual: **1.70 dB**
- Significant bins: **19**

| Offset | A | B | Δ | residual | σ |
| ---: | ---: | ---: | ---: | ---: | ---: |
| +712.9 kHz | -35.0 dBFS | -26.9 dBFS | -8.18 dB | -8.32 dB | -27.0 |
| +673.8 kHz | -36.6 dBFS | -28.7 dBFS | -7.96 dB | -8.10 dB | -26.3 |
| +748 kHz | -37.9 dBFS | -30.2 dBFS | -7.70 dB | -7.84 dB | -25.5 |
| +642.6 kHz | -39.0 dBFS | -31.4 dBFS | -7.59 dB | -7.72 dB | -25.1 |
| +3.5293 MHz | -42.7 dBFS | -35.3 dBFS | -7.39 dB | -7.52 dB | -24.5 |
| +3.4707 MHz | -42.5 dBFS | -35.1 dBFS | -7.36 dB | -7.50 dB | -24.4 |
| +3.43945 MHz | -44.6 dBFS | -38.1 dBFS | -6.56 dB | -6.70 dB | -21.8 |
| +3.56445 MHz | -45.3 dBFS | -39.4 dBFS | -5.82 dB | -5.96 dB | -19.4 |
| +2.70898 MHz | -40.7 dBFS | -36.3 dBFS | -4.41 dB | -4.54 dB | -14.8 |
| -3.30078 MHz | -34.9 dBFS | -30.7 dBFS | -4.14 dB | -4.28 dB | -13.9 |
| -3.26367 MHz | -35.5 dBFS | -31.4 dBFS | -4.11 dB | -4.24 dB | -13.8 |
| -3.33789 MHz | -35.8 dBFS | -31.8 dBFS | -4.00 dB | -4.13 dB | -13.4 |
| -3.23242 MHz | -35.7 dBFS | -31.8 dBFS | -3.88 dB | -4.01 dB | -13.0 |
| -396.5 kHz | -40.0 dBFS | -44.0 dBFS | +4.03 dB | +3.90 dB | +12.7 |
| 0 Hz | 3.5 dBFS | -0.2 dBFS | +3.74 dB | +3.60 dB | +11.7 |
| +2.5918 MHz | -44.9 dBFS | -41.6 dBFS | -3.30 dB | -3.44 dB | -11.2 |
| +87.89 kHz | -44.4 dBFS | -41.3 dBFS | -3.18 dB | -3.32 dB | -10.8 |
| +3.375 MHz | -45.4 dBFS | -42.3 dBFS | -3.16 dB | -3.30 dB | -10.7 |
| +3.31641 MHz | -45.4 dBFS | -42.4 dBFS | -3.00 dB | -3.14 dB | -10.2 |

Residual spectrum (dB, median offset removed). `A` is hotter above the axis.

```
 +8.3                                                                         
                                                                              
                                                                              
                                                                              
                                     #                                        
      ...............................|...##...................................
                                     |   ||                                   
          **   **    *  * *      ** *|  *||          *  *** * ***** **    ****
  0.0 -----------*----**-*-*----*--*--**------*--------*-----*----------------
      ****  ***   *|*       *|**           |** ||||** |    *       *  *||*    
                   *         |             |   *|||   |                ||     
      .......................*.............#....|||...*................#|.....
                                                |||                     #     
                                                #||                           
                                                 ||                           
                                                 ||                           
 -8.3                                            ##                           
      -3.20117 MHz                       0                        +3.20117 MHz
```

## Band +433 MHz

- Welch windows: A 975, B 975
- Median offset A−B: **+0.71 dB**
- Residual σ (MAD): **0.216 dB**; threshold **3.00 dB**
- Noise floor: A **-45.3 dBFS**, B **-46.0 dBFS** (Δ +0.71 dB)
- LO leakage (DC): A **3.5 dBFS**, B **-0.3 dBFS** (Δ +3.82 dB, residual +3.11 dB)
- RMS residual: **0.30 dB**
- Significant bins: **2**

| Offset | A | B | Δ | residual | σ |
| ---: | ---: | ---: | ---: | ---: | ---: |
| +1 MHz | -45.2 dBFS | -39.2 dBFS | -6.09 dB | -6.80 dB | -31.5 |
| -1.953 kHz | -2.5 dBFS | -6.3 dBFS | +3.82 dB | +3.11 dB | +14.4 |

Residual spectrum (dB, median offset removed). `A` is hotter above the axis.

```
 +6.8                                                                         
                                                                              
                                                                              
                                                                              
      ...................................##...................................
                                         ||                                   
                                         ||                                   
         **  *    ** * *** *     ***** * ||* * **** *   *   **** ***   *   *  
  0.0 -----*--*-------*-----------------------*----*--*----*--------------*---
      **|   *  ***  *     * **|**     * *   *        | * **     *   *** **  **
        |                     |                      |                        
        *                     *                      |                        
      ...............................................|........................
                                                     |                        
                                                     |                        
                                                     |                        
 -6.8                                                #                        
      -3.20117 MHz                       0                        +3.20117 MHz
```

## Band +915 MHz

- Welch windows: A 975, B 975
- Median offset A−B: **+0.81 dB**
- Residual σ (MAD): **0.208 dB**; threshold **3.00 dB**
- Noise floor: A **-45.2 dBFS**, B **-46.0 dBFS** (Δ +0.82 dB)
- LO leakage (DC): A **3.6 dBFS**, B **-0.2 dBFS** (Δ +3.80 dB, residual +2.98 dB)
- RMS residual: **0.32 dB**
- Significant bins: **3**

| Offset | A | B | Δ | residual | σ |
| ---: | ---: | ---: | ---: | ---: | ---: |
| +2 MHz | -43.8 dBFS | -38.3 dBFS | -5.48 dB | -6.29 dB | -30.3 |
| -3 MHz | -34.3 dBFS | -29.5 dBFS | -4.85 dB | -5.67 dB | -27.3 |
| +3 MHz | -44.8 dBFS | -40.8 dBFS | -4.06 dB | -4.88 dB | -23.5 |

Residual spectrum (dB, median offset removed). `A` is hotter above the axis.

```
 +6.3                                                                         
                                                                              
                                                                              
                                                                              
      ...................................**...................................
                                         ||                                   
                                         ||                                   
            * *** ** * * ** *   ** ***  *|| ** * * *  *  **  *     **    *    
  0.0 -*---*----------*-*-----------------------*---*------*--*------*--------
      * |**  *   *  *      * ***  *   **   *  *   *  | **   *  *|**   *** *|**
        |                                            |          |          |  
        |                                            *          |          |  
      ..|.......................................................|..........|..
        |                                                       |          |  
        |                                                       |          #  
        #                                                       |             
 -6.3                                                           #             
      -3.20117 MHz                       0                        +3.20117 MHz
```

## Band +2.45 GHz

- Welch windows: A 975, B 975
- Median offset A−B: **+1.01 dB**
- Residual σ (MAD): **1.328 dB**; threshold **7.97 dB**
- Noise floor: A **-41.7 dBFS**, B **-42.7 dBFS** (Δ +0.93 dB)
- LO leakage (DC): A **5.6 dBFS**, B **2.4 dBFS** (Δ +3.25 dB, residual +2.23 dB)
- RMS residual: **1.32 dB**
- Significant bins: **0**

Residual spectrum (dB, median offset removed). `A` is hotter above the axis.

```
+15.9                                                                         
                                                                              
                                                                              
                                                                              
      ........................................................................
                                                                              
      ****       *****  *                                                     
      ||||*******|||||**|*********** *   **                                   
  0.0 --------------------------------***---*************-*-------------------
                                    *      *             * *****|*************
                                                                *             
                                                                              
      ........................................................................
                                                                              
                                                                              
                                                                              
-15.9                                                                         
      -3.20117 MHz                       0                        +3.20117 MHz
```

