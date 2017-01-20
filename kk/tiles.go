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
			Size: float64((sz.X + sz.Y) / 4),
		})
	}
	v := t.m[n]
	if v == nil {
		v = t.genTex(n)
		t.m[n] = v
	}
	return v
}

func (t *Tiles) genTex(n int) *glutil.Image {
	img := t.ims.NewImage(t.sz.X, t.sz.Y)

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
