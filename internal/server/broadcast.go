package server

import (
	"log/slog"
	"sync"
	"time"
)

func (s *server) broadcastSender(timeout time.Duration, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	ticker := time.NewTicker(timeout)
	defer ticker.Stop()
	for {
		select {
		case <-s.ctx.Done():
			return
		case data := <-s.broadcastCh:
			for k := range s.reg.GetClientsPos() {
				if _, err := s.conn.WriteToUDPAddrPort([]byte(data), k); err != nil {
					slog.Info("write to client", "error", err)
				}
			}
		case <-ticker.C:
			positions := s.reg.GetClientsPos()
			for k := range positions {
				for _, v := range positions {
					if _, err := s.conn.WriteToUDPAddrPort([]byte(v), k); err != nil {
						slog.Info("write to client", "error", err)
					}
				}
			}
		default:
		}
	}
}

func (s *server) sendBroadcast(data string) {
	select {
	case s.broadcastCh <- data:
	default:
	}
}
