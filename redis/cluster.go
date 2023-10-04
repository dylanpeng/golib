package redis

import (
	"errors"
	"github.com/redis/go-redis/v9"
	"sync"
)

type ClusterConfig struct {
	Addrs    []string `toml:"addrs" json:"addrs" yaml:"addrs"`
	Password string   `toml:"password" json:"password" yaml:"password"`
}

type ClusterPool struct {
	locker  sync.RWMutex
	clients map[string]*redis.ClusterClient
}

func (c *ClusterPool) Add(name string, conf *ClusterConfig) {
	c.locker.Lock()
	defer c.locker.Unlock()

	options := &redis.ClusterOptions{
		Addrs:          conf.Addrs,
		ReadOnly:       true,
		RouteByLatency: true,
	}

	c.clients[name] = redis.NewClusterClient(options)
}

func (c *ClusterPool) Get(name string) (client *redis.ClusterClient, err error) {
	c.locker.RLock()
	defer c.locker.RUnlock()

	client, ok := c.clients[name]

	if !ok {
		err = errors.New("no redis cluster client")
	}

	return
}

func NewClusterPool() *ClusterPool {
	return &ClusterPool{clients: make(map[string]*redis.ClusterClient)}
}
