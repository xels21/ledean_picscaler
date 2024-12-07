package ledean_picscaler

import (
	"fmt"
	"image"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/sunshineplan/imgconv"
)

type PicScaler struct {
	inDir      string
	pixelCount int
	picNames   []string
	isReadDone bool
	asBytes    bool
	name       string
	outDir     string
}

func NewPicScaler(inDir string, name string, pixelCount int, asBytes bool) *PicScaler {
	self := PicScaler{
		inDir:      inDir,
		pixelCount: pixelCount,
		asBytes:    asBytes,
		name:       name,
		outDir:     filepath.Join(inDir, "gen_"+name),
	}

	return &self
}

func (self *PicScaler) CreateController() {
	output, err := os.Create(filepath.Join(self.outDir, "pics_"+self.name+".go"))
	defer output.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprint(output, "package "+self.name+"\n\n")
	if !self.asBytes {
		fmt.Fprint(output, "import \"image\"\n\n")
	}
	fmt.Fprint(output, "var PixelCount = "+strconv.Itoa(self.pixelCount)+"\n\n")

	if self.asBytes {
		fmt.Fprint(output, "var Pics = [][]string{")
	} else {
		fmt.Fprint(output, "var Pics = []*image.NRGBA{")
	}
	for _, picName := range self.picNames {
		picNameWoExtension := RemoveFileExtension(picName)
		if self.asBytes {
			fmt.Fprint(output, "\n\t"+self.name+"_"+picNameWoExtension+",")
		} else {
			fmt.Fprint(output, "\n\t"+"&"+self.name+"_"+picNameWoExtension+",")
		}
	}
	fmt.Fprint(output, "\n}\n")
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

// func NRGBAToStringArray(self *image.NRGBA) string {
// 	return fmt.Sprintf(`image.NRGBA{
// 	Pix:    []uint8{%s},
// 	Stride: %d,
// 	Rect:   image.Rect(%d, %d, %d, %d),
// }`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(self.Pix)), ", "), "[]"), self.Stride, self.Rect.Min.X, self.Rect.Min.Y, self.Rect.Max.X, self.Rect.Max.Y)
// }

func RgbaToRgbString(rgba []uint8, pixelPerRow int) string {
	pixelCount := len(rgba) / 4
	rowCount := pixelCount / pixelPerRow
	pixelAsString := "[]string{"
	for r := 0; r < rowCount; r++ {
		pixelAsString += "\n\t\""
		for i := 0; i < pixelPerRow; i++ {
			pixelAsString += fmt.Sprintf("\\x%02x\\x%02x\\x%02x", rgba[r*pixelPerRow*4+i*4+0], rgba[r*pixelPerRow*4+i*4+1], rgba[r*pixelPerRow*4+i*4+2])
			// pixelAsString += string(rgba[r*pixelPerRow*4+i*4 : r*pixelPerRow*4+i*4+3])
		}
		pixelAsString += "\","
	}
	return pixelAsString + "\n}\n"
}

func NRGBAToStringArray(self *image.NRGBA) string {
	return RgbaToRgbString(self.Pix, self.Rect.Max.X-self.Rect.Min.X)
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
	output, err := os.Create(filepath.Join(self.outDir, strings.ToLower(self.name)+"_"+picNameWoExtension+".go"))
	defer output.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprint(output, "package "+self.name+"\n\n")
	if !self.asBytes {
		fmt.Fprint(output, `
		import (
			"image"
		)
	`)
	}
	fmt.Fprint(output, "var ", self.name, "_", picNameWoExtension, " = ")
	if self.asBytes {
		fmt.Fprint(output, NRGBAToStringArray(resized))
	} else {
		fmt.Fprint(output, NRGBAToGo(resized))
	}
}
