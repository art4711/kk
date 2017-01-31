package kk

import (
	"image"
	"log"

	"golang.org/x/exp/shiny/unit"
	"golang.org/x/exp/shiny/widget"
	"golang.org/x/exp/shiny/widget/node"
	"golang.org/x/exp/shiny/widget/theme"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/geom"
	"golang.org/x/mobile/gl"
)

type State struct {
	tiles Tiles
	f     Field

	glctx gl.Context
	wsz   size.Event

	ful geom.Point
	fst geom.Point
	tsz image.Point

	buttons *widget.Uniform
	field   *widget.Uniform

	butImg *glutil.Image
	t      theme.Theme
}

func New() *State {
	s := &State{}
	s.f.Init()
	return s
}

type EvL struct{}
type EvR struct{}
type EvU struct{}
type EvD struct{}
type EvQ struct{}

const margin = 0.02

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

// helper.
func (s *State) ip2gp(ip image.Point) geom.Point {
	return geom.Point{
		X: geom.Pt(ip.X) / geom.Pt(s.wsz.PixelsPerPt),
		Y: geom.Pt(ip.Y) / geom.Pt(s.wsz.PixelsPerPt),
	}
}

func (s *State) setSize(e size.Event) {
	s.t = theme.Theme{DPI: float64(unit.PointsPerInch * e.PixelsPerPt)}

	horizontal := e.Orientation == size.OrientationLandscape

	padPx := e.WidthPx / 50
	bAx := widget.AxisHorizontal
	aAx := widget.AxisVertical
	if horizontal {
		padPx = e.HeightPx / 50
		bAx = widget.AxisVertical
		aAx = widget.AxisHorizontal
	}

	s.buttons = widget.NewUniform(theme.Neutral,
		widget.NewPadder(widget.AxisBoth, unit.Pixels(float64(padPx)),
			widget.NewFlow(bAx,
				widget.NewLabel("Foo"),
				stretch(widget.NewSpace(), 1),
				widget.NewLabel("Bar"),
			),
		),
	)
	s.field = widget.NewUniform(theme.Light, nil)
	all := widget.NewFlow(aAx,
		stretch(s.buttons, 0),
		stretch(widget.NewPadder(widget.AxisBoth, unit.Pixels(float64(padPx)), s.field), 1),
	)
	all.Measure(&s.t, e.WidthPx, e.HeightPx)
	all.Rect = image.Rectangle{Max: image.Pt(e.WidthPx, e.HeightPx)}
	all.Layout(&s.t)

	log.Print("Butt: ", s.buttons.Rect)
	log.Print("field: ", s.field.Rect)

	s.wsz = e

	r := s.field.Rect
	dx, dy := r.Dx(), r.Dy()
	if dx > dy {
		r.Min.X += dx - dy
	} else {
		r.Min.Y += dy - dx
	}

	s.tsz.X = r.Dx() / s.f.W()
	s.tsz.Y = r.Dy() / s.f.H()

	s.fst = s.ip2gp(s.tsz)
	s.ful = s.ip2gp(r.Min)

	s.tiles.SetSz(s.tsz)
	if s.butImg != nil {
		s.butImg.Release()
		s.butImg = nil
	}
}

func (s *State) drawButt() {
	if s.butImg == nil {
		s.butImg = s.tiles.ims.NewImage(s.buttons.Rect.Dx(), s.buttons.Rect.Dy())
		s.buttons.PaintBase(&node.PaintBaseContext{
			Theme: &s.t,
			Dst:   s.butImg.RGBA,
		}, image.Point{})
		s.butImg.Upload()
	}
	tl := s.buttons.Rect.Min
	br := s.buttons.Rect.Max
	tr := image.Pt(br.X, tl.Y)
	bl := image.Pt(tl.X, br.Y)
	s.butImg.Draw(s.wsz, s.ip2gp(tl), s.ip2gp(tr), s.ip2gp(bl), image.Rectangle{Max: s.buttons.Rect.Size()})
}

func (s *State) fRectBounds(x, y int) (geom.Point, geom.Point, geom.Point) {
	return geom.Point{s.ful.X + s.fst.X*geom.Pt(x), s.ful.Y + s.fst.X*geom.Pt(y)},
		geom.Point{s.ful.X + s.fst.X*geom.Pt(x+1), s.ful.Y + s.fst.X*geom.Pt(y)},
		geom.Point{s.ful.X + s.fst.X*geom.Pt(x), s.ful.Y + s.fst.X*geom.Pt(y+1)}
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
	case paint.Event:
	case size.Event:
		s.setSize(e)
	case error:
		log.Print(e)
		return true
	}
	s.draw(pub)
	return true
}
