package kk

import (
	"math/rand"

	"github.com/art4711/unpredictable"
)

const width = 4
const height = 4

type Field struct {
	f     [height][width]int
	score int
}

func (f *Field) Init() {
	*f = Field{}
	f.r()
	f.r()
}

func (f *Field) GameOver() bool {
	for y := 0; y < height; y++ {
		for x := 0; x < width-1; x++ {
			if f.f[y][x] == 0 || f.f[y][x] == f.f[y][x+1] {
				return false
			}
		}
	}
	for y := 0; y < height-1; y++ {
		for x := 0; x < width; x++ {
			if f.f[y][x] == 0 || f.f[y][x] == f.f[y+1][x] {
				return false
			}
		}
	}
	return true
}

func (f *Field) Full() bool {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if f.f[y][x] == 0 {
				return false
			}
		}
	}
	return true
}

// If you're allergic to unpredictable numbers, just replace the line
// below with:
//     var rsrc = rand.New(rand.NewSource(0))
var rsrc = rand.New(unpredictable.NewMathRandSource())

func (f *Field) r() {
	if f.Full() {
		return
	}
	for {
		i := rsrc.Intn(height * width)
		c, r := i/width, i%width
		if f.f[c][r] == 0 {
			f.f[c][r] = 1
			break
		}
	}
}

func (f *Field) set(n Field) {
	add := n.f != f.f
	f.f = n.f
	if add {
		f.score += n.score
		f.r()
	}
}

func (f *Field) W() int {
	return width
}

func (f *Field) H() int {
	return height
}

func (f *Field) merged(y, x, val int) {
	f.f[y][x] = val + 1
	f.score += val
}

func (f *Field) Left() {
	n := Field{}
	for y := 0; y < height; y++ {
		last := 0
		c := 0
		for x := 0; x < width; x++ {
			v := f.f[y][x]
			if v == 0 {
				continue
			}
			if last == v {
				n.merged(y, c, last)
				c++
				last = 0
			} else {
				if last != 0 {
					n.f[y][c] = last
					c++
				}
				last = v
			}
		}
		if last != 0 {
			n.f[y][c] = last
		}
	}
	f.set(n)
}

func (f *Field) Right() {
	n := Field{}
	for y := 0; y < height; y++ {
		last := 0
		c := width - 1
		for x := width - 1; x >= 0; x-- {
			v := f.f[y][x]
			if v == 0 {
				continue
			}
			if last == v {
				n.merged(y, c, last)
				c--
				last = 0
			} else {
				if last != 0 {
					n.f[y][c] = last
					c--
				}
				last = v
			}
		}
		if last != 0 {
			n.f[y][c] = last
		}
	}
	f.set(n)
}

func (f *Field) Up() {
	n := Field{}
	for x := 0; x < width; x++ {
		last := 0
		r := 0
		for y := 0; y < height; y++ {
			v := f.f[y][x]
			if v == 0 {
				continue
			}
			if last == v {
				n.merged(r, x, last)
				r++
				last = 0
			} else {
				if last != 0 {
					n.f[r][x] = last
					r++
				}
				last = v
			}
		}
		if last != 0 {
			n.f[r][x] = last
		}
	}
	f.set(n)
}

func (f *Field) Down() {
	n := Field{}
	for x := 0; x < width; x++ {
		last := 0
		r := height - 1
		for y := height - 1; y >= 0; y-- {
			v := f.f[y][x]
			if v == 0 {
				continue
			}
			if last == v {
				n.merged(r, x, last)
				r--
				last = 0
			} else {
				if last != 0 {
					n.f[r][x] = last
					r--
				}
				last = v
			}
		}
		if last != 0 {
			n.f[r][x] = last
		}
	}
	f.set(n)
}
