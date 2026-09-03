package main

import (
	"context"
	"flag"
	"fmt"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"iter"
	"log"
	"maps"
	"os"
	"os/signal"
	"slices"
	"strconv"
	"strings"
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
	flagXDensity := flag.Int("x-density", 1, "x-density")
	flagYDensity := flag.Int("y-density", 1, "y-density")
	flagSymbols := flag.String("symbol-type", "*", "symbologies, space/comma separated")
	flagRotate := flag.String("rotate", "0", "rotations to try, space/comma seprated")
	flagConcurrency := flag.Int("concurrency", 4, "concurrency")
	flag.Usage = func() {
		flag.PrintDefaults()
	}
	flag.Parse()
	if *flagConcurrency < 1 {
		*flagConcurrency = 1
	}
	//log.SetFlags(0)
	if flag.NArg() < 1 {
		return fmt.Errorf("Need specify the image file")
	}
	filePath := flag.Arg(0)

	img, err := grcode.OpenImage(filePath)
	if err != nil {
		return err
	}

	commaSpace := func(r rune) bool { return r == ' ' || r == '\t' || r == '\n' || r == ',' }
	splitCommaSpace := func(s string) iter.Seq[string] { return strings.FieldsFuncSeq(s, commaSpace) }

	var rotations []int
	for s := range splitCommaSpace(*flagRotate) {
		i, err := strconv.ParseInt(s, 10, 9)
		if err != nil {
			return fmt.Errorf("parse rotation %s: %w", s, err)
		}
		if i != 0 {
			rotations = append(rotations, int(i))
		}
	}
	slices.Sort(rotations)
	slices.Compact(rotations)

	options := []grcode.Option{
		{Value: 0},
		{Config: grcode.XDensity, Value: *flagXDensity},
		{Config: grcode.YDensity, Value: *flagYDensity},
	}
	if *flagSymbols == "*" || *flagSymbols == "" {
		options[0].Value = 1
	} else {
		for s := range splitCommaSpace(*flagSymbols) {
			sym, err := grcode.ParseSymbolType(s)
			if err != nil {
				return fmt.Errorf("unknown symbology %q: %s", s, err)
			}
			options = append(options, grcode.Option{
				Symbology: sym, Config: 0, Value: 1})
		}
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	result := make(map[int][]string, 8)
	var resultMu sync.Mutex
	grp, _ := errgroup.WithContext(ctx)
	grp.SetLimit(*flagConcurrency)

	for _, r := range append(rotations, 0) {
		img := img
		if r != 0 {
			if r == 90 {
				img = imaging.Rotate90(img)
			} else if r == 270 {
				img = imaging.Rotate270(img)
			} else {
				img = imaging.Rotate(img, float64(r), color.White)
			}
		}
		grp.Go(func() error {
			start := time.Now()
			results, err := grcode.GetDataFromImage(img, options...)
			log.Printf("%d: %s (%+v)", r, time.Since(start), err)
			if err != nil {
				return err
			}
			if len(results) == 0 {
				log.Printf("No qrcode detected from file: %s", filePath)
			}
			resultMu.Lock()
			result[r] = results
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
