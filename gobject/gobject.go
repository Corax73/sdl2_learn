package gobject

import (
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
	fromX, fromY, Width, Height int32
	// Holds image
	Texture *sdl.Texture
	// Part of the spritesheet
	Src sdl.Rect
	// Part of the screen where to draw
	Dest sdl.Rect
	// Is object moving
	IsMoving bool
	Direction sdl.FPoint
}

// NewGobject creates new game object
func NewGobject(r *sdl.Renderer, file, id string, x, y int32) *Gobject {
	gob := &Gobject{Filename: file, X: x, Y: y, Speed: 1}
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
func (gob *Gobject) Update() {
	keyStates := sdl.GetKeyboardState()
	gob.Speed = 20

	if keyStates[sdl.SCANCODE_LEFT] == 1 {
		gob.fromX = gob.X
		gob.Direction.X = -1
		gob.X -= gob.Speed
		sdl.Delay(50)
	} else if keyStates[sdl.SCANCODE_RIGHT] == 1 {
		gob.fromX = gob.X
		gob.Direction.X = 1
		gob.X += gob.Speed
		sdl.Delay(50)
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
