package grcode

// #cgo darwin pkg-config: zbar
// #cgo LDFLAGS: -lzbar
// #include <zbar.h>
import "C"

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
)

type RawData string

func GetDataFromFile(imagePath string) (results []string, err error) {
	img, err := OpenImage(imagePath)
	if err != nil {
		return nil, err
	}
	return GetDataFromImage(img)
}

func OpenImage(imagePath string) (image.Image, error) {
	fh, err := os.Open(imagePath)
	if err != nil {
		return nil, err
	}
	defer fh.Close()
	return ImageFromReader(fh)
}

func ImageFromReader(r io.Reader) (image.Image, error) {
	m, _, err := image.Decode(r)
	if err != nil {
		err = fmt.Errorf("decode: %w", err)
	}
	return m, err
}

func GetDataFromReader(r io.Reader) (results []string, err error) {
	m, err := ImageFromReader(r)
	return GetDataFromImage(m)
}

// GetDataFromImage read qrcode directly from golang Image class
func GetDataFromImage(image image.Image) (results []string, err error) {
	scanner := NewScanner()
	defer scanner.Close()
	scanner.SetConfig(0, C.ZBAR_CFG_ENABLE, 1)
	zImg := NewZbarImage(image)
	defer zImg.Close()
	scanner.Scan(zImg)
	symbol := zImg.GetSymbol()
	for ; symbol != nil; symbol = symbol.Next() {
		results = append(results, symbol.Data())
	}
	return results, nil
}
