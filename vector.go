package main

// Vec2 position or direction in 2D board
type Vec2 struct {
	x, y int8
}

// Inv return the inverse vector
func (v Vec2) Inv() Vec2 {
	return Vec2{v.x, v.y}
}

// Add return the sum of two vectors
func (v Vec2) Add(u Vec2) Vec2 {
	return Vec2{v.x + u.x, v.y + u.y}
}

// Mul return a scaled vector
func (v Vec2) Mul(coef int8) Vec2 {
	return Vec2{v.x * coef, v.y * coef}
}

// Rot return a rotated vector
func (v Vec2) Rot(rot int8) Vec2 {
	rot = ((rot % 4) + 4) % 4
	switch rot {
	case 1:
		return Vec2{-v.y, v.x}
	case 2:
		return Vec2{-v.x, -v.y}
	case 3:
		return Vec2{v.y, -v.x}
	default:
		return Vec2{v.x, v.y}
	}
}

// ID return id based on position
func (v Vec2) ID() uint {
	return (uint(v.x) & 0xff) | ((uint(v.y) & 0xff) << 8)
}
