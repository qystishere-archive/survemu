package object

import "github.com/qystishere/survemu/game/utils"

type Emote struct {
	Type       uint8
	Pos        *utils.Vector2D
	UseLoadout bool
	IsPing     bool

	PlayerID uint16
}
