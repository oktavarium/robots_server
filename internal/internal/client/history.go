package client

import (
	"github.com/oktavarium/sgs/internal/internal/vector"
)

type StateHistory struct {
	pos vector.Vector
}

func NewStateHistory(pos vector.Vector) StateHistory {
	return StateHistory{
		pos: pos,
	}
}
