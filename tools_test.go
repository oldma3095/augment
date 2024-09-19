package augment

import (
	"image"
	"math/rand"
	"testing"
	"time"
)

func TestUniformInt(t *testing.T) {
	maxRotationAngle := 5
	n1, n2 := -maxRotationAngle, maxRotationAngle
	rand.NewSource(time.Now().UnixNano())
	num := n1 + rand.Intn(n2-n1+1)
	t.Log(num)
}

func TestIOU(t *testing.T) {
	t.Log(iou(
		image.Rect(10, 10, 50, 50),
		image.Rect(30, 30, 70, 70),
	))

	t.Log(iou(
		image.Rect(30, 30, 70, 70),
		image.Rect(10, 10, 50, 50),
	))
}
