package kk

import (
	"log"

	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
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
		Draw(s.glctx, s.wsz, s.f, s.tiles)
		publish = true
	case size.Event:
		s.wsz = e

	case error:
		log.Print(e)
	}
	return
}
