package server

import (
	"log/slog"
	"net/netip"
	"strconv"

	"github.com/oktavarium/sgs/internal/internal/client"
	"github.com/oktavarium/sgs/internal/internal/parser"
	"github.com/oktavarium/sgs/internal/internal/vector"
)

func (s *server) handle(clientAddr netip.AddrPort, data string) {
	switch data[0] {
	case 'n':
		if s.reg.Len() >= s.maxClients {
			return
		}

		s.clientsIdIndex++
		id := "c" + strconv.Itoa(s.clientsIdIndex) + "t"
		pos, err := parser.ParseInitialPosition(data)
		if err != nil {
			slog.Info("parsing initial position", "error", err)
			return
		}

		if _, err := s.conn.WriteToUDPAddrPort([]byte("a "+id), clientAddr); err != nil {
			slog.Info("write to client", "error", err)
		}

		newClient := client.NewClient(id, pos)
		if err := s.reg.NewClient(clientAddr, newClient); err != nil {
			slog.Info("creating client", "error", err)
			return
		}
	case 'e':
		s.reg.DeleteClient(clientAddr)
		// tell everyone that this player is deleted
		s.sendBroadcast(data)
	default:
		_, err := parser.ParseID(data)
		if err != nil {
			slog.Info("parse id", "error", err)
			return
		}
		seqNumber, err := parser.ParseSequenceNumber(data)
		if err != nil {
			slog.Info("parse sequence number", "error", err)
			return
		}

		userInput, err := parser.ParseInput(data)
		if err != nil {
			slog.Info("parse user input", "error", err)
			return
		}

		s.handleUserMoveInput(clientAddr, userInput, seqNumber)
	}
}

func (s *server) handleUserMoveInput(clientAddr netip.AddrPort, userInput string, seqNumber int) {
	cl, err := s.reg.GetClient(clientAddr)
	if err != nil {
		slog.Info("unknown client")
		return
	}

	if cl.LastSeqNumber() > seqNumber {
		return
	}

	if !cl.CheckHistory(seqNumber) {
		cl.UpdateStateHistory(seqNumber)
	}

	move := vector.Vector{}
	switch userInput {
	case "a":
		move.X = -0.1
	case "d":
		move.X = 0.1
	case "w":
		move.Y = 0.1
	case "s":
		move.Y = -0.1
	}

	cl.UpdatePosition(move)
	s.reg.UpdateClient(clientAddr, cl)
}
