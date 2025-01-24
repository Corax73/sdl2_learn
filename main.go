package main

import (
	"fmt"
	"sdl_learn/gobject"
	"sdl_learn/inputs"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
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
	isExit    bool
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
	player := gobject.NewGobject(
		rend,
		"assets/battleship.png",
		"assets/exp.png",
		"assets/bullet.png",
		"player",
		WindowWidth/2-10,
		int32(float64(WindowHeight)*0.8),
		WindowWidth,
		WindowHeight,
		true,
	)
	ufo1 := gobject.NewGobject(
		rend,
		"assets/ufo.png",
		"assets/exp.png",
		"",
		"ufo1",
		WindowWidth/2-10,
		int32(10),
		WindowWidth,
		WindowHeight,
		true,
	)
	ufo2 := gobject.NewGobject(
		rend,
		"assets/ufo.png",
		"assets/exp.png",
		"",
		"ufo2",
		WindowWidth/2-200,
		int32(100),
		WindowWidth,
		WindowHeight,
		true,
	)
	ufo3 := gobject.NewGobject(
		rend,
		"assets/ufo.png",
		"assets/exp.png",
		"",
		"ufo3",
		WindowWidth/2+10,
		int32(200),
		WindowWidth,
		WindowHeight,
		true,
	)
	ufo4 := gobject.NewGobject(
		rend,
		"assets/ufo.png",
		"assets/exp.png",
		"",
		"ufo4",
		WindowWidth/2+200,
		int32(300),
		WindowWidth,
		WindowHeight,
		true,
	)

	bullet := gobject.NewGobject(
		rend,
		"assets/bullet.png",
		"",
		"",
		"bullet",
		WindowWidth/2+200,
		int32(300),
		WindowWidth,
		WindowHeight,
		true,
	)

	// Init gameObjects map
	players = make(map[string]*gobject.Gobject)
	players[player.Id] = player
	enemies := make(map[string]*gobject.Gobject)
	enemies[ufo1.Id] = ufo1
	enemies[ufo2.Id] = ufo2
	enemies[ufo3.Id] = ufo3
	enemies[ufo4.Id] = ufo4

	manager := gobject.NewManager(player, bullet, rend, enemies)

	var left int32
	var right int32

startGame:
	// Game loop
	for isRunning {
		frameStartTime := sdl.GetTicks()

		isRunning, isExit = inputs.Listen(isRunning)
		if !isRunning {
			goto paused
		}

		// Handle event queue
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				fmt.Println(t)
				isExit = true
				isRunning = false
				break
			}
		}

		// Clear screen
		rend.SetDrawColor(0, 100, 155, 0)
		rend.Clear()

		if player.IsMoving {
			player.Draw(rend)
			player.Update(rend)
			manager.ScanShoot()
		} else {
			player.Draw(rend)
		}

		if left < WindowWidth/4 {
			for _, val := range enemies {
				if val.IsMoving && val.X > 100 {
					val.Draw(rend)
					val.LeftMoving(rend, player)
					left++
				} else {
					val.Draw(rend)
				}
			}
		} else if left >= WindowWidth/4 && right < WindowWidth-100 {
			for _, val := range enemies {
				if val.IsMoving && val.X < WindowWidth-100 {
					val.Draw(rend)
					val.RightMoving(rend, player)
					right++
				} else {
					val.Draw(rend)
				}
			}
		} else {
			if right >= WindowWidth-100 && left >= WindowWidth/4 {
				left, right = -WindowWidth, 0
			}
		}
		rend.Present()

		// If too fast add delay
		frameTime := sdl.GetTicks() - frameStartTime
		if frameTime < DelayTime {
			sdl.Delay(uint32(DelayTime - frameTime))
		}

	} // End of isRunning

paused:
	for !isExit {
		isRunning, isExit = inputs.Listen(isRunning)
		if isExit {
			break
		}
		if isRunning {
			goto startGame
		}
		drawText(win, "Pause", rend)
	}
	if isExit {
		for _, val := range enemies {
			val.Free()
		}
		player.Free()
		sdl.Quit()
	}
}

func drawText(win *sdl.Window, drawingText string, rend *sdl.Renderer) {
	if drawingText != "" {
		var font *ttf.Font
		//var surface *sdl.Surface
		var text *sdl.Surface

		var textRect sdl.Rect
		var textImage *sdl.Texture
		textRect.X = WindowWidth/2
		textRect.Y = WindowHeight/2
		if err = ttf.Init(); err != nil {
			return
		}

		/*if surface, err = win.GetSurface(); err != nil {
			return
		}*/

		// Load the font for our text
		if font, err = ttf.OpenFont("assets/test.ttf", 48); err != nil {
			return
		}
		defer font.Close()

		// Create text with the font
		if text, err = font.RenderUTF8Blended(drawingText, sdl.Color{R: 155, G: 0, B: 100, A: 255}); err != nil {
			return
		}
		defer text.Free()

		textRect.W = text.W
		textRect.H = text.H

		if textImage, err = rend.CreateTextureFromSurface(text); err != nil {
			fmt.Println(err)
		}

		rend.Copy(textImage, nil, &textRect)
		fmt.Println(textImage)
		// Draw the text around the center of the window
		/*if err = text.Blit(nil, surface, &sdl.Rect{X: WindowWidth/2 - 50, Y: WindowHeight/2 - 50, W: 0, H: 0}); err != nil {
			return
		}*/

		// Update the window surface with what we have drawn
		win.UpdateSurface()
	}
}
