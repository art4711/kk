package kk

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"sync"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/math/fixed"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/gl"
)

type Tiles struct {
	ims *glutil.Images
	m   map[int]*glutil.Image
	mtx sync.Mutex
	sz  image.Point

	face font.Face
}

func NewTiles(glctx gl.Context) *Tiles {
	return &Tiles{ims: glutil.NewImages(glctx)}
}

func (t *Tiles) Release() {
	t.drop()
	t.ims.Release()
}

func (t *Tiles) drop() {
	for _, t := range t.m {
		t.Release()
	}
	t.m = make(map[int]*glutil.Image)
}

func (t *Tiles) Get(n int, sz image.Point) *glutil.Image {
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
			Size: float64((sz.X + sz.Y) / 8),
		})
	}
	if t.m[n] == nil {
		t.m[n] = t.genTex(n)
	}
	return t.m[n]
}

var pal = [...][3]float32{
	{0.0, 0.0, 1.0},
	{0.0, 1.0, 0.5},
	{1.0, 0.5, 0.0},
	{1.0, 0.0, 0.0},
}

func (t *Tiles) genTex(n int) *glutil.Image {
	img := t.ims.NewImage(t.sz.X, t.sz.Y)

	p := n / 6
	d2 := float32(n%6) / 6.0
	d1 := 1.0 - d2
	ic := image.NewUniform(color.RGBA{
		uint8((pal[p][0]*d1 + pal[p+1][0]*d2) * 255),
		uint8((pal[p][1]*d1 + pal[p+1][1]*d2) * 255),
		uint8((pal[p][2]*d1 + pal[p+1][2]*d2) * 255),
		255})
	if n == 0 {
		ic = image.NewUniform(color.RGBA{204, 204, 204, 255})
	}
	im := img.RGBA
	ul := t.sz.Div(20)
	lr := t.sz.Sub(ul)
	draw.Draw(im, im.Bounds(), image.NewUniform(color.RGBA{255, 255, 255, 255}), image.Point{}, draw.Src)
	draw.Draw(im, image.Rectangle{ul, lr}, ic, image.Point{}, draw.Src)

	if n > 0 {
		s := fmt.Sprintf("%d", 1<<uint(n))

		dot := fixed.P(t.sz.X/2, t.sz.Y/2)
		dot.Y += t.face.Metrics().Ascent / 2
		dot.X -= font.MeasureString(t.face, s) / 2
		d := font.Drawer{
			Dst:  im,
			Src:  image.Black,
			Face: t.face,
			Dot:  dot,
		}
		d.DrawString(s)
	}
	img.Upload()
	return img
}
