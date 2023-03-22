package eth

import (
	"context"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type Config struct {
	Addr string `toml:"addr" json:"addr"`
}

type Client struct {
	c      *Config
	client *ethclient.Client
}

func (c *Client) GetClient() *ethclient.Client {
	return c.client
}

func (c *Client) GetConfig() *Config {
	return c.c
}

func NewClient(config *Config) (*Client, error) {
	return NewClientWithOptions(context.Background(), config)
}

func NewClientWithOptions(ctx context.Context, config *Config, options ...rpc.ClientOption) (*Client, error) {
	client := &Client{
		c: config,
	}

	rpcClient, err := rpc.DialOptions(ctx, config.Addr, options...)

	if err != nil {
		return nil, err
	}

	client.client = ethclient.NewClient(rpcClient)

	return client, nil
}
