package main

import (
	"fmt"
	"sdl_learn/gobject"

	"github.com/veandco/go-sdl2/sdl"
)

// Global consts
const (
	FPS          uint32 = 60
	DelayTime    uint32 = uint32(1000.0 / FPS)
	WindowWidth  int32  = 1280
	WindowHeight int32  = 720
	WindowTitle         = "Game"
)

// Globals, maybe someday wrapped to struct but now less to type
var (
	win       *sdl.Window
	rend      *sdl.Renderer
	event     sdl.Event
	err       error
	isRunning = true
	// Maybe later
	players map[string]*gobject.Gobject
)

// Error checker
func perror(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	// Init SDL and create window
	err = sdl.Init(sdl.INIT_VIDEO)
	perror(err)

	win, err = sdl.CreateWindow(
		WindowTitle,
		sdl.WINDOWPOS_CENTERED,
		sdl.WINDOWPOS_CENTERED,
		WindowWidth,
		WindowHeight,
		sdl.WINDOW_SHOWN,
	)
	perror(err)
	defer win.Destroy()

	// Create renderer
	rend, err = sdl.CreateRenderer(win, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	perror(err)
	defer rend.Destroy()

	// Create player
	player := gobject.NewGobject(rend, "assets/battleship.png", "assets/exp.png", "player", WindowWidth/2-10, int32(float64(WindowHeight)*0.8), WindowWidth, WindowHeight, true)
	ufo1 := gobject.NewGobject(rend, "assets/ufo.png", "assets/exp.png", "ufo1", WindowWidth/2-10, int32(10), WindowWidth, WindowHeight, true)
	ufo2 := gobject.NewGobject(rend, "assets/ufo.png", "assets/exp.png", "ufo2", WindowWidth/2-200, int32(100), WindowWidth, WindowHeight, true)
	// Init gameObjects map
	players = make(map[string]*gobject.Gobject)
	players[player.Id] = player
	enemies := make(map[string]*gobject.Gobject)
	enemies[ufo1.Id] = ufo1
	enemies[ufo2.Id] = ufo2

	// Game loop
	for isRunning {
		frameStartTime := sdl.GetTicks()

		// Handle event queue
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				fmt.Println(t)
				isRunning = false
			}
		}

		// Clear screen
		rend.SetDrawColor(0, 100, 155, 0)
		rend.Clear()

		if player.IsMoving {
			player.Draw(rend)
			player.Update(rend, enemies)
		} else {
			player.Draw(rend)
		}

		for _, val := range enemies {
			if val.IsMoving {
				val.Draw(rend)
				val.RandomMoving(rend, player)
			} else {
				val.Draw(rend)
			}
		}
		rend.Present()

		// If too fast add delay
		frameTime := sdl.GetTicks() - frameStartTime
		if frameTime < DelayTime {
			sdl.Delay(uint32(DelayTime - frameTime))
		}

	} // End of isRunning

	for _, val := range enemies {
		val.Free()
	}
	player.Free()
	sdl.Quit()
}
