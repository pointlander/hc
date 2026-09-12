package hackrf

import (
	"os"
	"path/filepath"
	"strings"
)

const sysUSB = "/sys/bus/usb/devices"

// ScanUSB finds HackRF-like USB devices via sysfs, independent of libhackrf.
// This catches devices libusb missed (wrong context, kernel driver, etc.).
func ScanUSB() []DeviceInfo {
	ents, err := os.ReadDir(sysUSB)
	if err != nil {
		return nil
	}
	var out []DeviceInfo
	for _, e := range ents {
		name := e.Name()
		if strings.Contains(name, ":") {
			continue
		}
		dir := filepath.Join(sysUSB, name)
		vid := readSysFile(filepath.Join(dir, "idVendor"))
		pid := readSysFile(filepath.Join(dir, "idProduct"))
		board, dfu, ok := usbBoard(vid, pid)
		if !ok {
			continue
		}
		serial := readSysFile(filepath.Join(dir, "serial"))
		out = append(out, DeviceInfo{
			Serial:     serial,
			USBBoardID: board,
			Index:      -1,
			USBPath:    name,
			DFU:        dfu,
		})
	}
	return out
}

func usbBoard(vid, pid string) (USBBoardID, bool, bool) {
	vid = strings.ToLower(vid)
	pid = strings.ToLower(pid)
	if vid == "1fc9" && pid == "000c" {
		return USBBoardIDInvalid, true, true
	}
	if vid != "1d50" {
		return 0, false, false
	}
	switch pid {
	case "604b":
		return USBBoardIDJawbreaker, false, true
	case "6089":
		return USBBoardIDHackRFOne, false, true
	case "cc15":
		return USBBoardIDRad1o, false, true
	default:
		return 0, false, false
	}
}

func readSysFile(path string) string {
	b, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(b))
}
