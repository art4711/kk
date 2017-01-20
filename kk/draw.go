package kk

import (
	"image"

	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/geom"
	"golang.org/x/mobile/gl"
)

func Draw(glctx gl.Context, wsz size.Event, f *Field, tiles *Tiles) {
	glctx.ClearColor(0, 0, 0, 1)
	glctx.Clear(gl.COLOR_BUFFER_BIT)
	// We actually want to create the tiles with correct pixel sizes so the font looks good ...
	tsz := image.Point{wsz.WidthPx / f.W(), wsz.HeightPx / f.H()}
	// ..., but we need to render them in geom space.
	tpsz := geom.Point{wsz.WidthPt / geom.Pt(f.W()), wsz.HeightPt / geom.Pt(f.H())}
	for y := 0; y < f.H(); y++ {
		for x := 0; x < f.W(); x++ {
			t := tiles.Get(f[y][x], tsz)
			tl := geom.Point{tpsz.X * geom.Pt(x), tpsz.Y * geom.Pt(y)}
			tr := geom.Point{tpsz.X * geom.Pt(x+1), tpsz.Y * geom.Pt(y)}
			bl := geom.Point{tpsz.X * geom.Pt(x), tpsz.Y * geom.Pt(y+1)}
			t.Draw(wsz, tl, tr, bl, image.Rectangle{Max: tsz})
		}
	}
}
