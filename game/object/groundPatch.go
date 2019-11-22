package object

import "github.com/qystishere/survemu/game/utils"

type GroundPatch struct {
	Min        *utils.Vector2D
	Max        *utils.Vector2D
	Color      uint32
	Roughness  float32
	OffsetDist float32
}
