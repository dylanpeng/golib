package zookeeper

import (
	"context"
	"errors"
	"fmt"
	"github.com/dylanpeng/golib/logger"
	"github.com/go-zookeeper/zk"
	"strings"
	"sync"
	"time"
)

const (
	EventTypeAll = 0
)

const (
	WatchTypeData = iota
	WatchTypeExist
	WatchTypeChild
)

type Config struct {
	Addrs   []string `toml:"addrs" json:"addrs" yaml:"addrs"`
	Timeout int64    `toml:"timeout" json:"timeout" yaml:"timeout"`
}

type zkLogger struct {
	*logger.Logger
}

func (l *zkLogger) Printf(format string, args ...interface{}) {
	l.Errorf(format, args...)
}

type Client struct {
	conf     *Config
	conn     *zk.Conn
	logger   *logger.Logger
	callback func(zk.Event)
	ctx      context.Context
	wg       *sync.WaitGroup
}

func NewClient(conf *Config, logger *logger.Logger, callback func(event zk.Event)) (*Client, error) {
	if conf.Timeout == 0 {
		conf.Timeout = 10
	}

	client := &Client{
		conf:     conf,
		logger:   logger,
		callback: callback,
		wg:       &sync.WaitGroup{},
		ctx:      context.TODO(),
	}

	conn, _, err := zk.Connect(conf.Addrs, time.Duration(conf.Timeout)*time.Second,
		zk.WithLogger(&zkLogger{logger}),
		zk.WithLogInfo(false),
		zk.WithEventCallback(callback),
	)

	if err != nil {
		return nil, err
	}

	client.conn = conn

	return client, nil
}

func (c *Client) Create(path string, data []byte, flags int32, acl []zk.ACL) error {
	if acl == nil {
		acl = zk.WorldACL(zk.PermAll)
	}

	pathSpilt := strings.Split(path, "/")

	for i := 2; i < len(pathSpilt); i++ {
		p := strings.Join(pathSpilt[:i], "/")

		exits, err := c.Exist(p)

		if err != nil {
			c.logger.Errorf("zookeeper client create exists fail. | err: %s", err)
			return err
		}

		if !exits {
			_, err = c.conn.CreateContainer(p, nil, zk.FlagTTL, acl)

			if err != nil {
				c.logger.Errorf("zookeeper client create container fail. | err: %s", err)
				return err
			}
		}
	}

	_, err := c.conn.Create(path, data, flags, acl)

	if err != nil {
		c.logger.Errorf("zookeeper client create fail. | err: %s", err)
		return err
	}

	return nil
}

func (c *Client) Exist(path string) (exists bool, err error) {
	exists, _, err = c.conn.Exists(path)
	return
}

func (c *Client) Update(path string, data []byte) error {
	exists, stat, err := c.conn.Exists(path)

	if err != nil {
		c.logger.Infof("zookeeper client update exists fail. | err: %s", err)
		return err
	}

	if !exists {
		c.logger.Infof("zookeeper client update path not exists. | err: %s", err)
		return errors.New(fmt.Sprintf("path: %s not exists", path))
	}

	_, err = c.conn.Set(path, data, stat.Version)

	if err != nil {
		c.logger.Infof("zookeeper client update path fail. | err: %s", err)
		return err
	}

	return nil
}

func (c *Client) Delete(path string) error {
	exists, stat, err := c.conn.Exists(path)

	if err != nil {
		c.logger.Infof("zookeeper client delete exists fail. | err: %s", err)
		return err
	}

	if !exists {
		c.logger.Infof("zookeeper client delete path not exists. | err: %s", err)
		return errors.New(fmt.Sprintf("path: %s not exists", path))
	}

	err = c.conn.Delete(path, stat.Version)

	if err != nil {
		c.logger.Infof("zookeeper client delete path fail. | err: %s", err)
		return err
	}

	return nil
}

func (c *Client) GetNode(path string) (data []byte, err error) {
	data, _, err = c.conn.Get(path)

	if err != nil {
		c.logger.Infof("zookeeper client get node path fail. | err: %s", err)
		return nil, err
	}

	return
}

func (c *Client) GetChildrenNodes(path string) (dataMap map[string][]byte, err error) {
	nodes, _, err := c.conn.Children(path)

	if err != nil {
		return
	}

	if len(nodes) == 0 {
		return nil, nil
	}

	dataMap = make(map[string][]byte, len(nodes))

	for _, node := range nodes {
		var data []byte
		nodePath := fmt.Sprintf("%s/%s", path, node)
		data, _, err = c.conn.Get(nodePath)

		if err != nil {
			c.logger.Infof("zookeeper client get nodes get fail. | err: %s", err)
			return nil, err
		}

		dataMap[node] = data
	}

	return
}

func (c *Client) GetAllSubNodes(path string) (dataMap map[string][]byte, err error) {
	nodes, _, err := c.conn.Children(path)

	if err != nil {
		return
	}

	if len(nodes) == 0 {
		return nil, nil
	}

	dataMap = make(map[string][]byte)

	for _, node := range nodes {
		var data []byte
		nodePath := fmt.Sprintf("%s/%s", path, node)
		data, _, err = c.conn.Get(nodePath)

		if err != nil {
			c.logger.Infof("zookeeper client get nodes get fail. | err: %s", err)
			return nil, err
		}

		dataMap[nodePath] = data
		var childMap map[string][]byte

		childMap, err = c.GetChildrenNodes(nodePath)

		if err != nil {
			c.logger.Infof("zookeeper client get child nodes fail. | err: %s", err)
			return nil, err
		}

		for k, v := range childMap {
			dataMap[nodePath+"/"+k] = v
		}

	}

	return
}

func (c *Client) Children(path string) (children []string, err error) {
	children, _, err = c.conn.Children(path)
	return
}

func (c *Client) Close() {
	c.conn.Close()
	c.wg.Wait()
}

func (c *Client) WatchNode(path string, watchType int, do func(event zk.Event)) {
	c.wg.Add(1)

	go func() {
		defer c.wg.Done()

		for {
			select {
			case <-c.ctx.Done():
				return
			default:
				var eventCh <-chan zk.Event
				var err error

				switch watchType {
				case WatchTypeData:
					_, _, eventCh, err = c.conn.GetW(path)
				case WatchTypeExist:
					_, _, eventCh, err = c.conn.ExistsW(path)
				case WatchTypeChild:
					_, _, eventCh, err = c.conn.ChildrenW(path)
				}

				if err != nil {
					c.logger.Warningf("zookeeper watch failed | path: %s | error: %s", path, err)
					continue
				}

				event := <-eventCh

				if event.Err != nil {
					if event.Err != zk.ErrClosing {
						c.logger.Warningf("zookeeper watch failed | path: %s | error: %s", path, event.Err)
					}

					continue
				}

				do(event)
			}
		}
	}()
}
