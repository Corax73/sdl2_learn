package main

// #cgo windows linux freebsd darwin pkg-config: sdl2

import (
	"math"
	"math/rand"
	"time"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	winWidth      = 1280
	winHeight     = 720
	minWarpFactor = 0.1
	numStars      = 300
	centerX       = winWidth / 2
	centerY       = winHeight / 2
)

type position struct {
	x float64
	y float64
}

type star struct {
	pos        position
	vel        position
	brightness byte
}

type stars struct {
	stars []star
}

func randFloat64(min float64, max float64) float64 {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Float64()*(max-min)
}

// clear a slice of pixels.
func clear(pixels []byte) {
	for i := range pixels {
		pixels[i] = 0
	}
}

func setPixel(x, y int, c byte, pixels []byte) {
	index := (y*winWidth + x) * 4

	if index < len(pixels)-4 && index >= 0 {
		pixels[index] = c
		pixels[index+1] = c
		pixels[index+2] = c
	}
}

func newStar() star {

	// # Pick a direction and speed
	// angle = random.uniform(-math.pi, math.pi)
	angle := randFloat64(float64(-3.14), float64(3.14))

	// speed = 255 * random.uniform(0.3, 1.0) ** 2
	speed := 255 * math.Pow(randFloat64(float64(0.3), float64(1.0)), 2)

	// # Turn the direction into position and velocity vectors
	// dx = math.cos(angle)
	dx := math.Cos(angle)

	// dy = math.sin(angle)
	dy := math.Sin(angle)

	// d = random.uniform(25 + TRAIL_LENGTH, 100)
	d := rand.Intn(100) + 25 //+ traillength

	// pos = centerx + dx * d, centery + dy * d
	pos := position{
		x: centerX + dx*float64(d),
		y: centerY + dy*float64(d),
	}

	// vel = speed * dx, speed * dy
	vel := position{
		x: speed * dx,
		y: speed * dy,
	}

	s := star{
		pos:        pos,
		vel:        vel,
		brightness: 0,
	}

	return s
}

func (s *stars) update(elapsedTime float32) {

	// calculate the stars new position
	for i := 0; i < len(s.stars); i++ {
		newPosX := s.stars[i].pos.x + (s.stars[i].vel.x * minWarpFactor) //* dt
		newPosY := s.stars[i].pos.y + (s.stars[i].vel.y * minWarpFactor) //* dt

		// if we're off the screen with the new position reset else update position
		if newPosX > winWidth || newPosY > winHeight || newPosX < 0 || newPosY < 0 {
			s.stars[i] = newStar()
		} else {
			s.stars[i].pos.x = newPosX
			s.stars[i].pos.y = newPosY

			// # Grow brighter
			// s.brightness = min(s.brightness + warp_factor * 200 * dt, s.speed)
			if s.stars[i].brightness < 255 {
				s.stars[i].brightness += 40
			}
		}
	}
}

func (s *stars) draw(pixels []byte) {
	for i := 0; i < len(s.stars); i++ {
		if int(s.stars[i].pos.x) >= 0 {
			setPixel(int(s.stars[i].pos.x), int(s.stars[i].pos.y), s.stars[i].brightness, pixels)
		}
	}
}

func main() {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}
	defer sdl.Quit()

	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")

	window, err := sdl.CreateWindow("Stars", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(winWidth), int32(winHeight), sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()

	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(winWidth), int32(winHeight))
	if err != nil {
		panic(err)
	}
	defer tex.Destroy()

	var elapsedTime float32
	pixels := make([]byte, winWidth*winHeight*4)
	starField := make([]star, numStars)
	all := &stars{}

	for i := 0; i < len(starField); i++ {
		all.stars = append(all.stars, newStar())
	}

	for {
		frameStart := time.Now()
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}

		all.update(elapsedTime)
		all.draw(pixels)

		tex.Update(nil, unsafe.Pointer(&pixels[0]), winWidth*4)
		renderer.Copy(tex, nil, nil)
		renderer.Present()
		clear(pixels)
		elapsedTime = float32(time.Since(frameStart).Seconds() * 1000)
		if elapsedTime < 7 {
			sdl.Delay(7 - uint32(elapsedTime))
			elapsedTime = float32(time.Since(frameStart).Seconds() * 1000)
		}
	}
}
