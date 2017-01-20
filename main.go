package main

import (
	"kk/kk"
	"log"

	"golang.org/x/exp/shiny/driver/gldriver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/gl"
)

type state struct {
	tiles *kk.Tiles
	f     *kk.Field
	wsz   size.Event
	glctx gl.Context
}

func st() *state {
	s := state{f: kk.NewField()}
	return &s
}

func (s *state) handle(ei interface{}) (repaint bool, quit bool, publish bool) {
	switch e := ei.(type) {
	case lifecycle.Event:
		if e.To == lifecycle.StageDead {
			quit = true
			return
		}
		switch e.Crosses(lifecycle.StageVisible) {
		case lifecycle.CrossOn:
			s.glctx, _ = e.DrawContext.(gl.Context)
			s.tiles = kk.NewTiles(s.glctx)
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
		kk.Draw(s.glctx, s.wsz, s.f, s.tiles)
		publish = true
	case size.Event:
		s.wsz = e

	case error:
		log.Print(e)
	}
	return
}

func main() {
	smain()
}

func smain() {
	gldriver.Main(func(s screen.Screen) {
		st := st()
		w, err := s.NewWindow(&screen.NewWindowOptions{400, 400})
		if err != nil {
			log.Fatal(err)
		}
		defer w.Release()
		for {
			repaint, quit, publish := st.handle(w.NextEvent())
			if quit {
				return
			}
			if repaint {
				w.Send(paint.Event{})
			}
			if publish {
				w.Publish()
			}
		}
	})
}

/*
func mmain() {
	app.Main(func(a app.App) {
		s := st()

		for e := range a.Events() {
			repaint, quit, publish := s.handle(a.Filter(e))
			if quit {
				return
			}
			if repaint {
				a.Send(paint.Event{})
			}
			if publish {
				a.Publish()
			}
		}
	})
}
*/
