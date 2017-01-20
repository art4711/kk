package kk

import "math/rand"

const width = 4
const height = 4

type Field [height][width]int

func (f *Field) r() {
	for {
		i := rand.Intn(height * width)
		c, r := i/width, i%width
		if f[c][r] == 0 {
			f[c][r] = 1
			break
		}
	}
}

func (f *Field) set(n Field) {
	add := n != *f
	*f = n
	if add {
		f.r()
	}
}

func (f *Field) W() int {
	return width
}

func (f *Field) H() int {
	return height
}

func (f *Field) Left() {
	n := Field{}
	for y := 0; y < height; y++ {
		last := 0
		c := 0
		for x := 0; x < width; x++ {
			v := (*f)[y][x]
			if v == 0 {
				continue
			}
			if last == v {
				n[y][c] = last + 1
				c++
				last = 0
			} else {
				if last != 0 {
					n[y][c] = last
					c++
				}
				last = v
			}
		}
		if last != 0 {
			n[y][c] = last
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
			v := (*f)[y][x]
			if v == 0 {
				continue
			}
			if last == v {
				n[y][c] = last + 1
				c--
				last = 0
			} else {
				if last != 0 {
					n[y][c] = last
					c--
				}
				last = v
			}
		}
		if last != 0 {
			n[y][c] = last
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
			v := (*f)[y][x]
			if v == 0 {
				continue
			}
			if last == v {
				n[r][x] = last + 1
				r++
				last = 0
			} else {
				if last != 0 {
					n[r][x] = last
					r++
				}
				last = v
			}
		}
		if last != 0 {
			n[r][x] = last
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
			v := (*f)[y][x]
			if v == 0 {
				continue
			}
			if last == v {
				n[r][x] = last + 1
				r--
				last = 0
			} else {
				if last != 0 {
					n[r][x] = last
					r--
				}
				last = v
			}
		}
		if last != 0 {
			n[r][x] = last
		}
	}
	f.set(n)
}

func NewField() *Field {
	f := &Field{}
	f.r()
	f.r()
	return f
}
