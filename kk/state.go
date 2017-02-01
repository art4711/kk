package kk

import (
	"image"
	"log"
	"math"

	"golang.org/x/exp/shiny/unit"
	"golang.org/x/exp/shiny/widget"
	"golang.org/x/exp/shiny/widget/node"
	"golang.org/x/exp/shiny/widget/theme"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/geom"
	"golang.org/x/mobile/gl"
)

type butt struct {
	b int
	w *widget.Uniform
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

	buttons []butt

	t theme.Theme

	touchStart image.Point
}

func New() *State {
	s := &State{}
	s.f.Init()
	return s
}

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

func (s *State) clickOrTouch(ps, pe image.Point) interface{} {
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
	case ps.In(s.buttons[0].r):
		return EvSave{}
	case ps.In(s.buttons[1].r):
		return EvLoad{}
	case ps.In(s.buttons[2].r):
		return EvReset{}
	}
	return nil
}

func (s *State) EvFilter(ei interface{}) interface{} {
	if e, ok := ei.(touch.Event); ok {
		switch e.Type {
		case touch.TypeBegin:
			s.touchStart = image.Pt(int(e.X), int(e.Y))
		case touch.TypeEnd:
			ne := s.clickOrTouch(s.touchStart, image.Pt(int(e.X), int(e.Y)))
			if ne != nil {
				return ne
			}
		}
	}
	if e, ok := ei.(mouse.Event); ok && e.Button == mouse.ButtonLeft {
		switch e.Direction {
		case mouse.DirPress:
			s.touchStart = image.Pt(int(e.X), int(e.Y))
		case mouse.DirRelease:
			ne := s.clickOrTouch(s.touchStart, image.Pt(int(e.X), int(e.Y)))
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

// helper.
func stretch(n node.Node, alongWeight int) node.Node {
	return widget.WithLayoutData(n, widget.FlowLayoutData{
		AlongWeight:  alongWeight,
		ExpandAlong:  true,
		ShrinkAlong:  true,
		ExpandAcross: true,
		ShrinkAcross: true,
	})
}

// helper to extract screen coords for a widget.
func widgetScreenRect(e *node.Embed) image.Rectangle {
	r := e.Rect
	for e.Parent != nil {
		e = e.Parent
		r = r.Add(e.Rect.Min)
	}
	return r
}

// helper to translate image points to geom points for glutil.Images.
func (s *State) ip2gp(ip image.Point) geom.Point {
	return geom.Point{
		X: geom.Pt(ip.X) / geom.Pt(s.wsz.PixelsPerPt),
		Y: geom.Pt(ip.Y) / geom.Pt(s.wsz.PixelsPerPt),
	}
}

func (s *State) setSize(e size.Event) {
	s.wsz = e
	s.t = theme.Theme{DPI: float64(unit.PointsPerInch * e.PixelsPerPt)}

	horizontal := e.Orientation == size.OrientationLandscape || e.WidthPx > e.HeightPx

	padPx := float64(e.WidthPx / 50)
	butPx := float64(e.WidthPx / 6)
	bAx := widget.AxisHorizontal
	aAx := widget.AxisVertical
	if horizontal {
		padPx = float64(e.HeightPx / 50)
		butPx = float64(e.HeightPx / 6)
		bAx = widget.AxisVertical
		aAx = widget.AxisHorizontal
	}

	// We abuse shiny widgets to do the layout for us.

	s.buttons = []butt{{b: Save}, {b: Load}, {b: Reset}}

	for i := range s.buttons {
		s.buttons[i].w = widget.NewUniform(theme.Light, nil)
	}

	bb := widget.NewUniform(theme.Neutral,
		widget.NewPadder(widget.AxisBoth, unit.Pixels(padPx),
			widget.NewFlow(bAx,
				widget.NewSizer(unit.Pixels(butPx), unit.Pixels(butPx), s.buttons[0].w),
				widget.NewSizer(unit.Pixels(butPx), unit.Pixels(butPx), s.buttons[1].w),
				stretch(widget.NewSpace(), 1),
				widget.NewSizer(unit.Pixels(butPx), unit.Pixels(butPx), s.buttons[2].w),
			),
		),
	)
	// field
	f := widget.NewUniform(theme.Light, nil)

	all := widget.NewFlow(aAx,
		stretch(bb, 0),
		stretch(widget.NewPadder(widget.AxisBoth, unit.Pixels(padPx), f), 1),
	)
	// do the layout.
	all.Measure(&s.t, e.WidthPx, e.HeightPx)
	all.Rect = image.Rectangle{Max: image.Pt(e.WidthPx, e.HeightPx)}
	all.Layout(&s.t)

	r := widgetScreenRect(&f.Embed)
	dx, dy := r.Dx(), r.Dy()

	// square the field rectangle.
	if dx > dy {
		r.Min.X += dx - dy
	} else {
		r.Min.Y += dy - dx
	}

	s.fr = r

	s.ful = r.Min
	s.tsz.X = r.Dx() / s.f.W()
	s.tsz.Y = r.Dy() / s.f.H()

	for i := range s.buttons {
		s.buttons[i].r = widgetScreenRect(&s.buttons[i].w.Embed)
	}

	s.tiles.SetSz(s.tsz, s.buttons[0].r.Size())
}

func (s *State) drawButt() {
	for i := range s.buttons {
		r := widgetScreenRect(&s.buttons[i].w.Embed)
		img := s.tiles.Get(s.buttons[i].b)
		tl := r.Min
		br := r.Max
		tr := image.Pt(br.X, tl.Y)
		bl := image.Pt(tl.X, br.Y)
		img.Draw(s.wsz, s.ip2gp(tl), s.ip2gp(tr), s.ip2gp(bl), image.Rectangle{Max: r.Size()})
	}
}

func (s *State) fRectBounds(x, y int) (geom.Point, geom.Point, geom.Point) {
	return s.ip2gp(image.Pt(s.ful.X+s.tsz.X*x, s.ful.Y+s.tsz.Y*y)),
		s.ip2gp(image.Pt(s.ful.X+s.tsz.X*(x+1), s.ful.Y+s.tsz.Y*y)),
		s.ip2gp(image.Pt(s.ful.X+s.tsz.X*x, s.ful.Y+s.tsz.Y*(y+1)))
}

func (s *State) draw(pub func()) {
	if s.glctx == nil || s.ful.X == 0 {
		return
	}
	s.glctx.ClearColor(1, 1, 1, 1)
	s.glctx.Clear(gl.COLOR_BUFFER_BIT)

	s.drawButt()

	w, h := s.f.W(), s.f.H()
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			t := s.tiles.Get(s.f[y][x])
			tl, tr, bl := s.fRectBounds(x, y)
			t.Draw(s.wsz, tl, tr, bl, image.Rectangle{Max: s.tsz})
		}
	}
	pub()
}
