package gobject

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Manager struct {
	r         *sdl.Renderer
	PlayerObj *Gobject
	BulletObj *Gobject
	enemies   map[string]*Gobject
}

func NewManager(player *Gobject, bullet *Gobject, r *sdl.Renderer, enemies map[string]*Gobject) *Manager {
	return &Manager{
		r:         r,
		PlayerObj: player,
		BulletObj: bullet,
		enemies:   enemies,
	}
}

func (manger *Manager) MovingBullet(startX, startY, distance int32) {
	manger.BulletObj.X = startX + 32
	manger.BulletObj.Y = startY - 35
	manger.BulletObj.Draw(manger.r)
	for i := int32(1); i*100 <= distance; i++ {
		manger.BulletObj.UpMoving(manger.r)
		manger.BulletObj.Draw(manger.r)
		sdl.Delay(10)
	}
	manger.ShootUp(startX, startY, distance)
}

func (manger *Manager) ShootUp(startX, startY, distance int32) {
	for key, obj := range manger.enemies {
		if startX+50 >= obj.X && startX+50 <= obj.X+120 && obj.Y-100 >= startY-distance {
			obj.IsMoving = false
			obj.Destroy(manger.r)
			delete(manger.enemies, key)
		}
	}
}

func (manger *Manager) ScanShoot() {
	if manger.PlayerObj.IsShoot {
		manger.MovingBullet(manger.PlayerObj.X, manger.PlayerObj.Y, int32(900))
		manger.PlayerObj.IsShoot = false
	}
}
