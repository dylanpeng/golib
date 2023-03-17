package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
	"strings"
	"testing"
)

var cachePool *Pool
var user *User
var ctx context.Context
var cache *redis.Client

type User struct {
	Name    string    `json:"name"`
	Gender  int       `json:"gender"`
	Courses []*Course `json:"courses"`
}

func (u *User) String() string {
	return fmt.Sprintf("%+v", *u)
}

type Course struct {
	CourseName string `json:"course_name"`
	Score      int    `json:"score"`
}

func (c *Course) String() string {
	return fmt.Sprintf("%+v", *c)
}

func GetKey(prefix string, items ...interface{}) string {
	format := prefix + strings.Repeat(":%v", len(items))
	return fmt.Sprintf(format, items...)
}

// init redis pool
func TestInitPool(t *testing.T) {
	t.Skip()

	cachePool = NewPool()

	conf := &Config{
		Host:     "127.0.0.1",
		Port:     6379,
		Password: "",
	}

	cachePool.Add("test", conf)

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
	cache, _ = cachePool.Get("test")
}

// redis cache get set
func TestCacheGetSet(t *testing.T) {
	t.Skip()

	userMarshal, err := json.Marshal(user)

	if err != nil {
		t.Fatalf("marshal user failed. err: %s", err)
	}

	key := GetKey("test", "user", "course")

	result, err := cache.Set(ctx, key, userMarshal, 0).Result()

	if err != nil {
		t.Fatalf("set cache key: %s faile. | result: %s | err: %s", key, result, err)
	}

	t.Logf("set user success. result: %s", result)

	cacheUser := &User{}

	result, err = cache.Get(ctx, key).Result()

	if err != nil {
		t.Fatalf("set cache key: %s faile. | result: %s | err: %s", key, result, err)
	}

	err = json.Unmarshal([]byte(result), cacheUser)

	if err != nil {
		t.Fatalf("convert fail. err: %s", err)
	}

	t.Logf("cache get success. result: %s", cacheUser)
}

// redis incr
func TestIncr(t *testing.T) {
	t.Skip()

	key := GetKey("test", "cache", "incr")
	result, err := cache.Incr(ctx, key).Result()

	if err != nil {
		t.Fatalf("incr fail. | err: %s | result: %d", err, result)
	}

	strResult, err := cache.Get(ctx, key).Result()

	if err != nil {
		t.Fatalf("get incr fail. | err: %s | result: %d", err, result)
	}

	t.Logf("get incr result: %s", strResult)

	result, err = cache.IncrBy(ctx, key, 5).Result()

	if err != nil {
		t.Fatalf("incr fail. | err: %s | result: %d", err, result)
	}

	strResult, err = cache.Get(ctx, key).Result()

	if err != nil {
		t.Fatalf("get incr fail. | err: %s | result: %d", err, result)
	}

	intResult, _ := strconv.Atoi(strResult)

	t.Logf("get incr result: %d", intResult)
}

// redis hget hset
func TestHGetSet(t *testing.T) {
	t.Skip()

	key := GetKey("test", "cache", "hget")

	resultInt, err := cache.HSet(ctx, key, "key1", "val1", "key2", "val2").Result()

	if err != nil {
		t.Fatalf("hset fail. | result: %d | err: %s", resultInt, err)
	}

	result, err := cache.HGet(ctx, key, "key1").Result()

	if err != nil {
		t.Fatalf("hget fail. | result: %s | err: %s", result, err)
	}

	resultMap, err := cache.HGetAll(ctx, key).Result()

	if err != nil {
		t.Fatalf("hget fail. | result: %s | err: %s", result, err)
	}

	t.Logf("success. resultInt: %d | result: %s | resultMap: %s", resultInt, result, resultMap)
}

// redis lpush lpop
func TestLPushPop(t *testing.T) {
	t.Skip()

	key := GetKey("test", "cache", "LPush")

	_, err := cache.LPush(ctx, key, "aaa").Result()

	if err != nil {
		t.Fatalf("lpush fail. | err: %s", err)
	}

	_, err = cache.LPush(ctx, key, "bbb").Result()

	if err != nil {
		t.Fatalf("lpush fail. | err: %s", err)
	}

	result, err := cache.LPop(ctx, key).Result()

	if err != nil {
		t.Fatalf("lpop fail. | err: %s", err)
	}

	t.Logf("lpop success. | result: %s", result)
}

// redis sadd spop
func TestSAddPop(t *testing.T) {
	t.Skip()

	key := GetKey("test", "cache", "sAdd")

	_, err := cache.SAdd(ctx, key, "aaa", "bbb", "ccc").Result()

	if err != nil {
		t.Fatalf("sadd fail. | err: %s", err)
	}

	members, err := cache.SMembers(ctx, key).Result()

	if err != nil {
		t.Fatalf("smembers fail. | err: %s", err)
	}

	t.Logf("smembers: %s", members)

	membersMap, err := cache.SMembersMap(ctx, key).Result()

	if err != nil {
		t.Fatalf("smembersMap fail. | err: %s", err)
	}

	t.Logf("smembersMap: %s", membersMap)

	result, err := cache.SPop(ctx, key).Result()

	if err != nil {
		t.Fatalf("spop fail. | err: %s", err)
	}

	t.Logf("spop. result: %s", result)
}

// redis zadd
func TestZAdd(t *testing.T) {
	t.Skip()

	key := GetKey("test", "cache", "zadd")

	list := make([]redis.Z, 0)
	list = append(list, redis.Z{
		Score:  1,
		Member: "a",
	})

	list = append(list, redis.Z{
		Score:  10,
		Member: "b",
	})

	list = append(list, redis.Z{
		Score:  15,
		Member: "c",
	})

	list = append(list, redis.Z{
		Score:  15,
		Member: "d",
	})

	_, err := cache.ZAdd(ctx, key, list...).Result()

	if err != nil {
		t.Fatalf("zadd fail. | err: %s", err)
	}

	results, err := cache.ZRange(ctx, key, 0, -1).Result()

	if err != nil {
		t.Fatalf("zrange fail. | err: %s", err)
	}

	t.Logf("zrange list. %s", results)

	resultList, err := cache.ZRangeWithScores(ctx, key, 0, -1).Result()

	if err != nil {
		t.Fatalf("zrange fail. | err: %s", err)
	}

	t.Logf("zrange with scores list. %+v", resultList)

	cache.ZRem(ctx, key, "d")
}

// redis geo
func TestGeo(t *testing.T) {
	t.Skip()

	key := GetKey("test", "cache", "geo")

	locations := make([]*redis.GeoLocation, 0)

	locations = append(locations, &redis.GeoLocation{
		Name:      "beijing",
		Longitude: 116.38,
		Latitude:  39.92,
	})

	locations = append(locations, &redis.GeoLocation{
		Name:      "tiantan",
		Longitude: 116.48,
		Latitude:  39.94,
	})

	locations = append(locations, &redis.GeoLocation{
		Name:      "tianjin",
		Longitude: 117.15,
		Latitude:  39.12,
	})

	_, err := cache.GeoAdd(ctx, key, locations...).Result()

	if err != nil {
		t.Fatalf("geo add fail. | err: %s", err)
	}

	geoPos, err := cache.GeoPos(ctx, key, "beijing", "tianjin").Result()

	if err != nil {
		t.Fatalf("geo pos add fail. | err: %s", err)
	}

	t.Logf("geo pos. result: %+v", geoPos)

	bDist, err := cache.GeoDist(ctx, key, "beijing", "tianjin", "km").Result()

	if err != nil {
		t.Fatalf("geo dist add fail. | err: %s", err)
	}

	t.Logf("geo pos. result: %fkm", bDist)

	t.Logf("f:%10.2f", 11.22)

	prop := &redis.GeoRadiusQuery{
		Radius:      10000,
		Unit:        "km",
		WithCoord:   true,
		WithDist:    true,
		WithGeoHash: true,
		Sort:        "ASC",
		Store:       "",
		StoreDist:   "",
	}

	geoList, err := cache.GeoRadius(ctx, key, 116.1, 39.1, prop).Result()

	if err != nil {
		t.Fatalf("geo dist add fail. | err: %s", err)
	}

	t.Logf("geo radius. result: %+v", geoList)

	prop1 := &redis.GeoRadiusQuery{
		Radius:      10000,
		Unit:        "km",
		WithCoord:   true,
		WithDist:    true,
		WithGeoHash: true,
		Sort:        "ASC",
		Store:       "",
		StoreDist:   "",
	}

	geoListMem, err := cache.GeoRadiusByMember(ctx, key, "beijing", prop1).Result()

	if err != nil {
		t.Fatalf("geo radius member fail. | err: %s", err)
	}

	t.Logf("geo radius member. result: %+v", geoListMem)

	geoHash, err := cache.GeoHash(ctx, key, "beijing", "tianjin").Result()

	if err != nil {
		t.Fatalf("geo hash fail. | err: %s", err)
	}

	t.Logf("geo hash result: %+v", geoHash)
}
