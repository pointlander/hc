package hackrf

import "testing"

func TestListDevices(t *testing.T) {
	if err := Init(); err != nil {
		t.Skipf("hackrf init: %v", err)
	}
	t.Cleanup(func() { _ = Exit() })

	devs, err := List()
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(devs) < 2 {
		t.Skipf("need two HackRFs, found %d", len(devs))
	}
	for _, d := range devs {
		if d.USBBoardID != USBBoardIDHackRFOne {
			t.Errorf("unexpected board %s serial=%s", d.USBBoardID, d.Serial)
		}
		if d.Serial == "" {
			t.Errorf("device %d has empty serial", d.Index)
		}
	}
}

func TestScanUSBNoPanic(t *testing.T) {
	_ = ScanUSB()
}
