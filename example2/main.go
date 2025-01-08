package example2

import (
	"fmt"
	"os"
	"sdl_learn/logger"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	defaultTitle = "test"
	winWidth     = 1280
	winHeight    = 720
	centerX      = winWidth / 2
	centerY      = winHeight / 2
)

func main() {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		logger.Error("sdl unable to init: %s", err.Error())
		os.Exit(1)
	}
	defer sdl.Quit()

	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")

	window, err := sdl.CreateWindow(defaultTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(winWidth), int32(winHeight), sdl.WINDOW_SHOWN)
	if err != nil {
		logger.Error("sdl unable to create window: %s", err.Error())
		os.Exit(1)
	}
	setWindowTitle(window)
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		logger.Error("sdl unable to create renderer: %s", err.Error())
		os.Exit(1)
	}
	defer renderer.Destroy()

	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(winWidth), int32(winHeight))
	if err != nil {
		logger.Error("sdl unable to create texture: %s", err.Error())
		os.Exit(1)
	}
	defer tex.Destroy()
}

func setWindowTitle(window *sdl.Window) {
	cursor := sdl.Point{}
	var keyboard []uint8
	quit := true
	for quit {
		for e := sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
			if e.GetType() == sdl.QUIT {
				quit = false
				window.Destroy()
			}
			switch event := e.(type) {
			case *sdl.KeyboardEvent:
				fmt.Printf("%d \n", event.Keysym)
			}

			switch e.GetType() {
			case sdl.MOUSEMOTION:
				motion := e.(*sdl.MouseMotionEvent)
				cursor.X = motion.X
				cursor.Y = motion.Y

			case sdl.KEYDOWN:
				fallthrough
			case sdl.KEYUP:
				keyboard = sdl.GetKeyboardState()
			}
			title := defaultTitle
			title += " | "

			title += time.Now().Format("2006-01-02 15:04:05")
			title += " | "
			title += "MX=" + fmt.Sprint(cursor.X)
			title += " | "
			title += "MY=" + fmt.Sprint(cursor.Y)
			title += " | "
			title += "UP=" + fmt.Sprint(keyboard)

			window.SetTitle(title)
		}
	}
}
