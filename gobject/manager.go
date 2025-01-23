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

func (manager *Manager) MovingBullet(startX, startY, distance int32) {
	manager.BulletObj.X = startX + 32
	manager.BulletObj.Y = startY - 35
	manager.BulletObj.Draw(manager.r)
	for i := int32(1); i*100 <= distance; i++ {
		manager.BulletObj.UpMoving(manager.r, manager.enemies)
		sdl.Delay(10)
	}
}

func (manager *Manager) ShootUp(startX, startY, distance int32) {
	for key, obj := range manager.enemies {
		if startX+50 >= obj.X && startX+50 <= obj.X+120 && obj.Y-100 >= startY-distance {
			obj.IsMoving = false
			obj.Destroy(manager.r)
			delete(manager.enemies, key)
		}
	}
}

func (manager *Manager) ScanShoot() {
	if manager.PlayerObj.IsShoot {
		manager.MovingBullet(manager.PlayerObj.X, manager.PlayerObj.Y, int32(900))
		manager.PlayerObj.IsShoot = false
	}
}
