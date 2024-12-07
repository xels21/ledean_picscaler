package main

import (
	"ledean_pic_scaler/ledean_picscaler"
	"ledean_pic_scaler/parameter"
)

func main() {
	parm := parameter.GetParameter()
	picScaler := ledean_picscaler.NewPicScaler(parm.InDir, parm.Name, parm.PixelCount, parm.AsByte)
	picScaler.Scale()
	picScaler.CreateController()
}
