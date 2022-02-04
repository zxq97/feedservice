package server

import (
	"context"
	"feedservice/global"
	"feedservice/util/concurrent"
	"fmt"
)

func getSelfFeed(ctx context.Context, uid, cursor, offset int64) ([]int64, int64, error) {
	key := fmt.Sprintf(RedisKeyZInBox, uid)
	ids, hasMore, err := cacheGetFeed(ctx, key, cursor, offset)
	if err != nil {
		return nil, 0, err
	}
	if len(ids) == 0 {
		var feedMap map[int64]int64
		ids, feedMap, err = moGetSelfFeed(ctx, uid)
		if err != nil {
			return nil, 0, err
		}
		global.DebugLog.Printf("ctx %v uid %v cursor %v offset %v ids %v feedmap %v err %v", ctx, uid, cursor, offset, ids, feedMap, err)
		concurrent.Go(func() {
			_ = cacheSetFeed(ctx, key, feedMap)
		})
		if len(ids) > int(offset) {
			ids = ids[:offset]
			hasMore = true
		}
	}
	var nextCur int64
	if hasMore {
		nextCur = cursor + offset
	}
	return ids, nextCur, nil
}

func pushSelfFeed(ctx context.Context, uid, articleID int64) error {
	_ = refresh(ctx, uid)
	return moPushSelfFeed(ctx, uid, articleID)
}

func getFollowFeed(ctx context.Context, uid, cursor, offset int64) ([]int64, int64, error) {
	key := fmt.Sprintf(RedisKeyZOutBox, uid)
	ids, hasMore, err := cacheGetFeed(ctx, key, cursor, offset)
	if err != nil {
		return nil, 0, err
	}
	if len(ids) == 0 {
		var feedMap map[int64]int64
		ids, feedMap, err = moGetFollowFeed(ctx, uid)
		if err != nil {
			return nil, 0, err
		}
		concurrent.Go(func() {
			_ = cacheSetFeed(ctx, key, feedMap)
		})
		ids = ids[:offset]
		hasMore = true
	}
	var nextCur int64
	if hasMore {
		nextCur = cursor + offset
	}
	return ids, nextCur, nil
}

func pushFollowFeed(ctx context.Context, uids []int64, uid, articleID int64) error {
	return moPushFollowFeed(ctx, uids, uid, articleID)
}

func refresh(ctx context.Context, uid int64) error {
	key := fmt.Sprintf(RedisKeyZInBox, uid)
	return cacheClearFeed(ctx, key)
}

func followAfter(ctx context.Context, uid, toUID int64) error {
	_ = refresh(ctx, uid)
	_, feedMap, err := moGetSelfFeed(ctx, toUID)
	if err != nil {
		return err
	}
	return moInsertFollowFeed(ctx, uid, toUID, feedMap)
}

func unfollowAfter(ctx context.Context, uid, toUID int64) error {
	_ = refresh(ctx, uid)
	return moPullFollowFeed(ctx, uid, toUID)
}
