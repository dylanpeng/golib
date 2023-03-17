package zookeeper

import (
	"fmt"
	"github.com/dylanpeng/golib/logger"
	"github.com/go-zookeeper/zk"
	"testing"
	"time"
)

var (
	hosts        = []string{"127.0.0.1:2181"}
	path         = "/test/dir/next/node"
	flags  int32 = zk.FlagEphemeral
	data         = []byte("zk data 001")
	acls         = zk.WorldACL(zk.PermAll)
	zkLog  *zkLogger
	Log    *logger.Logger
	cfg    *Config
	client *Client
)

func init() {
	var err error
	cfg = &Config{
		Addrs:   hosts,
		Timeout: 100,
	}

	Log, err = logger.NewLogger(&logger.Config{
		FilePath:   "./logs/confuse",
		Level:      "debug",
		TimeFormat: "2006-01-02 15:04:05.000",
		MaxAgeDay:  30,
	})

	if err != nil {
		fmt.Printf("init logger fail. | err: %s", err)
		return
	}

	client, err = NewClient(cfg, Log, nil)

	if err != nil {
		fmt.Printf("init NewClient fail. | err: %s", err)
		return
	}

	Log.Warningf("test")

	return
}

func TestClient_Create(t *testing.T) {
	t.Skip()
	err := client.Create(path, data, flags, acls)

	if err != nil {
		t.Fatalf("TestClient_Create Create fail. | err: %s", err)
		return
	}
}

func TestClient_Exist(t *testing.T) {
	t.Skip()
	exist, err := client.Exist(path)

	if err != nil {
		t.Fatalf("TestClient_Exist Exist fail. | err: %s", err)
		return
	}

	if exist {
		t.Logf("path: %s exist", path)
	} else {
		t.Logf("path: %s not exist", path)
	}

	return
}

func TestClient_Update(t *testing.T) {
	t.Skip()
	err := client.Update(path, []byte("update"))

	if err != nil {
		t.Fatalf("TestClient_Update Update fail. | err: %s", err)
		return
	}
}

func TestClient_Delete(t *testing.T) {
	t.Skip()
	err := client.Delete(path)

	if err != nil {
		t.Fatalf("TestClient_Delete Delete fail. | err: %s", err)
		return
	}
}

func TestClient_GetNode(t *testing.T) {
	t.Skip()
	result, err := client.GetNode(path)

	if err != nil {
		t.Fatalf("TestClient_GetNode GetNode fail. | err: %s", err)
		return
	}

	t.Logf("getnode: %s", string(result))

	return
}

func TestClient_GetNodes(t *testing.T) {
	t.Skip()
	mapNode, err := client.GetChildrenNodes("/root")

	if err != nil {
		t.Fatalf("TestClient_GetNodes GetNodes fail. | err: %s", err)
		return
	}

	for k, v := range mapNode {
		t.Logf("k: %s | v: %s", k, string(v))
	}
}

func TestClient_GetAllSubNodes(t *testing.T) {
	t.Skip()
	mapNode, err := client.GetAllSubNodes("/root")

	if err != nil {
		t.Fatalf("TestClient_GetNodes GetNodes fail. | err: %s", err)
		return
	}

	for k, v := range mapNode {
		t.Logf("k: %s | v: %s", k, string(v))
	}
}

func TestClient_Children(t *testing.T) {
	t.Skip()
	paths, err := client.Children("/root")

	if err != nil {
		t.Fatalf("TestClient_GetNodes GetNodes fail. | err: %s", err)
		return
	}

	for _, v := range paths {
		t.Logf("node: %s", v)
	}
}

func TestClient_WatchNode(t *testing.T) {
	client.WatchNode("/root/aaaa", WatchTypeExist, Do)

	time.Sleep(300 * time.Minute)
}

func Do(event zk.Event) {
	Log.Infof("event: %+v", event)
}
