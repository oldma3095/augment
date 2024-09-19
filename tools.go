package augment

import (
	"image"
	"math/rand"
	"time"
)

// 计算两个矩形框的 IOU
func iou(rect1, rect2 image.Rectangle) float64 {
	// 找到相交矩形
	intersection := rect1.Intersect(rect2)
	if intersection.Empty() {
		return 0.0
	}
	rect1Area := rect1.Dx() * rect1.Dy()
	rect2Area := rect2.Dx() * rect2.Dy()

	// 计算交集面积
	intersectionArea := intersection.Dx() * intersection.Dy()
	// 计算并集面积
	unionArea := rect1Area + rect2Area - intersectionArea

	iou := float64(intersectionArea) / float64(unionArea)
	return iou
}

func randomUniformFloat64(n1, n2 float64) float64 {
	rand.NewSource(time.Now().UnixNano())
	return n1 + rand.Float64()*(n2-n1)
}

func randomUniformInt(n1, n2 int) int {
	rand.NewSource(time.Now().UnixNano())
	return n1 + rand.Intn(n2-n1+1)
}
