package kk

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/math/fixed"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/gl"
)

/*
 * This file contains all the bitmap rendering.
 * We use glutil.Images for sprites, but at some point
 * it might be a better idea to take control over the
 * shaders and maybe do something interesting with them.
 */

type Tiles struct {
	ims *glutil.Images
	m   map[T]*glutil.Image

	buttSz  image.Point
	sz      image.Point
	scoreSz image.Point

	face      font.Face
	buttFace  font.Face
	scoreFace font.Face
}

type T interface {
	Gen(*Tiles) *glutil.Image
}

func (t *Tiles) SetCtx(ctx gl.Context) {
	t.Release()
	t.ims = glutil.NewImages(ctx)
}

var fnt *truetype.Font

func init() {
	f, err := truetype.Parse(gobold.TTF)
	if err != nil {
		log.Fatal(err)
	}
	fnt = f
}

func (t *Tiles) SetSz(sz, buttSz image.Point, scoreSz image.Point) {
	if t.sz != sz || t.buttSz != buttSz || t.scoreSz != scoreSz {
		t.drop()
		t.sz = sz
		t.buttSz = buttSz
		t.scoreSz = scoreSz
		t.face = truetype.NewFace(fnt, &truetype.Options{
			Size: float64((sz.X + sz.Y) / 8),
		})
		t.buttFace = truetype.NewFace(fnt, &truetype.Options{
			Size: float64(buttSz.X / 4),
		})
		ss := float64(scoreSz.X) * 1.75 // XXX - Why?
		if scoreSz.X > scoreSz.Y {
			ss = float64(scoreSz.Y)
		}
		t.scoreFace = truetype.NewFace(fnt, &truetype.Options{
			Size: ss,
		})
	}
}

func (t *Tiles) Release() {
	t.drop()
	if t.ims != nil {
		t.ims.Release()
		t.ims = nil
	}
}

func (t *Tiles) drop() {
	for _, t := range t.m {
		t.Release()
	}
	t.m = make(map[T]*glutil.Image)
}

func (t *Tiles) Get(tl T) *glutil.Image {
	if t.m[tl] == nil {
		t.m[tl] = tl.Gen(t)
	}
	return t.m[tl]
}

type Butt struct {
	Label  string
	Invert bool
	Fade   bool
}

func (b Butt) Gen(t *Tiles) *glutil.Image {
	// can't be a pointer receiver because we want the interface
	// value compare the struct contents, not pointers for the
	// map.

	tile := t.ims.NewImage(t.buttSz.X, t.buttSz.Y)
	s := b.Label

	img := tile.RGBA
	r := img.Bounds()

	fg := image.Black
	bg := image.Transparent

	if b.Fade {
		fg = image.NewUniform(color.RGBA{0, 0, 0, 90})
	}
	if b.Invert {
		fg, bg = bg, fg
	}

	draw.Draw(img, r, bg, image.Point{}, draw.Src)

	draw.Draw(img, image.Rect(r.Min.X+2, r.Min.Y+2, r.Max.X-2, r.Min.Y+4), fg, image.Point{}, draw.Src)
	draw.Draw(img, image.Rect(r.Min.X+2, r.Min.Y+2, r.Min.X+4, r.Max.Y-2), fg, image.Point{}, draw.Src)
	draw.Draw(img, image.Rect(r.Max.X-4, r.Min.Y+2, r.Max.X-2, r.Max.Y-2), fg, image.Point{}, draw.Src)
	draw.Draw(img, image.Rect(r.Min.X+2, r.Max.Y-4, r.Max.X-2, r.Max.Y-2), fg, image.Point{}, draw.Src)

	dot := fixed.P(t.buttSz.X/2, t.buttSz.Y/2)
	dot.Y += t.buttFace.Metrics().Ascent / 3
	dot.X -= font.MeasureString(t.buttFace, s) / 2
	d := font.Drawer{
		Dst:  img,
		Src:  fg,
		Face: t.buttFace,
		Dot:  dot,
	}
	d.DrawString(s)
	tile.Upload()

	return tile
}

var pal = [...][3]float32{
	{0.0, 0.0, 1.0},
	{0.0, 1.0, 0.5},
	{1.0, 0.5, 0.0},
	{1.0, 0.0, 0.0},
}

type FT struct {
	n    int
	over bool
}

func (ft FT) Gen(t *Tiles) *glutil.Image {
	img := t.ims.NewImage(t.sz.X, t.sz.Y)

	p := ft.n / 6
	d2 := float32(ft.n%6) / 6.0
	d1 := 1.0 - d2
	ic := image.NewUniform(color.RGBA{
		uint8((pal[p][0]*d1 + pal[p+1][0]*d2) * 255),
		uint8((pal[p][1]*d1 + pal[p+1][1]*d2) * 255),
		uint8((pal[p][2]*d1 + pal[p+1][2]*d2) * 255),
		255})
	if ft.n == 0 {
		ic = image.NewUniform(color.RGBA{204, 204, 204, 255})
	}
	im := img.RGBA
	ul := t.sz.Div(20)
	lr := t.sz.Sub(ul)
	borderColor := image.NewUniform(color.RGBA{255, 255, 255, 204})
	draw.Draw(im, im.Bounds(), borderColor, image.Point{}, draw.Src)
	draw.Draw(im, image.Rectangle{ul, lr}, ic, image.Point{}, draw.Src)

	if ft.n > 0 {
		s := fmt.Sprintf("%d", 1<<uint(ft.n))

		dot := fixed.P(t.sz.X/2, t.sz.Y/2)
		dot.Y += t.face.Metrics().Ascent / 3
		dot.X -= font.MeasureString(t.face, s) / 2
		d := font.Drawer{
			Dst:  im,
			Face: t.face,
			Dot:  dot,
		}
		if ft.over {
			d.Src = image.NewUniform(color.RGBA{255, 0, 0, 255})
		} else {
			d.Src = image.Black
		}
		d.DrawString(s)
	}
	img.Upload()
	return img
}

type ST int

func (n ST) Gen(t *Tiles) *glutil.Image {
	tile := t.ims.NewImage(t.scoreSz.X, t.scoreSz.Y)

	s := fmt.Sprint(n)

	img := tile.RGBA

	dot := fixed.P(t.scoreSz.X/2, t.scoreSz.Y/2)
	dot.Y += t.scoreFace.Metrics().Ascent / 3
	dot.X -= font.MeasureString(t.scoreFace, s) / 2
	d := font.Drawer{
		Dst:  img,
		Src:  image.Black,
		Face: t.scoreFace,
		Dot:  dot,
	}
	d.DrawString(s)
	tile.Upload()

	return tile
}
