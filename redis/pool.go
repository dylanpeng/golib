package redis

import (
	"errors"
	"github.com/redis/go-redis/v9"
	"net"
	"strconv"
	"sync"
)

type Config struct {
	Host     string `toml:"host" json:"host"`
	Port     int    `toml:"port" json:"port"`
	Password string `toml:"password" json:"password"`
}

func (c *Config) GetAddr() string {
	return net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
}

type Pool struct {
	locker  sync.RWMutex
	clients map[string]*redis.Client
}

func (p *Pool) Add(name string, conf *Config) {
	p.locker.Lock()
	defer p.locker.Unlock()

	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.GetAddr(),
		Password: conf.Password,
	})

	p.clients[name] = rdb
	return
}

func (p *Pool) Get(name string) (client *redis.Client, err error) {
	p.locker.RLock()
	defer p.locker.RUnlock()

	client, ok := p.clients[name]

	if !ok {
		err = errors.New("no redis client")
	}

	return
}

func NewPool() *Pool {
	return &Pool{clients: make(map[string]*redis.Client)}
}
