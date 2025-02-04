package grcode

// #cgo pkg-config: zbar
// #include <zbar.h>
import "C"
import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
)

type RawData string

func GetDataFromFile(imagePath string) (results []string, err error) {
	// TODO: read via libjpeg, libpng instead of Go
	//filePath := C.CString(imagePath)
	reader, err := os.Open(imagePath)
	if err != nil {
		log.Printf("open file error: %v", err)
		return results, err
	}
	defer reader.Close()
	m, _, err := image.Decode(reader)
	if err != nil {
		log.Printf("decode file error: %v", err)
		return results, err
	}
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
