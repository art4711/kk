// +build android

package main

import (
	"kk/kk"

	"golang.org/x/mobile/app"
)

func main() {
	app.Main(func(a app.App) {
		s := kk.New()
		a.RegisterFilter(s.EvFilter)

		for e := range a.Events() {
			if !s.Handle(a.Filter(e), func() { a.Publish() }) {
				return
			}
		}
	})
}
