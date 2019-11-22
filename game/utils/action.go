package utils

import "time"

type Action struct {
	Time     float32
	Duration float32
	Type     uint8

	CreatedAt time.Time
	TargetID  uint16
}
