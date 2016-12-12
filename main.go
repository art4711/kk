package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"math/rand"
	"sync"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/math/fixed"

	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/geom"
	"golang.org/x/mobile/gl"

	"github.com/golang/freetype/truetype"
)

const width = 4
const height = 4

type field [height][width]int

func (f *field) r() {
	for {
		i := rand.Intn(height * width)
		c, r := i/width, i%width
		if f[c][r] == 0 {
			f[c][r] = 1
			break
		}
	}
}

func (f *field) set(n field) {
	add := n != *f
	*f = n
	if add {
		f.r()
	}
}

func (f *field) left() {
	n := field{}
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

func (f *field) right() {
	n := field{}
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

func (f *field) up() {
	n := field{}
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

func (f *field) down() {
	n := field{}
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

func newField() *field {
	f := &field{}
	f.r()
	f.r()
	return f
}

func main() {
	app.Main(func(a app.App) {
		var glctx gl.Context

		repaint := false
		wsz := size.Event{}
		f := newField()
		tiles := tiles{}

		var ims *glutil.Images

		for e := range a.Events() {
			switch e := a.Filter(e).(type) {
			case lifecycle.Event:
				switch e.Crosses(lifecycle.StageVisible) {
				case lifecycle.CrossOn:
					glctx, _ = e.DrawContext.(gl.Context)
					ims = glutil.NewImages(glctx)
					a.Send(paint.Event{})
				case lifecycle.CrossOff:
					glctx = nil
					return
				}

			case key.Event:
				if e.Code == key.CodeEscape {
					return
				}
				if e.Direction == key.DirPress {
					switch e.Code {
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
						a.Send(paint.Event{})
					}
				}

			case paint.Event:
				glctx.ClearColor(0, 0, 0, 1)
				glctx.Clear(gl.COLOR_BUFFER_BIT)
				// We actually want to create the tiles with correct pixel sizes so that the font looks good ...
				tsz := image.Point{wsz.WidthPx / width, wsz.HeightPx / height}
				// ..., but we need to render them in geom space.
				tpsz := geom.Point{wsz.WidthPt / geom.Pt(width), wsz.HeightPt / geom.Pt(height)}
				for y := 0; y < height; y++ {
					for x := 0; x < width; x++ {
						t := tiles.get(f[y][x], ims, tsz)
						tl := geom.Point{tpsz.X * geom.Pt(x), tpsz.Y * geom.Pt(y)}
						tr := geom.Point{tpsz.X * geom.Pt(x+1), tpsz.Y * geom.Pt(y)}
						bl := geom.Point{tpsz.X * geom.Pt(x), tpsz.Y * geom.Pt(y+1)}
						t.Draw(wsz, tl, tr, bl, image.Rectangle{Max: tsz})
					}
				}
				a.Publish()
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
	m   map[int]*glutil.Image
	mtx sync.Mutex
	sz  image.Point

	face font.Face
}

func (t *tiles) drop() {
	for _, t := range t.m {
		t.Release()
	}
	t.m = make(map[int]*glutil.Image)
}

func (t *tiles) get(n int, ims *glutil.Images, sz image.Point) *glutil.Image {
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
			Size: float64((sz.X + sz.Y) / 4),
		})
	}
	v := t.m[n]
	if v == nil {
		v = t.genTex(n, ims)
		t.m[n] = v
	}
	return v
}

func (t *tiles) genTex(n int, ims *glutil.Images) *glutil.Image {
	img := ims.NewImage(t.sz.X, t.sz.Y)

	ic := image.NewUniform(color.RGBA{20, 20 * uint8(n), 120, 255})
	if n == 0 {
		ic = image.NewUniform(color.RGBA{80, 80, 80, 255})
	}
	im := img.RGBA
	ul := t.sz.Div(20)
	lr := t.sz.Sub(ul)
	draw.Draw(im, im.Bounds(), image.NewUniform(color.RGBA{133, 95, 15, 255}), image.Point{}, draw.Src)
	draw.Draw(im, image.Rectangle{ul, lr}, ic, image.Point{}, draw.Src)

	if n > 0 {
		fc := image.Black
		dot := fixed.P(t.sz.X/3, t.sz.Y/2)
		dot.Y += t.face.Metrics().Ascent / 2
		d := font.Drawer{
			Dst:  im,
			Src:  fc,
			Face: t.face,
			Dot:  dot,
		}
		d.DrawString(fmt.Sprint(n))
	}
	img.Upload()
	return img
}
