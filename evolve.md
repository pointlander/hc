# Evolved Casimir center sheet

Evolved 2026-09-18 00:26:22 UTC.

A genetic algorithm searched binary metal occupancy on a rectangular lattice. The center sheet may **overhang the anodized plates**: only metal under the plates sees the 1 µm MIM gap; overhang is still connected (higher inductance, no MIM C). Ranking fitness is **⟨P₅₀⟩·A** (mean available thermal power into 50 Ω times total metal area). Power is evaluated at 100 MHz, 433 MHz, 915 MHz, and 2.45 GHz. Sheet span is a gene; plate size is fixed. The E-shaped seed is the baseline.

| | |
| --- | ---: |
| Grid | 16 × 16 |
| Population / generations | 40 / 60 |
| Mutation | 4% bits / generation |
| Plates (centered, square) | **50.0 mm** |
| Sheet span range | 2.0–250 mm |

## Best sheet

| | |
| --- | ---: |
| Sheet span | **247.07 mm** |
| Plates | 50.00 mm |
| Fill | **60.9%** (156 cells) |
| Metal area | 37198.223 mm² |
| Overlap (under plates) | 2622.952 mm² |
| Overhang | 34575.271 mm² |
| Mean P(50 Ω) | **-210.11 dBm/Hz** |
| Fitness ⟨P⟩·A | **3.626e-26** W·m²/Hz |
| E-shape baseline P | -207.47 dBm/Hz |
| E-shape baseline ⟨P⟩·A | 3.115e-27 W·m²/Hz |
| Fitness gain | **×11.64** |

Band-by-band available power:

| f | Re(Z) | Im(Z) | P(50 Ω) |
| ---: | ---: | ---: | ---: |
| 100.00 MHz | 0.00728 Ω | -0.00677 Ω | **-206.3 dBm/Hz** |
| 433.00 MHz | 0.00322 Ω | -0.005 Ω | **-209.8 dBm/Hz** |
| 915.00 MHz | 0.0014 Ω | -0.00366 Ω | **-213.4 dBm/Hz** |
| 2.45 GHz | 0.000153 Ω | -0.00162 Ω | **-223.0 dBm/Hz** |

Plan view (`#` = metal under the plates, `+` = overhang, `.` = empty):

```
++..........++++
+++++..+.+..++++
++.+..++.+...++.
++.+.+.+.+...+++
.+++++++.+++..++
+++..+.+++.+.++.
++++++.#..+..+++
+.....#.##+...+.
+++...##.#+..+++
.+++++####+.++.+
++++..+.+.++++++
..+++...+.+..+++
.++.+..++.+.+++.
.++.++.+.+++.+++
..+.+++++.+++++.
...++.++.+++++..
```

![Evolved center sheet](sheet.png)

![Evolved sheet, studio view](sheet_photo.jpg)
