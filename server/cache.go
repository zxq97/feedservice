package server

import (
	"context"
	"feedservice/global"
	"feedservice/util/cast"
	"github.com/go-redis/redis/v8"
	"time"
)

const (
	FeedTTL = 5 * time.Second

	RedisKeyZInBox  = "feed_service_in_box_%v"
	RedisKeyZOutBox = "feed_service_out_box_%v"
)

func cacheGetFeed(ctx context.Context, key string, cursor, offset int64) ([]int64, bool, error) {
	val, err := redisCli.ZRevRange(ctx, key, cursor, cursor+offset).Result()
	if err != nil && err != redis.Nil {
		global.ExcLog.Printf("ctx %v cacheGetFeed key %v cursor %v err %v", ctx, key, cursor, err)
		return nil, false, err
	}
	ids := make([]int64, 0, offset)
	var hasMore bool
	if len(val) > int(offset) {
		hasMore = true
	}
	for _, v := range val {
		if len(ids) == int(offset) {
			break
		}
		ids = append(ids, cast.ParseInt(v, 0))
	}
	return ids, hasMore, nil
}

func cacheSetFeed(ctx context.Context, key string, feedMap map[int64]int64) error {
	zs := make([]*redis.Z, 0, len(feedMap))
	for k, v := range feedMap {
		zs = append(zs, &redis.Z{
			Member: k,
			Score:  float64(v),
		})
	}
	global.DebugLog.Printf("ctx %v key %v feedmap %v", ctx, key, feedMap)
	err := redisCli.ZAdd(ctx, key, zs...).Err()
	global.DebugLog.Printf("ctx %v err %v", ctx, err)
	if err != nil {
		global.ExcLog.Printf("ctx %v cacheSetFeed key %v feedmap %v err %v", ctx, key, feedMap, err)
		return err
	}
	redisCli.Expire(ctx, key, FeedTTL)
	return nil
}

func cacheClearFeed(ctx context.Context, key string) error {
	err := redisCli.Del(ctx, key).Err()
	if err != nil {
		global.ExcLog.Printf("ctx %v cacheClearFeed key %v err %v", ctx, key, err)
	}
	return err
}
