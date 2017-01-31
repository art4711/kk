// +build android

package main

import (
	"kk/kk"
	"math"

	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/touch"
)

type tf struct {
	touchStart touch.Event
}

func (t *tf) TouchFilter(ei interface{}) interface{} {
	if e, ok := ei.(touch.Event); ok {
		switch e.Type {
		case touch.TypeBegin:
			t.touchStart = e
		case touch.TypeEnd:
			x, y := e.X-t.touchStart.X, e.Y-t.touchStart.Y
			if math.Abs(float64(x)) > math.Abs(float64(y)) {
				if x < 0 {
					return kk.EvL{}
				} else {
					return kk.EvR{}
				}
			} else {
				if y < 0 {
					return kk.EvU{}
				} else {
					return kk.EvD{}
				}
			}
		}
	}
	return ei
}

func main() {
	app.Main(func(a app.App) {
		s := kk.New()
		tf := &tf{}
		a.RegisterFilter(tf.TouchFilter)

		for e := range a.Events() {
			if !s.Handle(a.Filter(e), func() { a.Publish() }) {
				return
			}
		}
	})
}
