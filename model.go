package augment

import "image"

type DataAugment struct {
	// 旋转相关参数
	RotationRate     float64 // 旋转的概率
	MaxRotationAngle int     // 最大旋转角度

	// 裁剪相关参数
	CropRate float64 // 裁剪的概率

	// 平移相关参数
	ShiftRate float64 // 平移的概率

	// 改变光照相关参数
	ChangeLightRate float64 // 改变光照的概率

	// 添加噪声相关参数
	AddNoiseRate float64 // 添加噪声的概率

	// 翻转相关参数
	FlipRate float64 // 翻转的概率

	// CutOut 相关参数
	CutoutRate      float64 // CutOut 的概率
	CutOutLength    int     // CutOut 的长度
	CutOutHoles     int     // CutOut 的孔洞数量
	CutOutThreshold float64 // CutOut 的阈值

	// 是否启用某种增强方式
	IsAddNoise       bool // 是否启用添加噪声
	IsChangeLight    bool // 是否启用改变光照
	IsCutout         bool // 是否启用 CutOut
	IsRotateImgBbox  bool // 是否启用旋转
	IsCropImgBBoxes  bool // 是否启用裁剪
	IsShiftPicBBoxes bool // 是否启用平移
	IsFlipPicBBoxes  bool // 是否启用翻转
}

type Cell struct {
	LabelId   int
	Rectangle image.Rectangle
}

func NewDataAugment() *DataAugment {
	return &DataAugment{
		RotationRate:     0.5,
		MaxRotationAngle: 5,
		CropRate:         0.5,
		ShiftRate:        0.5,
		ChangeLightRate:  0.5,
		AddNoiseRate:     0.5,
		FlipRate:         0.5,
		CutoutRate:       0.5,
		CutOutLength:     50,
		CutOutHoles:      1,
		CutOutThreshold:  0.5,
		IsAddNoise:       true,
		IsChangeLight:    true,
		IsCutout:         true,
		IsRotateImgBbox:  true,
		IsCropImgBBoxes:  true,
		IsShiftPicBBoxes: true,
		IsFlipPicBBoxes:  false,
	}
}

func (a *DataAugment) SetRotationRate(rotationRate float64) {
	a.RotationRate = rotationRate
}
func (a *DataAugment) SetMaxRotationAngle(maxRotationAngle int) {
	a.MaxRotationAngle = maxRotationAngle
}
func (a *DataAugment) SetCropRate(cropRate float64) {
	a.CropRate = cropRate
}
func (a *DataAugment) SetShiftRate(shiftRate float64) {
	a.ShiftRate = shiftRate
}
func (a *DataAugment) SetChangeLightRate(changeLightRate float64) {
	a.ChangeLightRate = changeLightRate
}
func (a *DataAugment) SetAddNoiseRate(addNoiseRate float64) {
	a.AddNoiseRate = addNoiseRate
}
func (a *DataAugment) SetFlipRate(flipRate float64) {
	a.FlipRate = flipRate
}
func (a *DataAugment) SetCutoutRate(cutoutRate float64) {
	a.CutoutRate = cutoutRate
}
func (a *DataAugment) SetCutOutLength(cutOutLength int) {
	a.CutOutLength = cutOutLength
}
func (a *DataAugment) SetCutOutHoles(cutOutHoles int) {
	a.CutOutHoles = cutOutHoles
}
func (a *DataAugment) SetCutOutThreshold(cutOutThreshold float64) {
	a.CutOutThreshold = cutOutThreshold
}
func (a *DataAugment) SetIsAddNoise(isAddNoise bool) {
	a.IsAddNoise = isAddNoise
}
func (a *DataAugment) SetIsChangeLight(isChangeLight bool) {
	a.IsChangeLight = isChangeLight
}
func (a *DataAugment) SetIsCutout(isCutout bool) {
	a.IsCutout = isCutout
}
func (a *DataAugment) SetIsRotateImgBbox(isRotateImgBbox bool) {
	a.IsRotateImgBbox = isRotateImgBbox
}
func (a *DataAugment) SetIsCropImgBBoxes(isCropImgBBoxes bool) {
	a.IsCropImgBBoxes = isCropImgBBoxes
}
func (a *DataAugment) SetIsShiftPicBBoxes(isShiftPicBBoxes bool) {
	a.IsShiftPicBBoxes = isShiftPicBBoxes
}
func (a *DataAugment) SetIsFlipPicBBoxes(isFlipPicBBoxes bool) {
	a.IsFlipPicBBoxes = isFlipPicBBoxes
}
