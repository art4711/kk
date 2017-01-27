// +build !android

package main

import (
	"kk/kk"
	"log"

	"golang.org/x/exp/shiny/driver/gldriver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/paint"
)

func keyFilter(ei interface{}) interface{} {
	if e, ok := ei.(key.Event); ok && e.Direction == key.DirPress {
		switch e.Code {
		case key.CodeLeftArrow:
			return kk.EvL{}
		case key.CodeRightArrow:
			return kk.EvR{}
		case key.CodeUpArrow:
			return kk.EvU{}
		case key.CodeDownArrow:
			return kk.EvD{}
		case key.CodeEscape:
			return kk.EvQ{}
		}
	}
	return ei
}

func main() {
	gldriver.Main(func(s screen.Screen) {
		st := kk.New()
		w, err := s.NewWindow(&screen.NewWindowOptions{1080, 1776})
		if err != nil {
			log.Fatal(err)
		}
		defer w.Release()
		for {
			e := w.NextEvent()
			repaint, quit, publish := st.Handle(keyFilter(e))
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
