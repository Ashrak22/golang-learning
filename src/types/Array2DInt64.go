package types

import "bettererror"

//Array2DInt64 is a two dimensional array of int64, internally it is stored as
//one dimensional array of length xSize * ySize and proper element is always computed
//in corresponding Get or Set procedure
type Array2DInt64 struct {
	xSize, ySize int
	array        []int64
}

//NewArray2DInt64 return new object of type Array2DInt64 which has the size
//of x*y
func NewArray2DInt64(x, y int) (res *Array2DInt64) {
	res = &Array2DInt64{xSize: x, ySize: y}
	res.array = make([]int64, x*y)
	return res
}

//GetValue gets the value at coordinates x and y, the proper element is counted as:
//y*xSize + x
func (a *Array2DInt64) GetValue(x, y int) (int64, error) {
	if x >= a.xSize || x < 0 || y >= a.ySize || y < 0 {
		return -1, bettererror.NewBetterError(myFacility, 0x0001, myErrors[0x0001])
	}
	return a.array[y*a.xSize+x], nil
}

//SetValue sets the value at coordinates x and y, the proper element is counted as:
//y*xSize + x
func (a *Array2DInt64) SetValue(x, y int, val int64) error {
	if x >= a.xSize || x < 0 || y >= a.ySize || y < 0 {
		return bettererror.NewBetterError(myFacility, 0x0001, myErrors[0x0001])
	}
	a.array[y*a.xSize+x] = val
	return nil
}

//Size returns the x and y sizes of the array
func (a *Array2DInt64) Size() (x int, y int) {
	x = a.xSize
	y = a.ySize
	return x, y
}
