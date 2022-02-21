package util

type Vec2 struct {
	X float64
	Y float64
}

func (v *Vec2) Add(s *Vec2) Vec2 {
	return Vec2{v.X + s.X, v.Y + s.Y}
}