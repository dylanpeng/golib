package etcd

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"testing"
	"time"
)

var cfg *Config
var client *Client
var dir = "/test/node"
var user = &User{
	Name: "user_name",
	Age:  20,
}

type User struct {
	Name string
	Age  int
}

func init() {
	cfg = &Config{
		Addrs:   []string{"127.0.0.1:2379"},
		Timeout: 10,
	}

	var err error
	client, err = NewClient(cfg)

	if err != nil {
		fmt.Printf("new client fail. | err: %s\n", err)
	}

	return
}

func TestAddNode(t *testing.T) {
	err := client.AddNode(dir, user)

	if err != nil {
		t.Fatalf("TestAddNode AddNode fail. | err: %s", err)
	}
}

func TestGetNode(t *testing.T) {
	u := new(User)

	err := client.GetNode(dir, u)

	if err != nil {
		t.Fatalf("TestGetNode fail. | err: %s", err)
	}
}

func TestGetRangeNode(t *testing.T) {
	result, err := client.GetRangeNode("/test/")

	if err != nil {
		t.Fatalf("TestGetRangeNode GetRangeNode fail. | err: %s", err)
	}

	for k, v := range result {
		t.Logf("dir: %s | value: %s", k, string(v))
	}
}

func TestDeleteNode(t *testing.T) {
	t.Skip()
	err := client.DeleteNode("/delete/test/aaa")

	if err != nil {
		t.Fatalf("TestGetRangeNode GetRangeNode fail. | err: %s", err)
	}
}

func TestClient_DeleteNodeWithPrefix(t *testing.T) {
	t.Skip()
	err := client.DeleteNodeWithPrefix("/delete")

	if err != nil {
		t.Fatalf("TestGetRangeNode GetRangeNode fail. | err: %s", err)
	}
}

func TestAddNodeWithLeaseKeepAlive(t *testing.T) {
	t.Skip()
	err := client.AddNodeWithLeaseKeepAlive("/lease/node", "{}", 10)

	if err != nil {
		t.Fatalf("TestAddNodeWithLeaseKeepAlive AddNodeWithLeaseKeepAlive fail. | err: %s", err)
	}

	time.Sleep(1 * time.Minute)
}

func TestClient_WatchNode(t *testing.T) {
	t.Skip()
	err := client.WatchNode(dir, DealWithWatch)

	if err != nil {
		t.Fatalf("TestClient_WatchNode WatchNode fail. | err: %s", err)
	}

	time.Sleep(10 * time.Minute)
}

func TestClient_WatchNodesWithPrefix(t *testing.T) {
	t.Skip()
	err := client.WatchNodesWithPrefix("/test", DealWithWatch)

	if err != nil {
		t.Fatalf("TestClient_WatchNodesWithPrefix WatchNodesWithPrefix fail. | err: %s", err)
	}

	time.Sleep(10 * time.Minute)
}

func DealWithWatch(rsp clientv3.WatchResponse, cancel context.CancelFunc) {
	kvMap := make(map[string]string)
	if len(rsp.Events) > 0 {
		for _, item := range rsp.Events {
			kvMap[string(item.Kv.Key)] = string(item.Kv.Value)
			if item.PrevKv != nil {
				kvMap[string(item.PrevKv.Key)+"pre"] = string(item.PrevKv.Value)
			}
		}
	}

	cancel()
	return
}
