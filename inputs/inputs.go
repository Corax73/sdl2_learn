package inputs

import (
	"github.com/veandco/go-sdl2/sdl"
)

func Listen(isRunning bool) (bool, bool) {
	var isExit bool
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event.GetType() {
		case sdl.QUIT:
			return false, true
		case sdl.KEYDOWN:
			if 1 == sdl.GetKeyboardState()[sdl.SCANCODE_PAUSE] {
				return !isRunning, isExit
			} else {
				return isRunning, isExit
			}
		}
	}
	return isRunning, isExit
}
