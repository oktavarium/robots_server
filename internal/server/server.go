package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/oktavarium/sgs/internal/internal/config"
	"github.com/oktavarium/sgs/internal/internal/registry"
)

const broadcastBufferSize = 1024
const udpBufferSize = 64
const maxConnectedClients = 2
const broadcastPositionTimeout = 100 * time.Millisecond

type server struct {
	ctx                      context.Context
	addr                     string
	conn                     *net.UDPConn
	reg                      *registry.Registry
	clientsIdIndex           int
	broadcastCh              chan string
	maxClients               int
	broadcastPositionTimeout time.Duration
}

func newServer(ctx context.Context, addr string, maxClients int, timeout time.Duration) server {
	if timeout == 0 {
		timeout = broadcastPositionTimeout
	}
	if maxClients <= 0 {
		maxClients = maxConnectedClients
	}

	return server{
		ctx:                      ctx,
		addr:                     addr,
		reg:                      registry.NewRegistry(),
		clientsIdIndex:           0,
		broadcastCh:              make(chan string, broadcastBufferSize),
		maxClients:               maxClients,
		broadcastPositionTimeout: timeout,
	}
}

func (s *server) listenAndServe() (err error) {
	udpAddr, err := net.ResolveUDPAddr("udp", s.addr)
	if err != nil {
		return fmt.Errorf("parse addr: %w", err)
	}
	s.conn, err = net.ListenUDP("udp", udpAddr)
	if err != nil {
		return fmt.Errorf("listen udp: %w", err)
	}
	defer func() {
		if err = s.conn.Close(); err != nil {
			err = fmt.Errorf("close udp socket: %w", err)
		}
	}()
	wg := new(sync.WaitGroup)
	go s.broadcastSender(s.broadcastPositionTimeout, wg)
Loop:
	for {
		select {
		case <-s.ctx.Done():
			break Loop
		default:
		}

		var buf [udpBufferSize]byte
		size, clientAddr, err := s.conn.ReadFromUDPAddrPort(buf[:])
		if err != nil {
			slog.Info("read from udp", "error", err)
			continue
		}

		go func() {
			wg.Add(1)
			defer wg.Done()
			s.handle(clientAddr, string(buf[:size]))
		}()
	}

	wg.Wait()
	close(s.broadcastCh)
	return nil
}

func Run() error {
	c := config.ParseConfig()
	s := newServer(context.Background(), c.GetAddress(), c.GetMaxClients(), c.GetTimeout())
	return s.listenAndServe()
}
