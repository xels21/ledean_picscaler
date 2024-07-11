package ledean_picscaler

import (
	"fmt"
	"image"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/sunshineplan/imgconv"
)

type PicScaler struct {
	inDir      string
	outDir     string
	pixelCount int
	picNames   []string
	isReadDone bool
}

func NewPicScaler(inDir string, outDirName string, pixelCount int) *PicScaler {
	self := PicScaler{
		inDir:      inDir,
		outDir:     filepath.Join(inDir, outDirName),
		pixelCount: pixelCount,
	}

	return &self
}

func (self *PicScaler) Scale() {
	err := os.Mkdir(self.outDir, os.ModeDir)
	if err != nil {
		log.Fatal(err)
	}
	self.ScaleToPixel()
}

func (self *PicScaler) readInDir() {
	entries, err := os.ReadDir(self.inDir)
	if err != nil {
		log.Fatal(err)
	}
	self.picNames = make([]string, 0, len(entries))
	for _, e := range entries {
		switch filepath.Ext(e.Name()) {
		case ".png", ".jpeg", ".jpg":
			self.picNames = append(self.picNames, e.Name())
			log.Print("Found image: " + e.Name())
		}
	}

}

func (self *PicScaler) ScaleToPixel() {
	if !self.isReadDone {
		self.readInDir()
	}
	for _, picName := range self.picNames {
		self.ScaleSingleToPixel(picName)
	}
}

func RemoveFileExtension(filename string) string {
	return strings.TrimSuffix(filename, filepath.Ext(filename))
}

func (self *PicScaler) ScaleSingleToPixel(picName string) {
	src, err := imgconv.Open(filepath.Join(self.inDir, picName))
	// defer src.Close()
	if err != nil {
		log.Fatal(err)
	}
	resized := imgconv.Resize(src, &imgconv.ResizeOption{Height: self.pixelCount})

	// Write the resulting image as TIFF.
	outPath := RemoveFileExtension(picName) + ".tiff"
	os.MkdirAll(self.outDir, os.ModePerm)
	self.ConvertToGo(resized.(*image.NRGBA), picName)
	output, err := os.Create(filepath.Join(self.outDir, outPath))
	defer output.Close()
	if err != nil {
		log.Fatal(err)
	}
	err = imgconv.Write(output, resized, &imgconv.FormatOption{Format: imgconv.TIFF})
	if err != nil {
		log.Fatalf("failed to write image: %v", err)
	}
}

func (self *PicScaler) ConvertToGo(resized *image.NRGBA, picName string) {
	/*
		format is:
		data.Pix -> Array R, G, B, A
		data.Rect.Max.X -> col
		data.Rect.Max.Y -> row
	*/

	output, err := os.Create(filepath.Join(self.outDir, RemoveFileExtension(picName)+".go"))
	defer output.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprint(output, `
	package picture

	import "ledean/color"
	
	var `, picName, " = PicRGB{col: []PicCol{\n")
	// var pic_1 = PicRGB{col: []PicCol{PicCol{row: []color.RGB{
	// color.RGB{R: 1, G: 1, B: 1},
	// color.RGB{R: 1, G: 1, B: 1},
	// color.RGB{R: 1, G: 1, B: 1}},
	// }, PicCol{row: []color.RGB{
	// color.RGB{R: 1, G: 1, B: 1},
	// color.RGB{R: 1, G: 1, B: 1},
	// color.RGB{R: 1, G: 1, B: 1}},
	// }}}
	for y := range resized.Rect.Max.Y {
		// output.Write("{\n")
		fmt.Fprint(output, "PicCol{row:")
		for x := range resized.Rect.Max.X {
			off := y*resized.Rect.Max.X + x*4
			fmt.Fprintf(output, " []color.RGB{%d,%d,%d},", resized.Pix[off], resized.Pix[off+1], resized.Pix[off+2])
			// output.Write("{}")
		}
		// output.Write("\n},")
		fmt.Fprint(output, "},\n")

	}
	fmt.Fprint(output, "}}")
}
