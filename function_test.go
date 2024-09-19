package augment

import (
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"os"
	"testing"
)

func TestAugment(t *testing.T) {
	_color := color.RGBA{R: 255}
	thickness := 2
	filename := "test.png"
	cells := []Cell{
		{LabelId: 3, Rectangle: image.Rect(163, 317, 397, 523)},
		{LabelId: 5, Rectangle: image.Rect(65, 0, 689, 269)},
		{LabelId: 7, Rectangle: image.Rect(116, 3, 429, 268)},
		{LabelId: 4, Rectangle: image.Rect(425, 180, 720, 409)},
	}
	img := gocv.IMRead(filename, gocv.IMReadColor)
	defer img.Close()
	augment := NewDataAugment()
	//augment := &DataAugment{
	//	RotationRate:     1,
	//	MaxRotationAngle: 5,
	//	CropRate:         1,
	//	ShiftRate:        1,
	//	ChangeLightRate:  1,
	//	AddNoiseRate:     1,
	//	FlipRate:         1,
	//	CutoutRate:       1,
	//	CutOutLength:     100,
	//	CutOutHoles:      1,
	//	CutOutThreshold:  0.5,
	//	IsAddNoise:       false, // 噪声 success
	//	IsChangeLight:    false, // 光照 success
	//	IsCutout:         false, // 抠图 success
	//	IsRotateImgBbox:  true,  // 旋转 success
	//	IsCropImgBBoxes:  false, // 裁剪 success
	//	IsShiftPicBBoxes: false, // 平移 success
	//	IsFlipPicBBoxes:  false, // 翻转 success
	//}
	img, cells = augment.Augment(img, cells)

	//for i := 0; i < 2; i++ {
	//	mat := img.Clone()
	//	newImage, newCells := augment.Augment(mat, cells)
	//	for _, c := range newCells {
	//		gocv.Rectangle(&newImage, c.Rectangle, _color, thickness)
	//	}
	//	gocv.IMWrite(fmt.Sprintf("%d.png", i), newImage)
	//	newImage.Close()
	//}

	t.Log(cells)
	for _, cell := range cells {
		gocv.Rectangle(&img, cell.Rectangle, _color, thickness)
	}
	showWin(img)
}

func showWin(mat gocv.Mat) {
	window := gocv.NewWindow("test")
	window.IMShow(mat)
	window.WaitKey(0)
	os.Exit(1)
}
