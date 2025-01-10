package gobject

import (
	"context"
	"crypto/rand"
	"math/big"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

// GameObject interface
type GameObject interface {
	Draw(r *sdl.Renderer)
	Update()
}

// Gobject represents game object
type Gobject struct {
	// Image filename
	Filename string
	// Key for mapping
	Id string
	// Position
	X, Y int32
	// Movement speed
	Speed int32
	// Sprite size, not full image size
	MaxX, MaxY, Width, Height int32
	// Holds image
	Texture *sdl.Texture
	// Part of the spritesheet
	Src sdl.Rect
	// Part of the screen where to draw
	Dest sdl.Rect
	// Is object moving
	IsMoving  bool
	Direction sdl.FPoint
}

// NewGobject creates new game object
func NewGobject(r *sdl.Renderer, file, id string, x, y, maxX, maxY int32) *Gobject {
	gob := &Gobject{Filename: file, X: x, Y: y, MaxX: maxX, MaxY: maxY, Speed: 1}
	gob.Load(r)
	return gob
}

// Load texture
func (gob *Gobject) Load(r *sdl.Renderer) {
	image, err := img.Load(gob.Filename)
	if err != nil {
		panic(err)
	}
	defer image.Free()

	gob.Texture, err = r.CreateTextureFromSurface(image)
	if err != nil {
		panic(err)
	}

	// Query image size and calculate frame width and height
	_, _, imageWidth, imageHeight, _ := gob.Texture.Query()
	gob.Width = imageWidth
	gob.Height = imageHeight
}

// Free resources
func (gob *Gobject) Free() {
	gob.Texture.Destroy()
}

// Update updates object state
func (gob *Gobject) Update(r *sdl.Renderer) {
	keyStates := sdl.GetKeyboardState()
	gob.Speed = 20
	if keyStates[sdl.SCANCODE_LEFT] == 1 {
		if (gob.X - gob.Speed) > 0 {
			gob.X -= gob.Speed
		}
		sdl.Delay(50)
	} else if keyStates[sdl.SCANCODE_RIGHT] == 1 {
		if (gob.X + gob.Speed + 70) < gob.MaxX {
			gob.X += gob.Speed
		}
		sdl.Delay(50)
	}
	if keyStates[sdl.SCANCODE_SPACE] == 1 {
		gob.ShootUp(r)
	}
}

func (gob *Gobject) Rect() sdl.Rect {
	x, y := int32(gob.X), int32(gob.Y)
	return sdl.Rect{
		X: x,
		Y: y,
		W: gob.Width,
		H: gob.Height,
	}
}

// Draw object
func (gob *Gobject) Draw(r *sdl.Renderer) {
	dst := gob.Rect()
	r.Copy(gob.Texture, nil, &dst)
}

func (gob *Gobject) ShootUp(r *sdl.Renderer) {
	r.SetDrawColor(255, 255, 0, 255)
	r.DrawLine(gob.X+50, gob.Y-5, gob.X+50, gob.Y-300)
}

func (gob *Gobject) ShootDown(r *sdl.Renderer) {
	r.SetDrawColor(0, 255, 0, 0)
	r.DrawLine(gob.X+50, gob.Y-20, gob.X+50, gob.Y+300)
}

func (gob *Gobject) RandomMoving(r *sdl.Renderer) {
	ctx, cancel := context.WithCancel(context.Background())

	go func(ctx context.Context) {
		var min int32 = 20
		gob.Speed = min
		select {
		case <-ctx.Done():
			return
		default:
			for i := 0; i < 2; i++ {
				sdl.Delay(1500)
				maxRand := big.NewInt(4)
				if val, err := rand.Int(rand.Reader, maxRand); err == nil && val.Int64() > 2 {
					gob.ShootDown(r)
					sdl.Delay(1500)
					if (gob.X - gob.Speed) > 0 {
						gob.X -= gob.Speed
						sdl.Delay(1500)
					}
					if (gob.Y + gob.Speed + 100) < gob.MaxY {
						gob.Y += gob.Speed
						sdl.Delay(1500)
					}
					if (gob.X + gob.Speed + 200) < gob.MaxX {
						gob.X += gob.Speed
						sdl.Delay(1500)
					}
					if (gob.Y - gob.Speed) > 0 {
						gob.Y -= gob.Speed
						sdl.Delay(1500)
					}
				} else {
					if (gob.X + gob.Speed + 200) < gob.MaxX {
						gob.X += gob.Speed
						sdl.Delay(1500)
					}
					if (gob.Y - gob.Speed) > 0 {
						gob.Y -= gob.Speed
						sdl.Delay(1500)
					}
					if (gob.X - gob.Speed) > 0 {
						gob.X -= gob.Speed
						sdl.Delay(1500)
					}
					if (gob.Y + gob.Speed + 100) < gob.MaxY {
						gob.Y += gob.Speed
						sdl.Delay(1500)
					}
				}
				maxSpeed := big.NewInt(41)
				val, err := rand.Int(rand.Reader, maxSpeed)
				if err == nil {
					gob.Speed = int32(val.Int64()) + min
				}
			}
		}
	}(ctx)
	go func() {
		sdl.Delay(1500)
		cancel()
	}()

}
