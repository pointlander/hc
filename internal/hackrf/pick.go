package hackrf

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// Find returns the first device matching want, which may be a serial,
// a serial suffix, or a decimal list index.
func Find(devs []DeviceInfo, want string) (DeviceInfo, bool) {
	if want == "" {
		return DeviceInfo{}, false
	}
	for _, d := range devs {
		if d.Serial != "" && (d.Serial == want || strings.HasSuffix(d.Serial, want) || strings.HasSuffix(want, d.Serial)) {
			return d, true
		}
		if d.USBPath == want {
			return d, true
		}
		if d.Index >= 0 && strconv.Itoa(d.Index) == want {
			return d, true
		}
	}
	return DeviceInfo{}, false
}

func usableSorted(devs []DeviceInfo) []DeviceInfo {
	usable := make([]DeviceInfo, 0, len(devs))
	for _, d := range devs {
		if d.DFU {
			continue
		}
		usable = append(usable, d)
	}
	sort.SliceStable(usable, func(i, j int) bool {
		if usable[i].Serial != usable[j].Serial {
			return usable[i].Serial < usable[j].Serial
		}
		return usable[i].Index < usable[j].Index
	})
	return usable
}

func sameDevice(a, b DeviceInfo) bool {
	if a.Serial != "" && b.Serial != "" {
		return a.Serial == b.Serial
	}
	if a.Index >= 0 && b.Index >= 0 {
		return a.Index == b.Index
	}
	if a.USBPath != "" && b.USBPath != "" {
		return a.USBPath == b.USBPath
	}
	return false
}

// PickTwo chooses two distinct HackRFs for simultaneous receive.
func PickTwo(devs []DeviceInfo, aWant, bWant string) (a DeviceInfo, b DeviceInfo, err error) {
	usable := usableSorted(devs)
	if aWant != "" && bWant != "" {
		var ok bool
		a, ok = Find(usable, aWant)
		if !ok {
			return DeviceInfo{}, DeviceInfo{}, fmt.Errorf("HackRF A %q not found", aWant)
		}
		b, ok = Find(usable, bWant)
		if !ok {
			return DeviceInfo{}, DeviceInfo{}, fmt.Errorf("HackRF B %q not found", bWant)
		}
		if sameDevice(a, b) {
			return DeviceInfo{}, DeviceInfo{}, fmt.Errorf("radios A and B must be different HackRFs")
		}
		return a, b, nil
	}
	if len(usable) < 2 {
		nUSB := len(ScanUSB())
		return DeviceInfo{}, DeviceInfo{}, fmt.Errorf("need two HackRF radios; found %d (usb scan %d)", len(usable), nUSB)
	}
	if aWant != "" {
		var ok bool
		a, ok = Find(usable, aWant)
		if !ok {
			return DeviceInfo{}, DeviceInfo{}, fmt.Errorf("HackRF A %q not found", aWant)
		}
		for _, d := range usable {
			if !sameDevice(d, a) {
				return a, d, nil
			}
		}
		return DeviceInfo{}, DeviceInfo{}, fmt.Errorf("no second HackRF for B")
	}
	if bWant != "" {
		var ok bool
		b, ok = Find(usable, bWant)
		if !ok {
			return DeviceInfo{}, DeviceInfo{}, fmt.Errorf("HackRF B %q not found", bWant)
		}
		for _, d := range usable {
			if !sameDevice(d, b) {
				return d, b, nil
			}
		}
		return DeviceInfo{}, DeviceInfo{}, fmt.Errorf("no second HackRF for A")
	}
	return usable[0], usable[1], nil
}

// Label is a short identifier for logs.
func (d DeviceInfo) Label() string {
	if d.Serial != "" {
		s := strings.TrimLeft(d.Serial, "0")
		if s == "" {
			s = "0"
		}
		if len(s) > 8 {
			s = s[len(s)-8:]
		}
		return s
	}
	if d.USBPath != "" {
		return d.USBPath
	}
	if d.Index >= 0 {
		return fmt.Sprintf("#%d", d.Index)
	}
	return "?"
}
