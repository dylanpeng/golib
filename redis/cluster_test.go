package redis

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"testing"
)

var cacheClusterPool *ClusterPool
var clusterCache *redis.ClusterClient

// init redis pool
func TestInitClusterPool(t *testing.T) {
	cacheClusterPool = NewClusterPool()

	conf := &ClusterConfig{
		// 只需要设置主节点，不需要设置从节点
		Addrs: []string{"127.0.0.1:8001", "127.0.0.1:8002", "127.0.0.1:8003"},
	}

	cacheClusterPool.Add("test", conf)

	user = &User{
		Name:    "user",
		Gender:  1,
		Courses: make([]*Course, 0),
	}

	user.Courses = append(user.Courses,
		&Course{
			CourseName: "english",
			Score:      80,
		},
		&Course{
			CourseName: "math",
			Score:      95,
		})

	ctx = context.TODO()
	clusterCache, _ = cacheClusterPool.Get("test")
}

func TestClusterGetSet(t *testing.T) {
	userMarshal, err := json.Marshal(user)

	if err != nil {
		t.Fatalf("marshal user failed. err: %s", err)
	}

	key := GetKey("a")

	ok, err := clusterCache.Set(ctx, key, userMarshal, 0).Result()

	if err != nil {
		t.Fatalf("set cache key: %s faile. | result: %s | err: %s", key, ok, err)
	}

	ok, err = clusterCache.Set(ctx, "b", userMarshal, 0).Result()

	if err != nil {
		t.Fatalf("set cache key: %s faile. | result: %s | err: %s", key, ok, err)
	}

	ok, err = clusterCache.Set(ctx, "c", userMarshal, 0).Result()

	if err != nil {
		t.Fatalf("set cache key: %s faile. | result: %s | err: %s", key, ok, err)
	}

	cacheUserA := &User{}

	resultA, err := clusterCache.Get(ctx, "a").Result()

	if err != nil {
		t.Fatalf("get cache key: %s faile. | result: %s | err: %s", key, ok, err)
	}

	err = json.Unmarshal([]byte(resultA), cacheUserA)

	if err != nil {
		t.Fatalf("convert fail. err: %s", err)
	}

	t.Logf("cache get success. result: %s", cacheUserA)

	cacheUserB := &User{}

	resultB, err := clusterCache.Get(ctx, "b").Result()

	if err != nil {
		t.Fatalf("get cache key: %s faile. | result: %s | err: %s", key, ok, err)
	}

	err = json.Unmarshal([]byte(resultB), cacheUserB)

	if err != nil {
		t.Fatalf("convert fail. err: %s", err)
	}

	t.Logf("cache get success. result: %s", cacheUserB)

	cacheUserC := &User{}

	resultC, err := clusterCache.Get(ctx, "c").Result()

	if err != nil {
		t.Fatalf("get cache key: %s faile. | result: %s | err: %s", key, ok, err)
	}

	err = json.Unmarshal([]byte(resultC), cacheUserC)

	if err != nil {
		t.Fatalf("convert fail. err: %s", err)
	}

	t.Logf("cache get success. result: %s", cacheUserC)

}
