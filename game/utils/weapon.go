package utils

import "time"

type Weapon struct {
	ID   uint8
	Ammo uint8

	LastUsedAt time.Time
	Duration   time.Duration
}
