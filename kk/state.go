package kk

import (
	"image"
	"log"
	"math"

	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/geom"
	"golang.org/x/mobile/gl"
)

type State struct {
	tiles *Tiles
	f     *Field
	wsz   size.Event
	glctx gl.Context

	touchStart touch.Event
}

func New() *State {
	return &State{f: NewField()}
}

func (s *State) Handle(ei interface{}) (repaint bool, quit bool, publish bool) {
	switch e := ei.(type) {
	case lifecycle.Event:
		if e.To == lifecycle.StageDead {
			quit = true
			return
		}
		switch e.Crosses(lifecycle.StageVisible) {
		case lifecycle.CrossOn:
			s.glctx, _ = e.DrawContext.(gl.Context)
			s.tiles = NewTiles(s.glctx)
			repaint = true
		case lifecycle.CrossOff:
			s.glctx = nil
			quit = true
			return
		}

	case key.Event:
		if e.Code == key.CodeEscape {
			quit = true
			return
		}
		if e.Direction == key.DirPress {
			switch e.Code {
			case key.CodeLeftArrow:
				s.f.Left()
			case key.CodeRightArrow:
				s.f.Right()
			case key.CodeUpArrow:
				s.f.Up()
			case key.CodeDownArrow:
				s.f.Down()
			}
			repaint = true
		}

	case paint.Event:
		s.Draw()
		publish = true
	case size.Event:
		s.wsz = e

	case touch.Event:
		switch e.Type {
		case touch.TypeBegin:
			s.touchStart = e
		case touch.TypeEnd:
			x, y := e.X-s.touchStart.X, e.Y-s.touchStart.Y
			if math.Abs(float64(x)) > math.Abs(float64(y)) {
				if x < 0 {
					s.f.Left()
				} else {
					s.f.Right()
				}
			} else {
				if y < 0 {
					s.f.Up()
				} else {
					s.f.Down()
				}
			}
			repaint = true
		}
	case error:
		log.Print(e)
	}
	return
}

func (s *State) Draw() {
	s.glctx.ClearColor(0, 0, 0, 1)
	s.glctx.Clear(gl.COLOR_BUFFER_BIT)

	w, h := s.f.W(), s.f.H()

	// We actually want to create the tiles with correct pixel sizes so the font looks good ...
	tsz := image.Point{s.wsz.WidthPx / w, s.wsz.HeightPx / h}
	// ..., but we need to render them in geom space.
	tpsz := geom.Point{s.wsz.WidthPt / geom.Pt(w), s.wsz.HeightPt / geom.Pt(h)}
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			t := s.tiles.Get(s.f[y][x], tsz)
			tl := geom.Point{tpsz.X * geom.Pt(x), tpsz.Y * geom.Pt(y)}
			tr := geom.Point{tpsz.X * geom.Pt(x+1), tpsz.Y * geom.Pt(y)}
			bl := geom.Point{tpsz.X * geom.Pt(x), tpsz.Y * geom.Pt(y+1)}
			t.Draw(s.wsz, tl, tr, bl, image.Rectangle{Max: tsz})
		}
	}
}
