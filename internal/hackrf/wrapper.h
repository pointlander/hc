#ifndef HC_HACKRF_WRAPPER_H
#define HC_HACKRF_WRAPPER_H

#include "libhackrf/hackrf.h"
#include <stdint.h>

int hc_rx_cb(hackrf_transfer *transfer);
int hc_tx_cb(hackrf_transfer *transfer);

int hc_start_rx(hackrf_device *device, uintptr_t ctx);
int hc_start_tx(hackrf_device *device, uintptr_t ctx);

int hc_hackrf_list_count(hackrf_device_list_t *list);
int hc_hackrf_list_usb_count(hackrf_device_list_t *list);
const char *hc_hackrf_list_serial(hackrf_device_list_t *list, int idx);
int hc_hackrf_list_board_id(hackrf_device_list_t *list, int idx);

#endif
