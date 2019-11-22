package utils

import "math"

type Vector2D struct {
	X float32
	Y float32
}

func (v *Vector2D) Distance(v2 *Vector2D) float64 {
	return math.Sqrt(math.Pow(float64(v2.X - v.X), 2) + math.Pow(float64(v2.Y - v.Y), 2))
}