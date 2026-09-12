#ifndef __HACKRF_H__
#define __HACKRF_H__

/*
 * Minimal libhackrf C API surface, vendored so this package can compile
 * against the runtime library (libhackrf.so.0) without libhackrf-dev.
 * Layout of hackrf_transfer and hackrf_device_list must match libhackrf.
 *
 * Copyright (c) 2012-2022 Great Scott Gadgets
 * BSD-3-Clause; see https://github.com/greatscottgadgets/hackrf
 */

#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

enum hackrf_error {
	HACKRF_SUCCESS = 0,
	HACKRF_TRUE = 1,
	HACKRF_ERROR_INVALID_PARAM = -2,
	HACKRF_ERROR_NOT_FOUND = -5,
	HACKRF_ERROR_BUSY = -6,
	HACKRF_ERROR_NO_MEM = -11,
	HACKRF_ERROR_LIBUSB = -1000,
	HACKRF_ERROR_THREAD = -1001,
	HACKRF_ERROR_STREAMING_THREAD_ERR = -1002,
	HACKRF_ERROR_STREAMING_STOPPED = -1003,
	HACKRF_ERROR_STREAMING_EXIT_CALLED = -1004,
	HACKRF_ERROR_USB_API_VERSION = -1005,
	HACKRF_ERROR_NOT_LAST_DEVICE = -2000,
	HACKRF_ERROR_OTHER = -9999,
};

enum hackrf_usb_board_id {
	USB_BOARD_ID_JAWBREAKER = 0x604B,
	USB_BOARD_ID_HACKRF_ONE = 0x6089,
	USB_BOARD_ID_RAD1O = 0xCC15,
	USB_BOARD_ID_INVALID = 0xFFFF,
};

typedef struct hackrf_device hackrf_device;

typedef struct {
	hackrf_device *device;
	uint8_t *buffer;
	int buffer_length;
	int valid_length;
	void *rx_ctx;
	void *tx_ctx;
} hackrf_transfer;

struct hackrf_device_list {
	char **serial_numbers;
	enum hackrf_usb_board_id *usb_board_ids;
	int *usb_device_index;
	int devicecount;
	void **usb_devices;
	int usb_devicecount;
};

typedef struct hackrf_device_list hackrf_device_list_t;
typedef int (*hackrf_sample_block_cb_fn)(hackrf_transfer *transfer);

extern int hackrf_init(void);
extern int hackrf_exit(void);
extern const char *hackrf_library_version(void);
extern const char *hackrf_library_release(void);
extern const char *hackrf_error_name(enum hackrf_error errcode);

extern hackrf_device_list_t *hackrf_device_list(void);
extern int hackrf_device_list_open(hackrf_device_list_t *list, int idx, hackrf_device **device);
extern void hackrf_device_list_free(hackrf_device_list_t *list);

extern int hackrf_open(hackrf_device **device);
extern int hackrf_open_by_serial(const char *desired_serial_number, hackrf_device **device);
extern int hackrf_close(hackrf_device *device);

extern int hackrf_start_rx(hackrf_device *device, hackrf_sample_block_cb_fn callback, void *rx_ctx);
extern int hackrf_stop_rx(hackrf_device *device);
extern int hackrf_start_tx(hackrf_device *device, hackrf_sample_block_cb_fn callback, void *tx_ctx);
extern int hackrf_stop_tx(hackrf_device *device);
extern int hackrf_is_streaming(hackrf_device *device);

extern int hackrf_set_baseband_filter_bandwidth(hackrf_device *device, const uint32_t bandwidth_hz);
extern int hackrf_set_freq(hackrf_device *device, const uint64_t freq_hz);
extern int hackrf_set_sample_rate(hackrf_device *device, const double freq_hz);
extern int hackrf_set_amp_enable(hackrf_device *device, const uint8_t value);
extern int hackrf_set_lna_gain(hackrf_device *device, uint32_t value);
extern int hackrf_set_vga_gain(hackrf_device *device, uint32_t value);
extern int hackrf_set_txvga_gain(hackrf_device *device, uint32_t value);
extern int hackrf_set_antenna_enable(hackrf_device *device, const uint8_t value);
extern uint32_t hackrf_compute_baseband_filter_bw(const uint32_t bandwidth_hz);

extern int hackrf_version_string_read(hackrf_device *device, char *version, uint8_t length);
extern int hackrf_usb_api_version_read(hackrf_device *device, uint16_t *version);

#ifdef __cplusplus
}
#endif

#endif
