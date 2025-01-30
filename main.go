package main

import (
	"fmt"
	"sdl_learn/gobject"
	"sdl_learn/inputs"
	"strconv"

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
	isExit    bool
	// Maybe later
	players, enemies, bullets map[string]*gobject.Gobject
	BulletId                  int
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

	// Init gameObjects map
	players = make(map[string]*gobject.Gobject)
	players[player.Id] = player
	enemies := make(map[string]*gobject.Gobject)
	enemies[ufo1.Id] = ufo1
	enemies[ufo2.Id] = ufo2
	enemies[ufo3.Id] = ufo3
	enemies[ufo4.Id] = ufo4

	bullets := make(map[string]*gobject.Gobject)

	manager := gobject.NewManager(player, rend, enemies, bullets)

	var left int32
	var right int32
	var shot bool

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
			shot = manager.ScanShoot()
		} else {
			player.Draw(rend)
		}

		if shot {
			bullet := NewBullet(rend, player.X, player.Y)
			manager.Bullets[bullet.Id] = bullet
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

		for key, val := range manager.Bullets {
			if val.Y > 0 && val.IsMoving {
				val.UpMoving(manager.R, manager.Enemies, manager.Bullets, manager.PlayerObj)
				val.Draw(rend)
			} else if val.Y <= 0 || !val.IsMoving {
				val.Draw(rend)
				val.Free()
				delete(manager.Bullets, key)
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
		isRunning = drawText()
	}
	if isExit {
		for _, val := range enemies {
			val.Free()
		}
		player.Free()
		sdl.Quit()
	}
}

func drawText() bool {
	buttons := []sdl.MessageBoxButtonData{
		{Flags: sdl.MESSAGEBOX_BUTTON_ESCAPEKEY_DEFAULT, ButtonID: 0, Text: "cancel PAUSE"},
	}

	colorScheme := sdl.MessageBoxColorScheme{
		Colors: [5]sdl.MessageBoxColor{
			sdl.MessageBoxColor{R: 255, G: 0, B: 0},
			sdl.MessageBoxColor{R: 0, G: 255, B: 0},
			sdl.MessageBoxColor{R: 255, G: 255, B: 0},
			sdl.MessageBoxColor{R: 0, G: 0, B: 255},
			sdl.MessageBoxColor{R: 255, G: 0, B: 255},
		},
	}

	messageboxdata := sdl.MessageBoxData{
		Flags:       sdl.MESSAGEBOX_INFORMATION,
		Window:      nil,
		Title:       "PAUSED",
		Message:     "",
		Buttons:     buttons,
		ColorScheme: &colorScheme,
	}

	buttonid, _ := sdl.ShowMessageBox(&messageboxdata)

	if buttonid == 1 {
		return false
	}
	return true
}

func NewBullet(r *sdl.Renderer, x, y int32) *gobject.Gobject {
	BulletId += 1
	strId := "bullet" + strconv.Itoa(BulletId)
	return gobject.NewGobject(
		r,
		"assets/bullet.png",
		"",
		"",
		strId,
		x,
		y-35,
		WindowWidth,
		WindowHeight,
		true,
	)
}
