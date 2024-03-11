package registry

import (
	"net/netip"
	"sync"

	"github.com/oktavarium/sgs/internal/internal/client"
	"github.com/oktavarium/sgs/internal/internal/sgserr"
)

type Registry struct {
	clients map[netip.AddrPort]client.Client
	mutex   sync.RWMutex
}

func NewRegistry() *Registry {
	return &Registry{
		clients: make(map[netip.AddrPort]client.Client),
	}
}

func (r *Registry) Len() int {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return len(r.clients)
}

func (r *Registry) NewClient(addr netip.AddrPort, client client.Client) error {
	if r.ClientExists(addr) {
		return sgserr.ErrConflict
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.clients[addr] = client

	return nil
}

func (r *Registry) ClientExists(addr netip.AddrPort) bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	_, ok := r.clients[addr]
	return ok
}

func (r *Registry) DeleteClient(addr netip.AddrPort) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	delete(r.clients, addr)
}

func (r *Registry) GetClientsPos() map[netip.AddrPort]string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	result := make(map[netip.AddrPort]string, len(r.clients))
	for k, v := range r.clients {
		result[k] = v.SendString()
	}
	return result
}

func (r *Registry) GetClient(addr netip.AddrPort) (client.Client, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	c, ok := r.clients[addr]
	if !ok {
		return c, sgserr.ErrNotFound
	}
	return c, nil
}

func (r *Registry) UpdateClient(addr netip.AddrPort, cl client.Client) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.clients[addr] = cl
}
