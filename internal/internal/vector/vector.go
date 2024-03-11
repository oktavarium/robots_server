package vector

type Vector struct {
	X float64
	Y float64
	Z float64
}

func NewVector(x, y, z float64) Vector {
	return Vector{
		X: x,
		Y: y,
		Z: z,
	}
}

func (v *Vector) Update(another Vector) {
	v.X += another.X
	v.Y += another.Y
	v.Z += another.Z
}
