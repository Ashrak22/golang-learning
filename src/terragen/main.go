package main

import (
	"args"
	"bettererror"
	"image"
	"image/color"
	"image/jpeg"
	"math"
	"os"
)

func init() {
	bettererror.RegisterFacility(myFacility, "Terrain Generator")
}

func main() {
	a := args.NewArg()
	a.EvalArgs(os.Args)
	t := NewTerrain(10, 0.5)
	t.SetStart()
	t.Divide(t.max)
	//t.print()
	for y := 0; y < t.size; y++ {
		for x := 0; x < t.size; x++ {
			t.setValue(x, y, float32(math.Floor(float64((t.getValue(x, y)/float32(t.max))*256))))
		}
	}
	//t.print()
	img := image.NewGray(image.Rect(0, 0, t.size, t.size))
	for y := 0; y < t.size; y++ {
		for x := 0; x < t.size; x++ {
			var gray color.Gray
			gray.Y = uint8(t.getValue(x, y))
			img.SetGray(x, y, gray)
		}
	}
	//fmt.Println(img.Pix)
	file, _ := os.Create("test.jpg")
	jpeg.Encode(file, img, &jpeg.Options{Quality: 100})
	//png.Encode(file, img)
}
