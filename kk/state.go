package kk

import (
	"image"
	"log"

	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/geom"
	"golang.org/x/mobile/gl"
)

type State struct {
	tiles *Tiles
	f     *Field
	wsz   size.Event
	glctx gl.Context
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

	case error:
		log.Print(e)
	}
	return
}

func (s *State) Draw() {
	s.glctx.ClearColor(0, 0, 0, 1)
	s.glctx.Clear(gl.COLOR_BUFFER_BIT)
	// We actually want to create the tiles with correct pixel sizes so the font looks good ...
	tsz := image.Point{s.wsz.WidthPx / s.f.W(), s.wsz.HeightPx / s.f.H()}
	// ..., but we need to render them in geom space.
	tpsz := geom.Point{s.wsz.WidthPt / geom.Pt(s.f.W()), s.wsz.HeightPt / geom.Pt(s.f.H())}
	for y := 0; y < s.f.H(); y++ {
		for x := 0; x < s.f.W(); x++ {
			t := s.tiles.Get(s.f[y][x], tsz)
			tl := geom.Point{tpsz.X * geom.Pt(x), tpsz.Y * geom.Pt(y)}
			tr := geom.Point{tpsz.X * geom.Pt(x+1), tpsz.Y * geom.Pt(y)}
			bl := geom.Point{tpsz.X * geom.Pt(x), tpsz.Y * geom.Pt(y+1)}
			t.Draw(s.wsz, tl, tr, bl, image.Rectangle{Max: tsz})
		}
	}
}
