package object

import (
	"time"

	"github.com/qystishere/survemu/game/data"
	"github.com/qystishere/survemu/game/utils"
)

type Player struct {
	ID   uint16
	Name string

	Team *Team

	Pos   *utils.Vector2D
	Ori   *utils.Vector2D
	Speed float32

	Dead   bool
	Downed bool

	AnimType byte
	AnimSeq  byte

	ActionType byte

	Emotes [data.ESCount]byte

	CurrentWSlotID byte
	WSlots         [data.WSCount]*utils.Weapon
	Items          [data.BSCount]uint16

	Scope    byte
	Skin     byte
	Backpack byte
	Helmet   byte
	Chest    byte

	ScopedIn bool
	Health   float32
	Boost    float32
	Action   *utils.Action

	WearingPan bool
	Frozen     bool
	FrozenOri  int32

	Layer byte

	Bullets []*Bullet

	LoadedPlayers map[uint16]*Player
	LoadedObjects map[uint16]*Object

	RequirePartUpdate bool
	RequireFullUpdate bool
}

func NewPlayer(name string, emotes [7]byte) *Player {
	return &Player{
		Name: name,

		Pos: &utils.Vector2D{X: 100, Y: 100},
		Ori: &utils.Vector2D{X: 100, Y: 100},

		Speed: 0.4,

		Emotes: emotes,

		CurrentWSlotID: 2,
		WSlots: [4]*utils.Weapon{
			{
				ID:   88,
				Ammo: 50,
			},
			{
				ID:   70,
				Ammo: 50,
			},
			{
				ID:   49,
				Ammo: 0,

				Duration: time.Millisecond * 250,
			},
			{
				ID:   114,
				Ammo: 50,
			},
		},
		Items: [21]uint16{
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0,
		},

		Scope:    141,
		Skin:     1,
		Backpack: 131,

		Health: 100,
		Boost:  0,

		LoadedPlayers: map[uint16]*Player{},
		LoadedObjects: map[uint16]*Object{},
	}
}
