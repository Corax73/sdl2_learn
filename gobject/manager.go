package gobject

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Manager struct {
	R         *sdl.Renderer
	PlayerObj *Gobject
	Enemies   map[string]*Gobject
	Bullets   map[string]*Gobject
}

func NewManager(player *Gobject, r *sdl.Renderer, enemies map[string]*Gobject, bullets map[string]*Gobject) *Manager {
	return &Manager{
		R:         r,
		PlayerObj: player,
		Enemies:   enemies,
		Bullets:   bullets,
	}
}

func (manager *Manager) ScanShoot() bool {
	if manager.PlayerObj.IsShoot {
		manager.PlayerObj.IsShoot = false
		return true
	}
	return false
}
