package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

type Terrain struct {
	size      int
	max       int
	rand      *rand.Rand
	heightMap []float32
	roughness float32
}

func NewTerrain(detail int, roughness float32) *Terrain {
	res := new(Terrain)
	res.size = int(math.Pow(2, float64(detail))) + 1
	res.max = res.size - 1
	res.heightMap = make([]float32, res.size*res.size)
	res.roughness = roughness
	res.rand = rand.New(rand.NewSource(time.Now().Unix()))
	return res
}

func (t *Terrain) setValue(x int, y int, val float32) {
	if x > t.max || y > t.max || x < 0 || y < 0 {
		return
	}
	t.heightMap[x+y*t.size] = val
}

func (t *Terrain) getValue(x int, y int) float32 {
	if x > t.max || y > t.max || x < 0 || y < 0 {
		return -1
	}
	return t.heightMap[x+y*t.size]
}

func (t *Terrain) SetStart() {
	t.setValue(0, 0, float32(t.max))
	t.setValue(0, t.max, float32(t.max/2))
	t.setValue(t.max, 0, float32(t.max/2))
	t.setValue(t.max, t.max, float32(0))
}

func (t *Terrain) Divide(size int) {
	half := size / 2

	scale := t.roughness * float32(size)
	for y := half; y < t.max; y += half {
		for x := half; x < t.max; x += half {
			t.square(x, y, half, t.rand.Float32()*scale*2-scale)
		}
	}
	//t.print()
	for y := 0; y <= t.max; y += half {
		for x := (y + half) % size; x <= t.max; x += size {
			t.diamond(x, y, half, t.rand.Float32()*scale*2-scale)
		}

	}
	//t.print()
	if half == 1 {
		return
	}
	t.Divide(half)
}

func (t *Terrain) square(x int, y int, size int, offset float32) {
	val := average(t.getValue(x-size, y-size),
		t.getValue(x+size, y-size),
		t.getValue(x+size, y+size),
		t.getValue(x-size, y+size))
	t.setValue(x, y, val+offset)
}

func (t *Terrain) diamond(x int, y int, size int, offset float32) {
	val := average(t.getValue(x, y-size),
		t.getValue(x+size, y),
		t.getValue(x, y+size),
		t.getValue(x-size, y))
	t.setValue(x, y, val+offset)
}

func average(vs ...float32) float32 {
	var sum float32
	var total int
	for _, val := range vs {
		if val == -1 {
			continue
		}
		total++
		sum += val
	}
	return sum / float32(total)
}

func (t *Terrain) print() {
	for y := 0; y < t.size; y++ {
		for x := 0; x < t.size; x++ {
			fmt.Print(t.getValue(x, y), "\t")
		}
		fmt.Println()
	}
	fmt.Println()
}
