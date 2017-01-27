package kk

import (
	"image"
	"log"

	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/geom"
	"golang.org/x/mobile/gl"
)

type State struct {
	tiles Tiles
	f     *Field
	wsz   size.Event
	glctx gl.Context

	ful geom.Point
	fst geom.Point
	tsz image.Point
}

func New() *State {
	return &State{f: NewField()}
}

type EvL struct{}
type EvR struct{}
type EvU struct{}
type EvD struct{}
type EvQ struct{}

const margin = 0.02

func (s *State) setSize(e size.Event) {
	s.wsz = e
	x, y := e.WidthPt, e.HeightPt

	if x > y {
		s.ful.X = x - y*(1-margin)
		s.ful.Y = y * margin
		s.fst.X = (y * (1 - 2*margin) / geom.Pt(s.f.W()))
		s.fst.Y = (y * (1 - 2*margin) / geom.Pt(s.f.H()))
	} else {
		s.ful.X = x * margin
		s.ful.Y = y - x*(1-margin)
		s.fst.X = (x * (1 - 2*margin) / geom.Pt(s.f.W()))
		s.fst.Y = (x * (1 - 2*margin) / geom.Pt(s.f.H()))
	}
	log.Print(e)
	s.tsz.X = int(s.fst.X.Px(e.PixelsPerPt))
	s.tsz.Y = int(s.fst.Y.Px(e.PixelsPerPt))
	s.tiles.SetSz(s.tsz)
}

func (s *State) fRectBounds(x, y int) (geom.Point, geom.Point, geom.Point) {
	return geom.Point{s.ful.X + s.fst.X*geom.Pt(x), s.ful.Y + s.fst.X*geom.Pt(y)},
		geom.Point{s.ful.X + s.fst.X*geom.Pt(x+1), s.ful.Y + s.fst.X*geom.Pt(y)},
		geom.Point{s.ful.X + s.fst.X*geom.Pt(x), s.ful.Y + s.fst.X*geom.Pt(y+1)}
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
			s.tiles.SetCtx(s.glctx)
			repaint = true
		case lifecycle.CrossOff:
			s.glctx = nil
			s.tiles.Release()
			return
		}
	case EvR:
		s.f.Right()
		repaint = true
	case EvL:
		s.f.Left()
		repaint = true
	case EvU:
		s.f.Up()
		repaint = true
	case EvD:
		s.f.Down()
		repaint = true
	case EvQ:
		quit = true
	case paint.Event:
		s.Draw()
		publish = true
	case size.Event:
		s.setSize(e)

	case touch.Event:
	case error:
		log.Print(e)
	}
	return
}

func (s *State) Draw() {
	s.glctx.ClearColor(1, 1, 1, 1)
	s.glctx.Clear(gl.COLOR_BUFFER_BIT)

	w, h := s.f.W(), s.f.H()

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			t := s.tiles.Get(s.f[y][x])
			tl, tr, bl := s.fRectBounds(x, y)
			t.Draw(s.wsz, tl, tr, bl, image.Rectangle{Max: s.tsz})
		}
	}
}
