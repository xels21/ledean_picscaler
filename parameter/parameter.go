package parameter

import (
	"flag"
)

type Parameter struct {
	InDir      string `json:"inDir"`
	PixelCount int    `json:"pixelCount"`
	OutDirName string `json:"outDirName"`
}

func GetParameter() *Parameter {
	var parm Parameter
	flag.StringVar(&parm.InDir, "in", ".", "Path to directory of to cenverting pictures")
	flag.IntVar(&parm.PixelCount, "pixelCount", 50, "Amount of pixel in one column")
	flag.StringVar(&parm.OutDirName, "outName", "gen", "Name of output directory")
	flag.Parse()
	return &parm
}
