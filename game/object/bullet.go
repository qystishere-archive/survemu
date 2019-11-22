package object

import "github.com/qystishere/survemu/game/utils"

type Bullet struct {
	PlayerID uint16
	Type     uint8

	Pos *utils.Vector2D
	Dir *utils.Vector2D

	Layer byte

	VarianceT float32
	DistAdjID byte

	ClipDistance bool
	Distance     float32

	ShotFX       uint8

	ShotOffHand byte
	LastShot    bool

	ReflectCount byte
	ReflectObjID uint16
}

func NewBullet(playerID uint16, t uint8, pos *utils.Vector2D, dir *utils.Vector2D, shotFX uint8) *Bullet {
	return &Bullet{
		PlayerID: playerID,
		Type:     t,

		Pos: pos,
		Dir: dir,

		Layer: 0,

		ShotFX:       shotFX,

		ShotOffHand: 1,
	}
}
