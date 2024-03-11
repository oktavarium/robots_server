package client

import (
	"fmt"

	"github.com/oktavarium/sgs/internal/internal/vector"
)

type Client struct {
	id         string
	pos        vector.Vector
	lastSeqNum int
	history    map[int]StateHistory
}

func NewClient(id string, pos vector.Vector) Client {
	return Client{
		id:      id,
		pos:     pos,
		history: make(map[int]StateHistory),
	}
}

func (c *Client) UpdateStateHistory(seqNum int) {
	c.history[seqNum] = NewStateHistory(c.pos)
	c.lastSeqNum = seqNum
}

func (c Client) SendString() string {
	return fmt.Sprintf("%d %s %f %f %f", c.lastSeqNum, c.id, c.pos.X, c.pos.Y, c.pos.Z)
}

func (c Client) LastSeqNumber() int {
	return c.lastSeqNum
}

func (c Client) CheckHistory(seqNumber int) bool {
	_, ok := c.history[seqNumber]
	return ok
}

func (c *Client) UpdatePosition(move vector.Vector) {
	c.pos.Update(move)
}
