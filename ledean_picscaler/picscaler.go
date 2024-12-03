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
	picPrefix  string
}

func NewPicScaler(inDir string, outDirName string, pixelCount int) *PicScaler {
	self := PicScaler{
		inDir:      inDir,
		outDir:     filepath.Join(inDir, outDirName),
		pixelCount: pixelCount,
		picPrefix:  "GetPoiPic_",
	}

	return &self
}

func (self *PicScaler) CreateController() {
	output, err := os.Create(filepath.Join(self.outDir, "poipics.go"))
	defer output.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(output, `package poi

import "image"

type GetPoiPic func() image.NRGBA

var PixelCount = %d

var PoiPicsCount = %d

var PoiPicsGetter = []GetPoiPic{`, self.pixelCount, len(self.picNames))
	for _, picName := range self.picNames {
		picNameWoExtension := RemoveFileExtension(picName)
		fmt.Fprint(output, self.picPrefix+picNameWoExtension+", ")
	}
	fmt.Fprint(output, "}")
}

func (self *PicScaler) Scale() {
	os.RemoveAll(self.outDir)
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
			if strings.HasPrefix(e.Name(), "_") {
				continue
			}
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
	resized := imgconv.Resize(src, &imgconv.ResizeOption{Width: self.pixelCount})

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

func NRGBAToGo(self *image.NRGBA) string {
	return fmt.Sprintf(`image.NRGBA{
	Pix:    []uint8{%s},
	Stride: %d,
	Rect:   image.Rect(%d, %d, %d, %d),
}`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(self.Pix)), ", "), "[]"), self.Stride, self.Rect.Min.X, self.Rect.Min.Y, self.Rect.Max.X, self.Rect.Max.Y)
}

func NRGBAToString(self *image.NRGBA) string {
	ret := "PicRGB{col: []PicCol{\n"
	// var pic_1 = PicRGB{col: []PicCol{PicCol{row: []color.RGB{
	// color.RGB{R: 1, G: 1, B: 1},
	// color.RGB{R: 1, G: 1, B: 1},
	// color.RGB{R: 1, G: 1, B: 1}},
	// }, PicCol{row: []color.RGB{
	// color.RGB{R: 1, G: 1, B: 1},
	// color.RGB{R: 1, G: 1, B: 1},
	// color.RGB{R: 1, G: 1, B: 1}},
	// }}}`+
	for y := range self.Rect.Max.Y {
		// output.Write("{\n")
		// fmt.Fprint(output, "PicCol{row:")
		ret += "	PicCol{row:"
		for x := range self.Rect.Max.X {
			off := y*self.Rect.Max.X*4 + x*4
			// fmt.Fprintf(output, " []color.RGB{%d,%d,%d},", self.Pix[off], self.Pix[off+1], self.Pix[off+2])
			ret += fmt.Sprintf(" []color.RGB{%d,%d,%d},", self.Pix[off], self.Pix[off+1], self.Pix[off+2])
			// output.Write("{}")
		}
		// output.Write("\n},")
		// fmt.Fprint(output, "},\n")
		ret += "},\n"

	}
	// fmt.Fprint(output, "}}")`
	ret += "	}\n}"
	return ret
}

func (self *PicScaler) ConvertToGo(resized *image.NRGBA, picName string) {
	/*
		format is:
		data.Pix -> Array R, G, B, A
		data.Rect.Max.X -> col
		data.Rect.Max.Y -> row
	*/

	picNameWoExtension := RemoveFileExtension(picName)
	output, err := os.Create(filepath.Join(self.outDir, strings.ToLower(self.picPrefix)+picNameWoExtension+".go"))
	defer output.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprint(output, `package poi

import (
	"image"
)

func `, self.picPrefix, picNameWoExtension, `() image.NRGBA {
	return `, NRGBAToGo(resized), `
}`)

	// var `, self.picPrefix, picNameWoExtension, " = ", NRGBAToGo(resized))
	// var `, picName, " := ", NRGBAToGo NRGBAToString(resized))
}
