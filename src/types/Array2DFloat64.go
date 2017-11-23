package types

import (
	"bettererror"
)

func init() {
	bettererror.RegisterFacility(0xF001, "Types")
}

type Array2DFloat64 struct {
	xSize, ySize int
	array        []float64
}

func NewArray2DFloat64(x, y int) (res *Array2DFloat64) {
	res = &Array2DFloat64{xSize: x, ySize: y}
	res.array = make([]float64, x*y)
	return res
}

func (a *Array2DFloat64) GetValue(x, y int) (float64, error) {
	if x >= a.xSize || x < 0 || y >= a.ySize || y < 0 {
		return -1, bettererror.NewBetterError(0xF001, 0x0001, "Out of bounds")
	}
	return a.array[y*a.xSize+x], nil
}

func (a *Array2DFloat64) SetValue(x, y int, val float64) error {
	if x >= a.xSize || x < 0 || y >= a.ySize || y < 0 {
		return bettererror.NewBetterError(0xF001, 0x0001, "Out of bounds")
	}
	a.array[y*a.xSize+x] = val
	return nil
}

func (a *Array2DFloat64) Size() (x int, y int) {
	x = a.xSize
	y = a.ySize
	return x, y
}
