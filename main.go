package main

import (
	"kk/kk"
	"log"

	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/gl"
)

func main() {
	mmain()
}

func mmain() {
	app.Main(func(a app.App) {
		var glctx gl.Context

		repaint := false
		wsz := size.Event{}
		f := kk.NewField()

		var tiles *kk.Tiles

		for e := range a.Events() {
			switch e := a.Filter(e).(type) {
			case lifecycle.Event:
				switch e.Crosses(lifecycle.StageVisible) {
				case lifecycle.CrossOn:
					glctx, _ = e.DrawContext.(gl.Context)
					tiles = kk.NewTiles(glctx)
					a.Send(paint.Event{})
				case lifecycle.CrossOff:
					glctx = nil
					return
				}

			case key.Event:
				if e.Code == key.CodeEscape {
					return
				}
				if e.Direction == key.DirPress {
					switch e.Code {
					case key.CodeLeftArrow:
						f.Left()
					case key.CodeRightArrow:
						f.Right()
					case key.CodeUpArrow:
						f.Up()
					case key.CodeDownArrow:
						f.Down()
					}
					if !repaint {
						repaint = true
						a.Send(paint.Event{})
					}
				}

			case paint.Event:
				kk.Draw(glctx, wsz, f, tiles)
				a.Publish()
				repaint = false
			case size.Event:
				wsz = e

			case error:
				log.Print(e)
			}
		}
	})
}
