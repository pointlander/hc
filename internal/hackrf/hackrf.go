package hackrf

/*
#cgo CFLAGS: -I${SRCDIR}
#cgo LDFLAGS: -l:libhackrf.so.0
#include "wrapper.h"
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"runtime/cgo"
	"unsafe"
)

// USBBoardID identifies the HackRF USB product.
type USBBoardID uint16

const (
	USBBoardIDJawbreaker USBBoardID = 0x604B
	USBBoardIDHackRFOne  USBBoardID = 0x6089
	USBBoardIDRad1o      USBBoardID = 0xCC15
	USBBoardIDInvalid    USBBoardID = 0xFFFF
)

func (u USBBoardID) String() string {
	switch u {
	case USBBoardIDJawbreaker:
		return "Jawbreaker"
	case USBBoardIDHackRFOne:
		return "HackRF One"
	case USBBoardIDRad1o:
		return "rad1o"
	case USBBoardIDInvalid:
		return "invalid"
	default:
		return fmt.Sprintf("unknown(0x%04x)", uint16(u))
	}
}

// DeviceInfo is a connected HackRF discovered before opening.
type DeviceInfo struct {
	Serial     string
	USBBoardID USBBoardID
	Index      int    // index in libhackrf's device list; -1 if sysfs-only
	USBPath    string // sysfs USB path, e.g. "3-1"
	DFU        bool   // true if the board is in NXP DFU (unusable until it enumerates as HackRF)
}

func toError(r C.int) error {
	if r == C.HACKRF_SUCCESS {
		return nil
	}
	name := C.GoString(C.hackrf_error_name(C.enum_hackrf_error(r)))
	if name == "" {
		return fmt.Errorf("hackrf: error %d", int(r))
	}
	return fmt.Errorf("hackrf: %s (%d)", name, int(r))
}

// Init initializes libhackrf. Call once at process start.
func Init() error {
	return toError(C.hackrf_init())
}

// Exit releases libhackrf. All devices must be closed first.
func Exit() error {
	return toError(C.hackrf_exit())
}

// LibraryVersion returns the linked libhackrf version string.
func LibraryVersion() string {
	return C.GoString(C.hackrf_library_version())
}

// LibraryRelease returns the linked libhackrf release string.
func LibraryRelease() string {
	return C.GoString(C.hackrf_library_release())
}

// List returns connected HackRF devices. It initializes libhackrf if needed,
// then enumerates via libusb. If libusb returns nothing, a sysfs USB scan is
// used as a fallback (kernel driver / context issues).
func List() ([]DeviceInfo, error) {
	if err := Init(); err != nil {
		return nil, err
	}
	devs, err := listLibUSB()
	if err != nil {
		return nil, err
	}
	usb := ScanUSB()
	if len(devs) == 0 && len(usb) > 0 {
		_ = Exit()
		if err := Init(); err != nil {
			return usb, nil
		}
		devs, err = listLibUSB()
		if err != nil || len(devs) == 0 {
			return usb, nil
		}
	}
	return mergeUSB(devs, usb), nil
}

func listLibUSB() ([]DeviceInfo, error) {
	clist := C.hackrf_device_list()
	if clist == nil {
		return nil, fmt.Errorf("hackrf: device list is nil")
	}
	defer C.hackrf_device_list_free(clist)

	n := int(C.hc_hackrf_list_count(clist))
	if n <= 0 {
		return nil, nil
	}
	out := make([]DeviceInfo, 0, n)
	for i := 0; i < n; i++ {
		info := DeviceInfo{
			USBBoardID: USBBoardID(C.hc_hackrf_list_board_id(clist, C.int(i))),
			Index:      i,
		}
		if s := C.hc_hackrf_list_serial(clist, C.int(i)); s != nil {
			info.Serial = C.GoString(s)
		}
		out = append(out, info)
	}
	return out, nil
}

func mergeUSB(lib, sys []DeviceInfo) []DeviceInfo {
	if len(sys) == 0 {
		return lib
	}
	used := map[string]bool{}
	for i := range lib {
		for _, s := range sys {
			if s.Serial != "" && s.Serial == lib[i].Serial {
				lib[i].USBPath = s.USBPath
				lib[i].DFU = s.DFU
				used[s.USBPath] = true
				break
			}
		}
	}
	for _, s := range sys {
		if used[s.USBPath] {
			continue
		}
		matched := false
		if s.Serial != "" {
			for _, d := range lib {
				if d.Serial == s.Serial {
					matched = true
					break
				}
			}
		}
		if !matched {
			lib = append(lib, s)
		}
	}
	return lib
}

// OpenDevice opens a HackRF identified by List(). Prefers serial, then list index.
func OpenDevice(info DeviceInfo) (*Device, error) {
	if info.DFU {
		return nil, fmt.Errorf("hackrf: device at %s is in DFU mode", info.USBPath)
	}
	if info.Serial != "" {
		d, err := OpenBySerial(info.Serial)
		if err == nil {
			d.serial = info.Serial
			return d, nil
		}
		if info.Index < 0 {
			return nil, err
		}
	}
	if info.Index >= 0 {
		d, err := OpenIndex(info.Index)
		if err == nil && info.Serial != "" {
			d.serial = info.Serial
		}
		return d, err
	}
	return Open()
}

// Open opens the first HackRF on the system.
func Open() (*Device, error) {
	var dev *C.hackrf_device
	if err := toError(C.hackrf_open(&dev)); err != nil {
		return nil, err
	}
	return &Device{dev: dev}, nil
}

// Device is an opened HackRF.
type Device struct {
	dev    *C.hackrf_device
	serial string
	rxH    cgo.Handle
	txH    cgo.Handle
	hasRX  bool
	hasTX  bool
	closed bool
}

// OpenBySerial opens a HackRF whose serial ends with the given string.
// An empty serial opens the first device.
func OpenBySerial(serial string) (*Device, error) {
	var dev *C.hackrf_device
	var cserial *C.char
	if serial != "" {
		cserial = C.CString(serial)
		defer C.free(unsafe.Pointer(cserial))
	}
	if err := toError(C.hackrf_open_by_serial(cserial, &dev)); err != nil {
		return nil, err
	}
	return &Device{dev: dev, serial: serial}, nil
}

// OpenIndex opens the device at the given index in List().
func OpenIndex(idx int) (*Device, error) {
	clist := C.hackrf_device_list()
	if clist == nil {
		return nil, fmt.Errorf("hackrf: device list is nil")
	}
	defer C.hackrf_device_list_free(clist)
	n := int(C.hc_hackrf_list_count(clist))
	if idx < 0 || idx >= n {
		return nil, fmt.Errorf("hackrf: index %d out of range (0..%d)", idx, n-1)
	}
	var dev *C.hackrf_device
	if err := toError(C.hackrf_device_list_open(clist, C.int(idx), &dev)); err != nil {
		return nil, err
	}
	d := &Device{dev: dev}
	if s := C.hc_hackrf_list_serial(clist, C.int(idx)); s != nil {
		d.serial = C.GoString(s)
	}
	return d, nil
}

// Serial returns the serial used to open the device, if known.
func (d *Device) Serial() string { return d.serial }

// Close stops streaming if needed and releases the device.
func (d *Device) Close() error {
	if d == nil || d.dev == nil || d.closed {
		return nil
	}
	_ = d.StopRX()
	_ = d.StopTX()
	err := toError(C.hackrf_close(d.dev))
	d.closed = true
	d.dev = nil
	return err
}

// Version reads the firmware version string.
func (d *Device) Version() (string, error) {
	buf := (*C.char)(C.malloc(256))
	defer C.free(unsafe.Pointer(buf))
	if err := toError(C.hackrf_version_string_read(d.dev, buf, 255)); err != nil {
		return "", err
	}
	return C.GoString(buf), nil
}

// USBAPIVersion reads the firmware USB API version.
func (d *Device) USBAPIVersion() (uint16, error) {
	var v C.uint16_t
	if err := toError(C.hackrf_usb_api_version_read(d.dev, &v)); err != nil {
		return 0, err
	}
	return uint16(v), nil
}

// SetFreq sets the center frequency in Hz.
func (d *Device) SetFreq(hz uint64) error {
	return toError(C.hackrf_set_freq(d.dev, C.uint64_t(hz)))
}

// SetSampleRate sets the ADC/DAC sample rate in Hz.
func (d *Device) SetSampleRate(hz float64) error {
	return toError(C.hackrf_set_sample_rate(d.dev, C.double(hz)))
}

// SetBasebandFilterBandwidth sets the analog baseband filter in Hz.
func (d *Device) SetBasebandFilterBandwidth(hz uint32) error {
	return toError(C.hackrf_set_baseband_filter_bandwidth(d.dev, C.uint32_t(hz)))
}

// ComputeBasebandFilterBW returns the nearest supported filter bandwidth.
func ComputeBasebandFilterBW(hz uint32) uint32 {
	return uint32(C.hackrf_compute_baseband_filter_bw(C.uint32_t(hz)))
}

// SetAmpEnable turns the RF amplifier (~14 dB) on or off.
func (d *Device) SetAmpEnable(on bool) error {
	var v C.uint8_t
	if on {
		v = 1
	}
	return toError(C.hackrf_set_amp_enable(d.dev, v))
}

// SetLNAGain sets RX IF/LNA gain. Range 0-40 dB in 8 dB steps.
func (d *Device) SetLNAGain(db int) error {
	return toError(C.hackrf_set_lna_gain(d.dev, C.uint32_t(db)))
}

// SetVGAGain sets RX baseband VGA gain. Range 0-62 dB in 2 dB steps.
func (d *Device) SetVGAGain(db int) error {
	return toError(C.hackrf_set_vga_gain(d.dev, C.uint32_t(db)))
}

// SetTXVGAGain sets TX IF gain. Range 0-47 dB in 1 dB steps.
func (d *Device) SetTXVGAGain(db int) error {
	return toError(C.hackrf_set_txvga_gain(d.dev, C.uint32_t(db)))
}

// SetAntennaEnable turns antenna port power (bias tee) on or off.
func (d *Device) SetAntennaEnable(on bool) error {
	var v C.uint8_t
	if on {
		v = 1
	}
	return toError(C.hackrf_set_antenna_enable(d.dev, v))
}

// IsStreaming reports whether the device is currently transferring samples.
func (d *Device) IsStreaming() bool {
	return C.hackrf_is_streaming(d.dev) == C.HACKRF_TRUE
}
