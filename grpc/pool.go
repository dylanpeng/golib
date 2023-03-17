package grpc

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"sync"
	"time"
)

var (
	// ErrClosed is the error when the client pool is closed
	ErrClosed = errors.New("grpc pool: client pool is closed")
	// ErrTimeout is the error when the client pool timed out
	ErrTimeout = errors.New("grpc pool: client pool timed out")
	// ErrAlreadyClosed is the error when the client conn was already closed
	ErrAlreadyClosed = errors.New("grpc pool: the connection was already closed")
	// ErrFullPool is the error when the pool is already full
	ErrFullPool = errors.New("grpc pool: closing a ClientConn into a full pool")
)

type Factory func(addr string) (*grpc.ClientConn, error)

func DefaultFactory(addr string) (*grpc.ClientConn, error) {
	return grpc.DialContext(context.TODO(), addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
}

type Pool struct {
	addr     string
	mu       sync.RWMutex
	ch       chan *ClientConn
	factory  Factory
	capacity int
	idle     time.Duration
	ttl      time.Duration
	isClose  bool
	increase int64
}

func NewPool(factory Factory, addr string, capacity int, idle time.Duration, ttl ...time.Duration) *Pool {
	if capacity <= 0 {
		capacity = 1
	}

	p := &Pool{
		addr:     addr,
		mu:       sync.RWMutex{},
		ch:       make(chan *ClientConn, capacity),
		factory:  factory,
		capacity: capacity,
		idle:     idle,
	}

	if len(ttl) > 0 {
		p.ttl = ttl[0]
	}

	return p
}

func (p *Pool) GetConnChan() chan *ClientConn {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.ch
}

func (p *Pool) Get() (*ClientConn, error) {
	connCh := p.GetConnChan()
	if connCh == nil {
		return nil, ErrClosed
	}

	var clientConn *ClientConn

	for {
		select {
		case clientConn = <-connCh:
			if clientConn == nil || clientConn.ClientConn == nil {
				continue
			}

			if p.ttl > 0 && clientConn.createAt.Add(p.ttl).Before(time.Now()) {
				_ = clientConn.Close()
				clientConn.ClientConn = nil
				continue
			}

			if p.idle > 0 && clientConn.lastUsed.Add(p.idle).Before(time.Now()) {
				_ = clientConn.Close()
				clientConn.ClientConn = nil
				continue
			}

			if clientConn.GetState() == connectivity.Ready {
				clientConn.lastUsed = time.Now()
				return clientConn, nil
			}

		default:
			conn, err := p.factory(p.addr)

			if err != nil {
				return nil, err
			}

			clientConn = &ClientConn{
				ClientConn: conn,
				pool:       p,
				createAt:   time.Now(),
				lastUsed:   time.Now(),
			}

			p.mu.Lock()
			clientConn.Id = p.increase
			p.increase++
			p.mu.Unlock()

			return clientConn, nil
		}
	}
}

func (p *Pool) Close() {
	p.mu.Lock()
	ch := p.ch
	p.ch = nil
	p.mu.Unlock()

	p.isClose = true

	if ch == nil {
		return
	}

	close(ch)

	for item := range ch {
		if item != nil && item.ClientConn != nil {
			_ = item.Close()
		}
	}
}

type ClientConn struct {
	*grpc.ClientConn
	pool     *Pool
	createAt time.Time
	lastUsed time.Time
	Id       int64
}

func (c *ClientConn) Release() {
	if c == nil || c.ClientConn == nil {
		return
	}

	if c.pool.isClose {
		_ = c.Close()
		return
	}

	ch := c.pool.GetConnChan()

	if ch == nil {
		_ = c.Close()
		return
	}

	select {
	// return to channel
	case ch <- c:
	default:
		// channel is full
		_ = c.Close()
	}
}
