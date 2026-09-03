package grcode

// #cgo pkg-config: zbar
// #include <zbar.h>
import "C"

import (
	"fmt"
	"strings"
)

// following is reference from zbar.h
const (
	ZBAR_NONE    = 0      /**< no symbol decoded */
	ZBAR_PARTIAL = 1      /**< intermediate status */
	ZBAR_EAN8    = 8      /**< EAN-8 */
	ZBAR_UPCE    = 9      /**< UPC-E */
	ZBAR_ISBN10  = 10     /**< ISBN-10 (from EAN-13). @since 0.4 */
	ZBAR_UPCA    = 12     /**< UPC-A */
	ZBAR_EAN13   = 13     /**< EAN-13 */
	ZBAR_ISBN13  = 14     /**< ISBN-13 (from EAN-13). @since 0.4 */
	ZBAR_I25     = 25     /**< Interleaved 2 of 5. @since 0.4 */
	ZBAR_CODE39  = 39     /**< Code 39. @since 0.4 */
	ZBAR_PDF417  = 57     /**< PDF417. @since 0.6 */
	ZBAR_QRCODE  = 64     /**< QR Code. @since 0.10 */
	ZBAR_CODE128 = 128    /**< Code 128 */
	ZBAR_SYMBOL  = 0x00ff /**< mask for base symbol type */
	ZBAR_ADDON2  = 0x0200 /**< 2-digit add-on flag */
	ZBAR_ADDON5  = 0x0500 /**< 5-digit add-on flag */
	ZBAR_ADDON   = 0x0700 /**< add-on flag mask */
	NONE         = C.ZBAR_NONE
	PARTIAL      = C.ZBAR_PARTIAL
	EAN8         = C.ZBAR_EAN8
	UPCE         = C.ZBAR_UPCE
	ISBN10       = C.ZBAR_ISBN10
	UPCA         = C.ZBAR_UPCA
	EAN13        = C.ZBAR_EAN13
	ISBN13       = C.ZBAR_ISBN13
	I25          = C.ZBAR_I25
	CODE39       = C.ZBAR_CODE39
	PDF414       = C.ZBAR_PDF417
	QRCODE       = C.ZBAR_QRCODE
	CODE128      = C.ZBAR_CODE128
	SYMBOL       = C.ZBAR_SYMBOL
	ADDON2       = C.ZBAR_ADDON2
	ADDON5       = C.ZBAR_ADDON5
	ADDON        = C.ZBAR_ADDON
	XDensity     = C.ZBAR_CFG_X_DENSITY
	YDensity     = C.ZBAR_CFG_Y_DENSITY
)

type SymbolType struct {
	t C.zbar_symbol_type_t
}

var All = SymbolType{}

func (t SymbolType) IsQRCODE() bool { return t.t == QRCODE }

func ParseSymbolType(text string) (SymbolType, error) {
	var s C.zbar_symbol_type_t
	switch strings.ToLower(text) {
	case "ean8":
		s = EAN8
	case "upce":
		s = UPCE
	case "isbn10":
		s = ISBN10
	case "upca":
		s = UPCA
	case "ean13":
		s = EAN13
	case "isbn13":
		s = ISBN13
	case "i25":
		s = I25
	case "code39":
		s = CODE39
	case "pdf414":
		s = PDF414
	case "qrcode":
		s = QRCODE
	case "code128":
		s = CODE128
	case "symbol":
		s = SYMBOL
	default:
		return SymbolType{}, fmt.Errorf("unknown symbol type %s", text)
	}
	return SymbolType{t: s}, nil
}
