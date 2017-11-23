package types

import "bettererror"

type Array2DInt64 struct {
	xSize, ySize int
	array        []int64
}

func NewArray2DInt64(x, y int) (res *Array2DInt64) {
	res = &Array2DInt64{xSize: x, ySize: y}
	res.array = make([]int64, x*y)
	return res
}

func (a *Array2DInt64) GetValue(x, y int) (int64, error) {
	if x >= a.xSize || x < 0 || y >= a.ySize || y < 0 {
		return -1, bettererror.NewBetterError(0xF001, 0x0001, "Out of bounds")
	}
	return a.array[y*a.xSize+x], nil
}

func (a *Array2DInt64) SetValue(x, y int, val int64) error {
	if x >= a.xSize || x < 0 || y >= a.ySize || y < 0 {
		return bettererror.NewBetterError(0xF001, 0x0001, "Out of bounds")
	}
	a.array[y*a.xSize+x] = val
	return nil
}

func (a *Array2DInt64) Size() (x int, y int) {
	x = a.xSize
	y = a.ySize
	return x, y
}
