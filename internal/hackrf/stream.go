package hackrf

/*
#include "wrapper.h"
*/
import "C"

import (
	"runtime/cgo"
	"unsafe"
)

// RXCallback is invoked from the libhackrf transfer thread with a slice of
// interleaved int8 I/Q bytes. The slice is only valid for the duration of
// the call; copy if the data must outlive the callback.
type RXCallback func(buf []byte) error

// TXCallback fills buf with interleaved int8 I/Q bytes to transmit.
type TXCallback func(buf []byte) error

//export hcGoRx
func hcGoRx(transfer *C.hackrf_transfer) C.int {
	h := cgo.Handle(uintptr(transfer.rx_ctx))
	cb, ok := h.Value().(RXCallback)
	if !ok || cb == nil {
		return -1
	}
	n := int(transfer.valid_length)
	if n <= 0 {
		return 0
	}
	buf := unsafe.Slice((*byte)(transfer.buffer), n)
	if err := cb(buf); err != nil {
		return -1
	}
	return 0
}

//export hcGoTx
func hcGoTx(transfer *C.hackrf_transfer) C.int {
	h := cgo.Handle(uintptr(transfer.tx_ctx))
	cb, ok := h.Value().(TXCallback)
	if !ok || cb == nil {
		return -1
	}
	n := int(transfer.buffer_length)
	if n <= 0 {
		return 0
	}
	buf := unsafe.Slice((*byte)(transfer.buffer), n)
	if err := cb(buf); err != nil {
		return -1
	}
	transfer.valid_length = C.int(n)
	return 0
}

// StartRX begins receiving interleaved int8 I/Q samples.
func (d *Device) StartRX(cb RXCallback) error {
	if d.hasRX {
		d.rxH.Delete()
		d.hasRX = false
	}
	d.rxH = cgo.NewHandle(cb)
	d.hasRX = true
	if err := toError(C.hc_start_rx(d.dev, C.uintptr_t(d.rxH))); err != nil {
		d.rxH.Delete()
		d.hasRX = false
		return err
	}
	return nil
}

// StopRX stops the receive stream.
func (d *Device) StopRX() error {
	if d.dev == nil {
		return nil
	}
	err := toError(C.hackrf_stop_rx(d.dev))
	if d.hasRX {
		d.rxH.Delete()
		d.hasRX = false
	}
	return err
}

// StartTX begins transmitting interleaved int8 I/Q samples supplied by cb.
func (d *Device) StartTX(cb TXCallback) error {
	if d.hasTX {
		d.txH.Delete()
		d.hasTX = false
	}
	d.txH = cgo.NewHandle(cb)
	d.hasTX = true
	if err := toError(C.hc_start_tx(d.dev, C.uintptr_t(d.txH))); err != nil {
		d.txH.Delete()
		d.hasTX = false
		return err
	}
	return nil
}

// StopTX stops the transmit stream.
func (d *Device) StopTX() error {
	if d.dev == nil {
		return nil
	}
	err := toError(C.hackrf_stop_tx(d.dev))
	if d.hasTX {
		d.txH.Delete()
		d.hasTX = false
	}
	return err
}
