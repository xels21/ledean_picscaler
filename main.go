package main

import (
	"ledean_pic_scaler/ledean_picscaler"
	"ledean_pic_scaler/parameter"
)

func main() {
	parm := parameter.GetParameter()
	picScaler := ledean_picscaler.NewPicScaler(parm.InDir, parm.OutDirName, parm.PixelCount)
	picScaler.Scale()
}
