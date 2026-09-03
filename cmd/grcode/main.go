package main

import (
	"context"
	"flag"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"maps"
	"os"
	"os/signal"
	"slices"
	"sync"
	"time"

	"github.com/UNO-SOFT/grcode"
	"github.com/disintegration/imaging"
	"golang.org/x/sync/errgroup"
)

// go build -ldflags "-linkmode external -extldflags -static"
func main() {
	if err := Main(); err != nil {
		log.Fatal(err)
	}
}

func Main() error {
	flag.Usage = func() {
		flag.PrintDefaults()
	}
	flag.Parse()
	//log.SetFlags(0)
	if flag.NArg() < 1 {
		return fmt.Errorf("Need specify the image file")
	}
	filePath := flag.Arg(0)

	img, err := grcode.OpenImage(filePath)
	if err != nil {
		return err
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	result := make(map[int][]string, 8)
	var resultMu sync.Mutex
	grp, _ := errgroup.WithContext(ctx)
	grp.SetLimit(4)
	i90 := imaging.Rotate90(img)
	i270 := imaging.Rotate270(img)
	for a, img := range map[int]image.Image{
		0:        img,
		90 - 15:  imaging.Rotate(i90, -15, color.White),
		90:       i90,
		90 + 15:  imaging.Rotate(i90, 15, color.White),
		270 - 15: imaging.Rotate(i270, -15, color.White),
		270:      i270,
		270 + 15: imaging.Rotate(i270, 15, color.White),
	} {
		grp.Go(func() error {
			start := time.Now()
			results, err := grcode.GetDataFromImage(img)
			log.Printf("%d: %s (%+v)", a, time.Since(start), err)
			if err != nil {
				return err
			}
			if len(results) == 0 {
				log.Printf("No qrcode detected from file: %s", filePath)
			}
			resultMu.Lock()
			result[a] = results
			resultMu.Unlock()
			return nil
		})
	}

	err = grp.Wait()
	for _, a := range slices.Sorted(maps.Keys(result)) {
		fmt.Printf("%d: %q\n", a, result[a])
	}
	return err
}
