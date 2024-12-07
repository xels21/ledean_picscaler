package ledean_picscaler

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const IN_PATH = "./testdata"
const NAME = "gen"
const OUT_NAME = "gen_" + NAME

func TestPicScale(t *testing.T) {
	out := filepath.Join(IN_PATH, OUT_NAME)
	// pixelCount := 4
	pixelCount := 50

	os.RemoveAll(out)
	// os.MkdirAll("/tmp/",FileMode)
	picScaler := NewPicScaler(IN_PATH, NAME, pixelCount, false)
	picScaler.Scale()
	picScaler.CreateController()
	entries, err := os.ReadDir(out)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, 13, len(entries))
	// assert.InDelta(t, expectedRgb.B, rgb.B, 1)
}

func TestPicScaleSingle(t *testing.T) {
	out := filepath.Join(IN_PATH, OUT_NAME)
	pixelCount := 2
	// pixelCount := 32

	os.RemoveAll(out)
	// os.MkdirAll("/tmp/",FileMode)
	picScaler := NewPicScaler(IN_PATH, NAME, pixelCount, false)
	picScaler.ScaleSingleToPixel("test_3x2.png")
	entries, err := os.ReadDir(out)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, 2, len(entries))
	// assert.InDelta(t, expectedRgb.B, rgb.B, 1)
}

func TestRgbaToRgbString(t *testing.T) {
	input := []uint8{0x10, 0x11, 0x12, 0xFF, 0x20, 0x21, 0x22, 0xFF}

	expected := "[]string{\"\\x10\\x11\\x12\\x20\\x21\\x22\",}"
	actual := RgbaToRgbString(input, 2)
	assert.Equal(t, expected, actual)

	expected = "[]string{\"\\x10\\x11\\x12\",\"\\x20\\x21\\x22\",}"
	actual = RgbaToRgbString(input, 1)
	assert.Equal(t, expected, actual)
}
