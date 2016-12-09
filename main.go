package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"math/rand"
	"sync"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/math/fixed"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"

	"github.com/golang/freetype/truetype"
)

type field struct {
	v [4][4]int
}

func (f *field) r() {
	for {
		i := rand.Intn(16)
		c, r := i/4, i%4
		if f.v[c][r] == 0 {
			f.v[c][r] = 1
			break
		}
	}
}

func (f *field) set(n [4][4]int) {
	add := n != f.v
	f.v = n
	if add {
		f.r()
	}
}

func (f *field) left() {
	n := [4][4]int{}
	for y := 0; y < 4; y++ {
		last := 0
		c := 0
		for x := 0; x < 4; x++ {
			v := f.v[y][x]
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

func (f *field) right() {
	n := [4][4]int{}
	for y := 0; y < 4; y++ {
		last := 0
		c := 3
		for x := 3; x >= 0; x-- {
			v := f.v[y][x]
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

func (f *field) up() {

}

func (f *field) down() {

}

func newField() *field {
	f := &field{}
	f.r()
	f.r()
	return f
}

func main() {
	driver.Main(func(s screen.Screen) {
		w, err := s.NewWindow(nil)
		if err != nil {
			log.Fatal(err)
		}
		defer w.Release()

		repaint := false
		wsz := size.Event{}
		f := newField()
		tiles := tiles{}

		for {
			switch e := w.NextEvent().(type) {
			case lifecycle.Event:
				if e.To == lifecycle.StageDead {
					return
				}

			case key.Event:
				if e.Direction == key.DirPress {
					switch e.Code {
					case key.CodeEscape:
						return
					case key.CodeLeftArrow:
						f.left()
					case key.CodeRightArrow:
						f.right()
					case key.CodeUpArrow:
						f.up()
					case key.CodeDownArrow:
						f.down()
					}
					if !repaint {
						repaint = true
						w.Send(paint.Event{})
					}
				}

			case paint.Event:
				var wg sync.WaitGroup
				tsz := image.Point{wsz.WidthPx / 4, wsz.HeightPx / 4}
				for y := 0; y < 4; y++ {
					for x := 0; x < 4; x++ {
						wg.Add(1)
						go func(x, y int) {
							defer wg.Done()
							t := tiles.get(f.v[y][x], s, tsz)
							w.Copy(image.Point{tsz.X * x, tsz.Y * y}, t, image.Rectangle{Max: tsz}, screen.Src, nil)
						}(x, y)
					}
				}
				wg.Wait()
				w.Publish()
				repaint = false
			case size.Event:
				wsz = e

			case error:
				log.Print(e)
			}
		}
	})
}

type tiles struct {
	m   map[int]screen.Texture
	mtx sync.Mutex
	sz  image.Point

	face font.Face
}

func (t *tiles) drop() {
	for _, t := range t.m {
		t.Release()
	}
	t.m = make(map[int]screen.Texture)
}

func (t *tiles) get(n int, s screen.Screen, sz image.Point) screen.Texture {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	if t.sz != sz {
		t.drop()
		t.sz = sz
		f, err := truetype.Parse(gobold.TTF)
		if err != nil {
			log.Fatal(err)
		}
		// notice that the use of face is protected by the mutex
		t.face = truetype.NewFace(f, &truetype.Options{
			Size: float64(sz.X / 2),
		})
	}
	v := t.m[n]
	if v != nil {
		return v
	}
	v, err := t.genTex(n, s)
	if err != nil {
		log.Fatal(err)
	}
	t.m[n] = v
	return v
}

func (t *tiles) genTex(n int, s screen.Screen) (screen.Texture, error) {
	tex, err := s.NewTexture(t.sz)
	if err != nil {
		return nil, err
	}
	buf, err := s.NewBuffer(t.sz)
	if err != nil {
		return nil, err
	}
	ic := image.NewUniform(color.RGBA{20, 20 * uint8(n), 120, 255})
	if n == 0 {
		ic = image.NewUniform(color.RGBA{80, 80, 80, 255})
	}
	im := buf.RGBA()
	draw.Draw(im, im.Bounds(), ic, image.Point{}, draw.Src)

	if n > 0 {
		fc := image.Black
		// These are completely made up numbers that seem to work, I
		// don't actually know why.
		dot := fixed.P(t.sz.X/3, t.sz.Y/8)
		dot.Y += t.face.Metrics().Ascent
		d := font.Drawer{
			Dst:  im,
			Src:  fc,
			Face: t.face,
			Dot:  dot,
		}
		d.DrawString(fmt.Sprint(n))
	}

	tex.Upload(image.Point{}, buf, image.Rectangle{Max: t.sz})
	buf.Release()
	return tex, nil
}
