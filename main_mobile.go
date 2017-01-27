// +build android

package main

import (
	"kk/kk"

	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/paint"
)

func main() {
	app.Main(func(a app.App) {
		s := kk.New()

		for e := range a.Events() {
			repaint, quit, publish := s.Handle(a.Filter(e))
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
