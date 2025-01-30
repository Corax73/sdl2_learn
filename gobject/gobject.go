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
	Filename, FilenameDestruction, FilenameBullet string
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
	// Holds image
	TextureDestruction *sdl.Texture
	// Holds image
	TextureBullet *sdl.Texture
	// Part of the spritesheet
	Src sdl.Rect
	// Part of the screen where to draw
	Dest sdl.Rect
	// Is object moving
	IsMoving          bool
	IsShoot           bool
	Direction         sdl.FPoint
	Score, ShootDelay int
}

// NewGobject creates new game object
func NewGobject(r *sdl.Renderer, file, filenameDestruction, filenameBullet, id string, x, y, maxX, maxY int32, isMoving bool) *Gobject {
	gob := &Gobject{
		Filename:            file,
		FilenameDestruction: filenameDestruction,
		FilenameBullet:      filenameBullet,
		Id:                  id,
		X:                   x,
		Y:                   y,
		MaxX:                maxX,
		MaxY:                maxY,
		Speed:               1,
		IsMoving:            isMoving,
	}
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

	if gob.FilenameDestruction != "" {
		imageDestruction, err := img.Load(gob.FilenameDestruction)
		if err != nil {
			panic(err)
		}
		defer imageDestruction.Free()
		gob.TextureDestruction, err = r.CreateTextureFromSurface(imageDestruction)
		if err != nil {
			panic(err)
		}
	}

	if gob.FilenameBullet != "" {
		imageBullet, err := img.Load(gob.FilenameBullet)
		if err != nil {
			panic(err)
		}
		defer imageBullet.Free()
		gob.TextureBullet, err = r.CreateTextureFromSurface(imageBullet)
		if err != nil {
			panic(err)
		}
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
	if gob.IsMoving {
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
			if gob.ShootDelay == 0 {
				gob.IsShoot = true
				gob.ShootDelay = 3
			} else {
				gob.ShootDelay -= 1
			}
		}
	}
}

func (gob *Gobject) Rect() sdl.Rect {
	x, y := int32(gob.X), int32(gob.Y)
	_, _, imageWidth, imageHeight, _ := gob.Texture.Query()
	return sdl.Rect{
		X: x,
		Y: y,
		W: imageWidth,
		H: imageHeight,
	}
}

// Draw object
func (gob *Gobject) Draw(r *sdl.Renderer) {
	dst := gob.Rect()
	if gob.IsMoving {
		r.Copy(gob.Texture, nil, &dst)
	} else {
		r.Copy(gob.TextureDestruction, nil, &dst)
	}
}

func (gob *Gobject) ShootDown(r *sdl.Renderer, player *Gobject) {
	if gob.IsMoving {
		r.SetDrawColor(0, 255, 0, 0)
		r.DrawLine(gob.X+60, gob.Y+100, gob.X+60, gob.Y+300)
		if gob.X+60 >= player.X && gob.X+60 <= player.X+50 && player.Y <= gob.Y+300 {
			player.IsMoving = false
			player.Destroy(r)
		}
	}
}

func (gob *Gobject) RandomMoving(r *sdl.Renderer, player *Gobject) {
	ctx, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		var min int32 = 20
		gob.Speed = min
		if gob.IsMoving {
			select {
			case <-ctx.Done():
				return
			default:
				for i := 0; i < 2; i++ {
					if gob.IsMoving {
						sdl.Delay(1500)
						maxRand := big.NewInt(4)
						if val, err := rand.Int(rand.Reader, maxRand); err == nil && val.Int64() > 2 && gob.IsMoving {
							gob.ShootDown(r, player)
							sdl.Delay(1500)
							if (gob.X-gob.Speed) > 0 && gob.IsMoving {
								gob.X -= gob.Speed
								sdl.Delay(1500)
							}
							if (gob.Y+gob.Speed+100) < gob.MaxY && gob.IsMoving {
								gob.Y += gob.Speed
								sdl.Delay(1500)
							}
							if (gob.X+gob.Speed+200) < gob.MaxX && gob.IsMoving {
								gob.X += gob.Speed
								sdl.Delay(1500)
							}
							if (gob.Y-gob.Speed) > 0 && gob.IsMoving {
								gob.Y -= gob.Speed
								sdl.Delay(1500)
							}
						} else {
							sdl.Delay(1500)
							if (gob.X+gob.Speed+200) < gob.MaxX && gob.IsMoving {
								gob.X += gob.Speed
								sdl.Delay(1500)
							}
							if (gob.Y-gob.Speed) > 0 && gob.IsMoving {
								gob.Y -= gob.Speed
								sdl.Delay(1500)
							}
							if (gob.X-gob.Speed) > 0 && gob.IsMoving {
								gob.X -= gob.Speed
								sdl.Delay(1500)
							}
							if (gob.Y+gob.Speed+100) < gob.MaxY && gob.IsMoving {
								gob.Y += gob.Speed
								sdl.Delay(1500)
							}
						}
						maxSpeed := big.NewInt(41)
						val, err := rand.Int(rand.Reader, maxSpeed)
						if err == nil {
							gob.Speed = int32(val.Int64()) + min
						}
					} else {
						gob.Destroy(r)
					}
				}
			}
		}
	}(ctx)
	go func() {
		sdl.Delay(1500)
		cancel()
	}()
}

func (gob *Gobject) Destroy(r *sdl.Renderer) {
	if !gob.IsMoving {
		dst := gob.Rect()
		r.Copy(gob.TextureDestruction, nil, &dst)
		sdl.Delay(500)
		gob.Free()
	}
}

func (gob *Gobject) LeftMoving(r *sdl.Renderer, player *Gobject) {
	ctx, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		gob.Speed = 2
		if gob.IsMoving {
			select {
			case <-ctx.Done():
				return
			default:
				if gob.IsMoving {
					sdl.Delay(1500)
					maxRand := big.NewInt(4)
					val, err := rand.Int(rand.Reader, maxRand)
					if err == nil && val.Int64() > 2 {
						gob.ShootDown(r, player)
					}
					sdl.Delay(1500)
					if (gob.X-gob.Speed) > 100 && gob.IsMoving {
						gob.X -= gob.Speed
						sdl.Delay(1500)
					}
				} else {
					gob.Destroy(r)
				}
			}
		}
	}(ctx)
	go func() {
		sdl.Delay(1500)
		cancel()
	}()
}

func (gob *Gobject) RightMoving(r *sdl.Renderer, player *Gobject) {
	ctx, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		gob.Speed = 2
		if gob.IsMoving {
			select {
			case <-ctx.Done():
				return
			default:
				if gob.IsMoving {
					sdl.Delay(1500)
					maxRand := big.NewInt(4)
					val, err := rand.Int(rand.Reader, maxRand)
					if err == nil && val.Int64() > 2 {
						gob.ShootDown(r, player)
					}
					sdl.Delay(1500)
					if (gob.X+gob.Speed) < 1100 && gob.IsMoving {
						gob.X += gob.Speed
						sdl.Delay(1500)
					}
				} else {
					gob.Destroy(r)
				}
			}
		}
	}(ctx)
	go func() {
		sdl.Delay(1500)
		cancel()
	}()
}

func (gob *Gobject) GetBulletRect(startX, startY int32) sdl.Rect {
	x, y := startX+32, startY-35
	_, _, imageWidth, imageHeight, _ := gob.TextureBullet.Query()
	return sdl.Rect{
		X: x,
		Y: y,
		W: imageWidth,
		H: imageHeight,
	}
}

func (gob *Gobject) UpMoving(r *sdl.Renderer, objects map[string]*Gobject, bullets map[string]*Gobject, player *Gobject) {
	ctx, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		gob.Speed = 1
		if gob.IsMoving {
			select {
			case <-ctx.Done():
				return
			default:
				if gob.Y > 0 && gob.IsMoving {
					sdl.Delay(100)
					gob.Y -= gob.Speed * 70
					for key, obj := range objects {
						if gob.IsMoving && !(gob.X >= obj.X+obj.Rect().W ||
							gob.X+gob.Rect().W <= obj.X ||
							gob.Y >= obj.Y+obj.Rect().H ||
							gob.Y+gob.Rect().H <= obj.Y) {
							obj.IsMoving = false
							gob.IsMoving = false
							obj.Destroy(r)
							player.Score += 100
							delete(objects, key)
						}
					}
					sdl.Delay(100)
				}
			}
		}
	}(ctx)
	go func() {
		sdl.Delay(500)
		cancel()
	}()
}
