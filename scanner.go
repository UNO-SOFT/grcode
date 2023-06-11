package grcode

// #cgo LDFLAGS: -lzbar
// #include <zbar.h>
import "C"

import (
	"runtime"
	"strconv"
	"unsafe"
)

// Scanner is wrapper of zbar_image_scanner
type Scanner struct {
	image_scanner *C.zbar_image_scanner_t
}

// NewScanner returns new Scanner
func NewScanner() *Scanner {
	r := &Scanner{image_scanner: C.zbar_image_scanner_create()}
	// runtime.SetFinalizer() works well for automatically free()'ing cgo memory allocations!
	// the finalizer will be called when the garbage collector is invoked.
	// we can derfer destroy on our own
	runtime.SetFinalizer(r, (*Scanner).Close)
	return r
}

// SetConfig sets scanner config
// 0 for success, non-0 for failure
func (s *Scanner) SetConfig(symbology C.zbar_symbol_type_t, config C.zbar_config_t, value int) int {
	return int(C.zbar_image_scanner_set_config(s.image_scanner, symbology, config, C.int(value)))
}

// Scan scans image, the image format must be "Y800" or "GRAY".
// returns:
// >0, if symbols were successfully decoded from the image
//
//	0, if no symbols were found or
//
// -1, if an error occurs
func (s *Scanner) Scan(img *ZbarImage) (int, error) {
	rc := int(C.zbar_scan_image(s.image_scanner, img.image))
	if rc < 0 {
		return rc, &zbarError{
			Code:    int(C.int(C._zbar_get_error_code(unsafe.Pointer(s.image_scanner)))),
			Message: C.GoString(C._zbar_error_string(unsafe.Pointer(s.image_scanner), 1)),
		}
	}
	return rc, nil
}

// Close suicides
func (s *Scanner) Close() {
	if s != nil && s.image_scanner != nil {
		C.zbar_image_scanner_destroy(s.image_scanner) // void function
		s.image_scanner = nil
	}
}

type zbarError struct {
	Code    int
	Message string
}

func (ze *zbarError) Error() string {
	if ze == nil {
		return ""
	}
	b := make([]byte, 0, 128)
	b = strconv.AppendInt(b, int64(ze.Code), 10)
	b = append(b, ": "...)
	b = append(b, ze.Message...)
	return string(b)
}

var _ = error((*zbarError)(nil))
