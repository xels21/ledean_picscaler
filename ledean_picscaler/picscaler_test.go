package ledean_picscaler

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const IN_PATH = "./testdata"
const OUT_NAME = "gen"

func TestPicScale(t *testing.T) {
	out := filepath.Join(IN_PATH, OUT_NAME)
	// pixelCount := 4
	pixelCount := 50

	os.RemoveAll(out)
	// os.MkdirAll("/tmp/",FileMode)
	picScaler := NewPicScaler(IN_PATH, OUT_NAME, pixelCount)
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
	picScaler := NewPicScaler(IN_PATH, OUT_NAME, pixelCount)
	picScaler.ScaleSingleToPixel("test_3x2.png")
	entries, err := os.ReadDir(out)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, 2, len(entries))
	// assert.InDelta(t, expectedRgb.B, rgb.B, 1)
}
