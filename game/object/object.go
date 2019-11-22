package object

import "github.com/qystishere/survemu/game/utils"

type Object struct {
	ID uint16

	Type    uint8

	ItemType uint16

	Pos *utils.Vector2D
	Ori int32

	Count uint8

	Layer byte
	Scale float32

	Old byte
}

func NewObject(id uint16, t uint8, itemType uint16, pos *utils.Vector2D, ori int32, count uint8, old byte) *Object {
	return &Object{
		ID:       id,
		Type:     t,
		ItemType: itemType,
		Pos: pos,
		Ori: ori,
		Count: count,
		Layer: 0,
		Scale: .5,
		Old: old,
	}
}