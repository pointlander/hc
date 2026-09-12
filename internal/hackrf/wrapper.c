#include "wrapper.h"
#include <stddef.h>

extern int hcGoRx(hackrf_transfer *transfer);
extern int hcGoTx(hackrf_transfer *transfer);

int hc_rx_cb(hackrf_transfer *transfer) {
	return hcGoRx(transfer);
}

int hc_tx_cb(hackrf_transfer *transfer) {
	return hcGoTx(transfer);
}

int hc_start_rx(hackrf_device *device, uintptr_t ctx) {
	return hackrf_start_rx(device, hc_rx_cb, (void *)ctx);
}

int hc_start_tx(hackrf_device *device, uintptr_t ctx) {
	return hackrf_start_tx(device, hc_tx_cb, (void *)ctx);
}

int hc_hackrf_list_count(hackrf_device_list_t *list) {
	if (list == NULL) {
		return 0;
	}
	return list->devicecount;
}

int hc_hackrf_list_usb_count(hackrf_device_list_t *list) {
	if (list == NULL) {
		return 0;
	}
	return list->usb_devicecount;
}

const char *hc_hackrf_list_serial(hackrf_device_list_t *list, int idx) {
	if (list == NULL || idx < 0 || idx >= list->devicecount || list->serial_numbers == NULL) {
		return NULL;
	}
	return list->serial_numbers[idx];
}

int hc_hackrf_list_board_id(hackrf_device_list_t *list, int idx) {
	if (list == NULL || idx < 0 || idx >= list->devicecount || list->usb_board_ids == NULL) {
		return (int)USB_BOARD_ID_INVALID;
	}
	return (int)list->usb_board_ids[idx];
}
