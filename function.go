package augment

import (
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"math"
	"math/rand"
	"time"
)

// 加噪音 高斯模糊
func (a *DataAugment) addNoise(img gocv.Mat) gocv.Mat {
	defer img.Close()
	result := gocv.NewMat()
	gocv.GaussianBlur(img, &result, image.Pt(11, 11), 0, 0, gocv.BorderDefault)
	return result
}

// 调整亮度
func (a *DataAugment) changeLight(img gocv.Mat) gocv.Mat {
	defer img.Close()
	// 生成一个在 [0.35, 1] 范围内的随机浮点数
	alpha := randomUniformFloat64(0.35, 1)
	blank := gocv.NewMatWithSize(img.Rows(), img.Cols(), img.Type())
	defer blank.Close() // Ensure blank image is released
	result := gocv.NewMat()
	gocv.AddWeighted(img, alpha, blank, 1-alpha, 0, &result)
	return result
}

// 抠图
func (a *DataAugment) cutout(img gocv.Mat, cells []Cell, length int, nHoles int, threshold float64) gocv.Mat {
	defer img.Close()
	if length == 0 {
		length = 100
	}
	if nHoles == 0 {
		nHoles = 1
	}
	if threshold == 0 {
		threshold = 0.5
	}
	h := img.Rows()
	w := img.Cols()

	// 初始化遮罩
	mask := gocv.NewMatWithSize(h, w, gocv.MatTypeCV8UC3)
	defer mask.Close()
	mask.SetTo(gocv.NewScalar(255, 255, 255, 0)) // white mask

	rand.NewSource(time.Now().UnixNano())

	for n := 0; n < nHoles; n++ {
		overlap := true // 看切割的区域是否与box重叠太多
		var (
			x1, x2, y1, y2 int
		)
		for overlap {
			y := rand.Intn(h)
			x := rand.Intn(w)

			// 定义要遮罩的区域
			y1 = max(0, y-length/2)
			y2 = min(h, y+length/2)
			x1 = max(0, x-length/2)
			x2 = min(w, x+length/2)

			overlap = false

			maskRect := image.Rect(x1, y1, x2, y2)
			for _, cell := range cells {
				if iou(maskRect, cell.Rectangle) > threshold {
					overlap = true
					break
				}
			}
		}

		// 蒙版 黑色矩形 遮挡区域
		gocv.Rectangle(&mask, image.Rect(x1, y1, x2, y2), color.RGBA{}, -1)
	}

	// 应用蒙版到图像
	result := gocv.NewMat()
	gocv.BitwiseAnd(img, mask, &result)
	return result
}

// 旋转
func (a *DataAugment) rotateImageBBbox(img gocv.Mat, cells []Cell, angle int, scale float64) (gocv.Mat, []Cell) {
	defer img.Close()
	if angle == 0 {
		angle = 5
	}
	if scale == 0 {
		scale = 1.0
	}
	h := img.Rows()
	w := img.Cols()
	// 旋转角度
	rangle := float64(angle) * math.Pi / 180.0

	// 计算新图像的宽度和高度
	nw := int(math.Abs(math.Sin(rangle)*float64(h)) + math.Abs(math.Cos(rangle)*float64(w))*scale)
	nh := int(math.Abs(math.Cos(rangle)*float64(h)) + math.Abs(math.Sin(rangle)*float64(w))*scale)

	// 旋转矩阵
	rotMat := gocv.GetRotationMatrix2D(image.Point{X: nw / 2, Y: nh / 2}, float64(angle), scale)
	defer rotMat.Close()

	// 结合旋转计算从旧中心到新中心的移动
	rotMove := []float64{float64(nw-w) * 0.5, float64(nh-h) * 0.5}

	// 更新旋转矩阵的平移部分
	rotMat.SetDoubleAt(0, 2, rotMat.GetDoubleAt(0, 2)+rotMove[0])
	rotMat.SetDoubleAt(1, 2, rotMat.GetDoubleAt(1, 2)+rotMove[1])

	// 应用仿射扭曲来旋转图像
	rotatedImg := gocv.NewMat()
	gocv.WarpAffine(img, &rotatedImg, rotMat, image.Point{X: nw, Y: nh})

	// 调整边界框
	rotBBoxes := make([]Cell, 0, len(cells))
	for _, cell := range cells {
		xmin, ymin, xmax, ymax := cell.Rectangle.Min.X, cell.Rectangle.Min.Y, cell.Rectangle.Max.X, cell.Rectangle.Max.Y

		// 计算边界框的四个角
		points := []gocv.Point2f{
			{X: float32((xmin + xmax) / 2), Y: float32(ymin)},
			{X: float32(xmax), Y: float32((ymin + ymax) / 2)},
			{X: float32((xmin + xmax) / 2), Y: float32(ymax)},
			{X: float32(xmin), Y: float32((ymin + ymax) / 2)},
		}

		// 将旋转矩阵应用于四个点
		rotatedPoints := make([]gocv.Point2f, 4)
		for i, pt := range points {
			x := float64(pt.X)
			y := float64(pt.Y)
			newX := rotMat.GetDoubleAt(0, 0)*x + rotMat.GetDoubleAt(0, 1)*y + rotMat.GetDoubleAt(0, 2)
			newY := rotMat.GetDoubleAt(1, 0)*x + rotMat.GetDoubleAt(1, 1)*y + rotMat.GetDoubleAt(1, 2)
			rotatedPoints[i] = gocv.Point2f{X: float32(newX), Y: float32(newY)}
		}

		// 从旋转点获取新的边界框
		rxMin, ryMin := int(rotatedPoints[0].X), int(rotatedPoints[0].Y)
		rxMax, ryMax := rxMin, ryMin
		for _, pt := range rotatedPoints {
			rxMin = int(math.Min(float64(rxMin), float64(pt.X)))
			ryMin = int(math.Min(float64(ryMin), float64(pt.Y)))
			rxMax = int(math.Max(float64(rxMax), float64(pt.X)))
			ryMax = int(math.Max(float64(ryMax), float64(pt.Y)))
		}

		// 将新的边界框附加到列表中
		rotBBoxes = append(rotBBoxes, Cell{LabelId: cell.LabelId, Rectangle: image.Rect(rxMin, ryMin, rxMax, ryMax)})
	}

	return rotatedImg, rotBBoxes
}

// 裁剪
func (a *DataAugment) cropImgAndBBoxes(img gocv.Mat, cells []Cell) (gocv.Mat, []Cell) {
	defer img.Close()
	h := img.Rows()
	w := img.Cols()

	// 将裁剪坐标初始化为图像的边界
	xMin, xMax := w, 0
	yMin, yMax := h, 0

	// 找到边界框的最小和最大坐标
	for _, cell := range cells {
		xMin = min(xMin, cell.Rectangle.Min.X)
		yMin = min(yMin, cell.Rectangle.Min.Y)
		xMax = max(xMax, cell.Rectangle.Max.X)
		yMax = max(yMax, cell.Rectangle.Max.Y)
	}

	// 边界框到图像边缘的距离
	dToLeft := xMin
	dToRight := w - xMax
	dToTop := yMin
	dToBottom := h - yMax

	// 随机扩大裁剪区域
	rand.NewSource(time.Now().UnixNano())
	cropXMin := int(float64(xMin) - rand.Float64()*float64(dToLeft))
	cropYMin := int(float64(yMin) - rand.Float64()*float64(dToTop))
	cropXMax := int(float64(xMax) + rand.Float64()*float64(dToRight))
	cropYMax := int(float64(yMax) + rand.Float64()*float64(dToBottom))

	// 确保裁剪坐标不超出范围
	cropXMin = max(0, cropXMin)
	cropYMin = max(0, cropYMin)
	cropXMax = min(w, cropXMax)
	cropYMax = min(h, cropYMax)

	// 裁剪图像
	croppedImg := img.Region(image.Rect(cropXMin, cropYMin, cropXMax, cropYMax))

	// 调整边界框
	croppedBBoxes := make([]Cell, len(cells))
	for i, cell := range cells {
		newXMin := cell.Rectangle.Min.X - cropXMin
		newYMin := cell.Rectangle.Min.Y - cropYMin
		newXMax := cell.Rectangle.Max.X - cropXMin
		newYMax := cell.Rectangle.Max.Y - cropYMin
		croppedBBoxes[i] = Cell{LabelId: cell.LabelId, Rectangle: image.Rect(newXMin, newYMin, newXMax, newYMax)}
	}

	return croppedImg, croppedBBoxes
}

// 平移
func (a *DataAugment) shiftImgAndBBoxes(img gocv.Mat, cells []Cell) (gocv.Mat, []Cell) {
	defer img.Close()
	h := img.Rows()
	w := img.Cols()

	// 初始化变量以查找包含所有框的边界框
	xMin, xMax := w, 0
	yMin, yMax := h, 0

	// 找到边界框的最小和最大坐标
	for _, cell := range cells {
		xMin = min(xMin, cell.Rectangle.Min.X)
		yMin = min(yMin, cell.Rectangle.Min.Y)
		xMax = max(xMax, cell.Rectangle.Max.X)
		yMax = max(yMax, cell.Rectangle.Max.Y)
	}

	// 计算最大移动距离
	dToLeft := xMin
	dToRight := w - xMax
	dToTop := yMin
	dToBottom := h - yMax

	// 随机选择移位值
	rand.NewSource(time.Now().UnixNano())
	xShift := rand.Float64()*((float64(dToRight)-1)/3) - ((float64(dToLeft) - 1) / 3)
	yShift := rand.Float64()*((float64(dToBottom)-1)/3) - ((float64(dToTop) - 1) / 3)

	// 定义用于平移的仿射变换矩阵
	M := gocv.NewMatWithSize(2, 3, gocv.MatTypeCV32F)
	defer M.Close()
	M.SetFloatAt(0, 0, 1)
	M.SetFloatAt(0, 1, 0)
	M.SetFloatAt(0, 2, float32(xShift))
	M.SetFloatAt(1, 0, 0)
	M.SetFloatAt(1, 1, 1)
	M.SetFloatAt(1, 2, float32(yShift))

	// 应用仿射变换来移动图像
	shiftedImg := gocv.NewMat()
	gocv.WarpAffine(img, &shiftedImg, M, image.Pt(w, h))

	// 调整边界框
	shiftedBBoxes := make([]Cell, len(cells))
	for i, cell := range cells {
		newXMin := int(float64(cell.Rectangle.Min.X) + xShift)
		newYMin := int(float64(cell.Rectangle.Min.Y) + yShift)
		newXMax := int(float64(cell.Rectangle.Max.X) + xShift)
		newYMax := int(float64(cell.Rectangle.Max.Y) + yShift)
		shiftedBBoxes[i] = Cell{LabelId: cell.LabelId, Rectangle: image.Rect(newXMin, newYMin, newXMax, newYMax)}
	}

	return shiftedImg, shiftedBBoxes
}

// 镜像
func (a *DataAugment) flipImgAndBBoxes(img gocv.Mat, cells []Cell) (gocv.Mat, []Cell) {
	defer img.Close()
	h := img.Rows()
	w := img.Cols()

	rand.NewSource(time.Now().UnixNano())
	sed := rand.Float64()

	var flipMode int
	if sed < 0.33 {
		flipMode = 0 // 0.33的概率水平翻转
	} else if sed < 0.66 {
		flipMode = 1 // 0.33的概率垂直翻转
	} else {
		flipMode = -1 // 0.33是对角反转
	}

	flipImg := gocv.NewMat()
	gocv.Flip(img, &flipImg, flipMode)

	// 调整边界框
	flipBBoxes := make([]Cell, len(cells))
	for i, cell := range cells {
		xMin := cell.Rectangle.Min.X
		yMin := cell.Rectangle.Min.Y
		xMax := cell.Rectangle.Max.X
		yMax := cell.Rectangle.Max.Y

		if flipMode == 0 {
			flipBBoxes[i] = Cell{LabelId: cell.LabelId, Rectangle: image.Rect(xMin, h-yMax, xMax, h-yMin)}
		} else if flipMode == 1 {
			flipBBoxes[i] = Cell{LabelId: cell.LabelId, Rectangle: image.Rect(w-xMax, yMin, w-xMin, yMax)}
		} else {
			flipBBoxes[i] = Cell{LabelId: cell.LabelId, Rectangle: image.Rect(w-xMax, h-yMax, w-xMin, h-yMin)}
		}
	}

	return flipImg, flipBBoxes
}

func (a *DataAugment) Augment(img gocv.Mat, cells []Cell) (gocv.Mat, []Cell) {
	rand.NewSource(time.Now().UnixNano())

	if a.IsRotateImgBbox {
		if rand.Float64() < a.RotationRate { // 旋转
			angle := randomUniformInt(-a.MaxRotationAngle, a.MaxRotationAngle)
			scale := randomUniformFloat64(0.7, 0.8)
			img, cells = a.rotateImageBBbox(img, cells, angle, scale)
		}
	}

	if a.IsShiftPicBBoxes {
		if rand.Float64() < a.ShiftRate { // 平移
			img, cells = a.shiftImgAndBBoxes(img, cells)
		}
	}

	if a.IsChangeLight {
		if rand.Float64() < a.ChangeLightRate { // 改变亮度
			img = a.changeLight(img)
		}
	}

	if a.IsAddNoise {
		if rand.Float64() < a.AddNoiseRate { // 加噪声
			img = a.addNoise(img)
		}
	}

	if a.IsCutout {
		if rand.Float64() < a.CutoutRate { // cutout
			img = a.cutout(img, cells, a.CutOutLength, a.CutOutHoles, a.CutOutThreshold)
		}
	}

	if a.IsFlipPicBBoxes {
		if rand.Float64() < a.FlipRate { // 翻转
			img, cells = a.flipImgAndBBoxes(img, cells)
		}
	}

	if a.IsCropImgBBoxes {
		if rand.Float64() < a.CropRate { // 裁剪
			img, cells = a.cropImgAndBBoxes(img, cells)
		}
	}

	return img, cells
}
