# kk #

Experiment to see how hard it is to implement something for gomobile.
This is pretty much a clone of the game 2048, except that I probably
don't count score the same way and there's a save and load button.

Called `kk` because... two k... 2k... 2048... Get it? Anyway.

## Run ##

It's runnable on a desktop with `go run main_shiny.go`. It's
installable on android at least with `gomobile install kk`. iOS not
tested, but there's nothing in theory that should prevent it from
working.

## Implementation details ##

The shiny library is abused in creative ways. First, the desktop
version uses shiny to create a window and feed us the events, but we
bypass everything else it does and pretend the events are just like
mobile events. This is because the mobile stuff doesn't allow resizing
windows on desktop and the shiny OpenGL widget doesn't work on OSX.
This happens to work for now, but I understand that it isn't the
intended way to use the library.

The second way the shiny library is abused is that I couldn't be
bothered to write a generic box layout library so I use shiny to
create and layout a bunch of widgets that only serve to provide me
rectangles on the screen. This is also probably not an intended use of
shiny either. Interestingly enough there doesn't seem to be a way yet
to use shiny to route events to the right boxes, something I actually
wanted to use it for (it seems to be done internally, but I haven't
found any way that was exposed to the user).

 * [kk/field.go](kk/field.go) - the game logic and scores.

 * [kk/tiles.go](kk/tiles.go) - generation of bitmaps/sprites.

 * [kk/layout.go](kk/layout.go) - layout of the ui elements.

 * [kk/state.go](kk/state.go) - main event handler and drawing.

 * [kk/persistent.go](kk/persistent.go) - persistent saving of state.

I'm fully aware that all this is silly. I don't need OpenGL for
anything, saving is pretty useless, saving the state persistently is
meaningless, I could have almost hardcoded the layout with some
trivial code, I didn't need to overgeneralize the event handlers so
much, etc. The point is to run into as many small problems as possible
to see how things are done.

## TODO ##

Get rid of glutil.Image. It's three layers of indirection and unit
obfuscation for something that will be shorter to just layout
ourselves without all this complexity of weird units. Besides, we'll
probably need our own shaders to draw the ui without going through
bitmaps all the time.

Maybe implement the relevant bits of shiny widgets without all the
baggage? If they can't be used on mobile, what's the point? This will
probably be relevant once shiny decides to break or unexpose the
interfaces that happen to be exposed to us.

Unbreak the paint events. We probably don't want to repaint for every
bloody resize event (there's lots of them in a window).

Better font for the score.

Acquire taste before ever picking colors for anything ever. Especially
the game over color is eye-watering.

It would be nice to have a portable way to figure out where to save
permanent state. Gomobile doesn't provide an abstraction for any
permanent storage, so it's up to us to guess paths. Maybe there is a
need for a package that implements one function:
`DirectoryForSmallPersistentState`.
