package etcd

import (
	"context"
	"encoding/json"
	"errors"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

type Config struct {
	Addrs    []string `toml:"addrs" json:"addrs"`
	Timeout  int64    `toml:"timeout" json:"timeout"`
	UserName string   `toml:"user_name" json:"user_name"`
	Password string   `toml:"password" json:"password"`
}

func GetDefaultConfig() *Config {
	return &Config{
		Addrs:   []string{"127.0.0.1:2379"},
		Timeout: 10,
	}
}

type Client struct {
	cfg        *Config
	etcdClient *clientv3.Client
	ctx        context.Context
}

func NewClient(cfg *Config) (*Client, error) {
	client := &Client{
		cfg: cfg,
		ctx: context.TODO(),
	}

	err := client.Init()

	return client, err
}

func (c *Client) Init() error {
	if c == nil {
		return errors.New("client is nil")
	}

	if c.cfg == nil {
		c.cfg = GetDefaultConfig()
	}

	client, err := clientv3.New(clientv3.Config{
		Endpoints:   c.cfg.Addrs,
		DialTimeout: time.Duration(c.cfg.Timeout) * time.Second,
		Username:    c.cfg.UserName,
		Password:    c.cfg.Password,
	})

	if err != nil {
		return err
	}

	c.etcdClient = client
	return nil
}

func (c *Client) GetEtcdClient() *clientv3.Client {
	return c.etcdClient
}

func (c *Client) GetConfig() *Config {
	return c.cfg
}

func (c *Client) AddNode(key string, data any) error {
	dataByte, err := json.Marshal(data)

	if err != nil {
		return err
	}

	_, err = c.etcdClient.Put(context.Background(), key, string(dataByte))

	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetNode(key string, data any) error {
	rsp, err := c.etcdClient.Get(context.Background(), key)

	if err != nil {
		return err
	}

	if rsp == nil || len(rsp.Kvs) == 0 {
		return errors.New("key not exist")
	}

	return json.Unmarshal(rsp.Kvs[0].Value, data)
}

func (c *Client) GetRangeNode(prefix string) (result map[string][]byte, err error) {
	result = make(map[string][]byte)

	rsp, err := c.etcdClient.Get(context.Background(), prefix, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend))

	if err != nil {
		return nil, err
	}

	if rsp != nil && len(rsp.Kvs) > 0 {
		for _, item := range rsp.Kvs {
			result[string(item.Key)] = item.Value
		}
	}

	return
}

func (c *Client) DeleteNode(key string) error {
	_, err := c.etcdClient.Delete(context.Background(), key)

	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteNodeWithPrefix(prefix string) error {
	_, err := c.etcdClient.Delete(context.Background(), prefix, clientv3.WithPrefix())

	if err != nil {
		return err
	}

	return nil
}

func (c *Client) AddNodeWithLeaseKeepAlive(key string, data any, ttl int64) error {
	if ttl == 0 {
		ttl = 10
	}

	leaseRsp, err := c.etcdClient.Grant(context.Background(), ttl)

	if err != nil {
		return err
	}

	dataByte, err := json.Marshal(data)

	if err != nil {
		return err
	}

	_, err = c.etcdClient.Put(context.Background(), key, string(dataByte), clientv3.WithLease(leaseRsp.ID))

	if err != nil {
		return err
	}

	ch, err := c.etcdClient.KeepAlive(context.Background(), leaseRsp.ID)

	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ch:
			case <-c.ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (c *Client) WatchNode(key string, do func(clientv3.WatchResponse, context.CancelFunc)) error {
	ctx, cancel := context.WithCancel(context.TODO())
	ch := c.etcdClient.Watch(ctx, key)

	go func() {
		for {
			select {
			case wtRsp := <-ch:
				do(wtRsp, cancel)
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (c *Client) WatchNodesWithPrefix(key string, do func(clientv3.WatchResponse, context.CancelFunc)) error {
	ctx, cancel := context.WithCancel(context.TODO())
	ch := c.etcdClient.Watch(ctx, key, clientv3.WithPrefix(), clientv3.WithPrevKV())

	go func() {
		for {
			select {
			case wtRsp := <-ch:
				do(wtRsp, cancel)
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (c *Client) Close() {
	_ = c.etcdClient.Close()
	c.ctx.Done()
}
