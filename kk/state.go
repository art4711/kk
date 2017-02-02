package kk

import (
	"image"
	"log"
	"math"

	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/geom"
	"golang.org/x/mobile/gl"
)

type butt struct {
	b Butt
	r image.Rectangle
}

type State struct {
	tiles Tiles
	f     Field

	saved *Field

	glctx gl.Context
	wsz   size.Event

	ful image.Point
	tsz image.Point
	fr  image.Rectangle

	buttons map[string]*butt
	scores  [8]image.Rectangle

	touchStart image.Point
}

func New() *State {
	s := &State{}
	s.f.Init()
	s.buttons = make(map[string]*butt)
	s.buttons["save"] = &butt{b: Butt{Label: "Save"}}
	s.buttons["load"] = &butt{b: Butt{Label: "Load", Fade: true}}
	s.buttons["reset"] = &butt{b: Butt{Label: "Reset"}}
	return s
}

// Event handler. Main entry point from platform specific code.
func (s *State) Handle(ei interface{}, pub func()) bool {
	switch e := ei.(type) {
	case lifecycle.Event:
		if e.To == lifecycle.StageDead {
			return false
		}
		switch e.Crosses(lifecycle.StageVisible) {
		case lifecycle.CrossOn:
			s.glctx, _ = e.DrawContext.(gl.Context)
			s.tiles.SetCtx(s.glctx)
		case lifecycle.CrossOff:
			s.glctx = nil
			s.tiles.Release()
			return true
		}
	case EvR:
		s.f.Right()
	case EvL:
		s.f.Left()
	case EvU:
		s.f.Up()
	case EvD:
		s.f.Down()
	case EvQ:
		return false
	case EvReset:
		s.f.Init()
	case EvSave:
		f := s.f
		s.buttons["load"].b.Fade = false
		s.saved = &f
	case EvLoad:
		if s.saved != nil {
			s.f = *s.saved
		}
	case paint.Event:
	case size.Event:
		s.setSize(e)
	case error:
		log.Print(e)
		return true
	default:
		return true
	}
	s.draw(pub)
	return true
}

func (s *State) clickOrTouch(start bool, p image.Point) interface{} {
	if start {
		s.touchStart = p
		for i := range s.buttons {
			if p.In(s.buttons[i].r) {
				s.buttons[i].b.Invert = true
			}
		}
		return paint.Event{}
	}
	ps := s.touchStart
	pe := p
	switch {
	case ps.In(s.fr):
		dp := pe.Sub(ps)
		if math.Abs(float64(dp.X)) > math.Abs(float64(dp.Y)) {
			if dp.X < 0 {
				return EvL{}
			} else {
				return EvR{}
			}
		} else {
			if dp.Y < 0 {
				return EvU{}
			} else {
				return EvD{}
			}
		}
	case ps.In(s.buttons["save"].r):
		s.buttons["save"].b.Invert = false
		return EvSave{}
	case ps.In(s.buttons["load"].r):
		s.buttons["load"].b.Invert = false
		return EvLoad{}
	case ps.In(s.buttons["reset"].r):
		s.buttons["reset"].b.Invert = false
		return EvReset{}
	}
	return nil
}

func (s *State) EvFilter(ei interface{}) interface{} {
	if e, ok := ei.(touch.Event); ok {
		switch e.Type {
		case touch.TypeBegin, touch.TypeEnd:
			ne := s.clickOrTouch(e.Type == touch.TypeBegin, image.Pt(int(e.X), int(e.Y)))
			if ne != nil {
				return ne
			}
		}
	}
	if e, ok := ei.(mouse.Event); ok && e.Button == mouse.ButtonLeft {
		switch e.Direction {
		case mouse.DirPress, mouse.DirRelease:
			ne := s.clickOrTouch(e.Direction == mouse.DirPress, image.Pt(int(e.X), int(e.Y)))
			if ne != nil {
				return ne
			}
		}
	}
	return ei
}

type EvL struct{}
type EvR struct{}
type EvU struct{}
type EvD struct{}
type EvQ struct{}
type EvReset struct{}
type EvLoad struct{}
type EvSave struct{}

func (s *State) fRectBounds(x, y int) (geom.Point, geom.Point, geom.Point) {
	return s.ip2gp(image.Pt(s.ful.X+s.tsz.X*x, s.ful.Y+s.tsz.Y*y)),
		s.ip2gp(image.Pt(s.ful.X+s.tsz.X*(x+1), s.ful.Y+s.tsz.Y*y)),
		s.ip2gp(image.Pt(s.ful.X+s.tsz.X*x, s.ful.Y+s.tsz.Y*(y+1)))
}

func (s *State) rect2gps(r image.Rectangle) (geom.Point, geom.Point, geom.Point) {
	tl := r.Min
	br := r.Max
	tr := image.Pt(br.X, tl.Y)
	bl := image.Pt(tl.X, br.Y)
	return s.ip2gp(tl), s.ip2gp(tr), s.ip2gp(bl)
}

func (s *State) draw(pub func()) {
	if s.glctx == nil || s.ful.X == 0 {
		return
	}
	s.glctx.ClearColor(0.95, 0.97, 1, 1)
	s.glctx.Clear(gl.COLOR_BUFFER_BIT)

	over := s.f.GameOver()

	// Draw the buttons.
	for i := range s.buttons {
		r := s.buttons[i].r
		img := s.tiles.Get(s.buttons[i].b)
		tl, tr, bl := s.rect2gps(r)
		img.Draw(s.wsz, tl, tr, bl, image.Rectangle{Max: r.Size()})
	}

	// Draw score.
	mask := 1
	for _ = range s.scores {
		mask *= 10
	}
	for i := range s.scores {
		r := s.scores[i]
		mask /= 10
		v := (s.f.score / mask) % 10
		img := s.tiles.Get(ST(v))
		tl, tr, bl := s.rect2gps(r)
		img.Draw(s.wsz, tl, tr, bl, image.Rectangle{Max: r.Size()})
	}

	// Draw the field.
	w, h := s.f.W(), s.f.H()
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			t := s.tiles.Get(FT{s.f.f[y][x], over})
			tl, tr, bl := s.fRectBounds(x, y)
			t.Draw(s.wsz, tl, tr, bl, image.Rectangle{Max: s.tsz})
		}
	}
	pub()
}
