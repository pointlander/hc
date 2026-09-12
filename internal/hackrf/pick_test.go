package hackrf

import "testing"

func TestPickTwoDefault(t *testing.T) {
	devs := []DeviceInfo{
		{Index: 0, Serial: "aaaa", USBBoardID: USBBoardIDHackRFOne},
		{Index: 1, Serial: "bbbb", USBBoardID: USBBoardIDHackRFOne},
	}
	a, b, err := PickTwo(devs, "", "")
	if err != nil {
		t.Fatal(err)
	}
	if a.Serial != "aaaa" || b.Serial != "bbbb" {
		t.Fatalf("a=%s b=%s", a.Serial, b.Serial)
	}
}

func TestPickTwoBySuffix(t *testing.T) {
	devs := []DeviceInfo{
		{Index: 0, Serial: "0000aaaa", USBBoardID: USBBoardIDHackRFOne},
		{Index: 1, Serial: "0000bbbb", USBBoardID: USBBoardIDHackRFOne},
	}
	a, b, err := PickTwo(devs, "bbbb", "aaaa")
	if err != nil {
		t.Fatal(err)
	}
	if a.Serial != "0000bbbb" || b.Serial != "0000aaaa" {
		t.Fatalf("a=%s b=%s", a.Serial, b.Serial)
	}
}

func TestPickTwoRejectsSame(t *testing.T) {
	devs := []DeviceInfo{
		{Index: 0, Serial: "aaaa", USBBoardID: USBBoardIDHackRFOne},
		{Index: 1, Serial: "bbbb", USBBoardID: USBBoardIDHackRFOne},
	}
	_, _, err := PickTwo(devs, "aaaa", "aaaa")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestPickTwoNeedsTwo(t *testing.T) {
	devs := []DeviceInfo{{Index: 0, Serial: "aaaa", USBBoardID: USBBoardIDHackRFOne}}
	_, _, err := PickTwo(devs, "", "")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestFindSkipsEmptySerialSuffix(t *testing.T) {
	devs := []DeviceInfo{
		{Index: 0, Serial: "", USBBoardID: USBBoardIDHackRFOne},
		{Index: 1, Serial: "f75461dc288e2dc3", USBBoardID: USBBoardIDHackRFOne},
	}
	d, ok := Find(devs, "2dc3")
	if !ok || d.Index != 1 {
		t.Fatalf("got %+v ok=%v", d, ok)
	}
}

func TestLabelTruncatesSerial(t *testing.T) {
	d := DeviceInfo{Serial: "0000000000000000f75461dc288e2dc3"}
	if g := d.Label(); g != "288e2dc3" {
		t.Fatalf("label=%q", g)
	}
}
