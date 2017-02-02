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
intended way to use the library. Not that I actually need OpenGL, but
let's use the powerful drugs from start.

The second way the shiny library is abused is that I couldn't be
bothered to write a generic box layout library so I use shiny to
create and layout a bunch of widgets that only serve to provide me
rectangles on the screen. This is also probably not an intended use of
shiny either. Interestingly enough there doesn't seem to be a way yet
to use shiny to route events to the right boxes, something I actually
wanted to use it for (it seems to be done internally, but I haven't
found any way that was exposed to the user).

 * [kk/field.go](kk/field.go) implements the game logic and scores.

 * [kk/tiles.go](kk/tiles.go) implements the bitmaps/sprites.

 * [kk/layout.go](kk/layout.go) implements the layout of the ui elements.

 * [kk/state.go](kk/state.go) implements the main event handler and drawing.

## TODO ##

Use unpredictable for randomness? I refuse to use the normal monkey
methods to seed a terrible RNG so right now the randomness is entirely
unseeded and generates the exact same game every time.

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

Save save state to something permanent. Gomobile doesn't provide an
abstraction for any permanent storage and I don't think it's worth it
to write specific Android code to plug into storage just to stuff a
handful of bytes somewhere.
