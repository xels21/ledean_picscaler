package parameter

import (
	"flag"
)

type Parameter struct {
	InDir      string `json:"inDir"`
	PixelCount int    `json:"pixelCount"`
	OutDirName string `json:"outDirName"`
	AsByte     bool   `json:"asByte"`
}

func GetParameter() *Parameter {
	var parm Parameter
	flag.StringVar(&parm.InDir, "in", ".", "Path to directory of to cenverting pictures")
	flag.IntVar(&parm.PixelCount, "pixelCount", 50, "Amount of pixel in one column")
	flag.StringVar(&parm.OutDirName, "outName", "gen_poi", "Name of output directory")
	flag.BoolVar(&parm.AsByte, "asByte", false, "Defines whether output should be generated as byte array (string)")
	flag.Parse()
	return &parm
}
