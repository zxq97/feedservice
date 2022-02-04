package server

import (
	"context"
	"feedservice/client/remind"
	"feedservice/conf"
	"feedservice/rpc/feed/pb"
	"feedservice/util/constant"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
)

type FeedService struct {
}

var (
	redisCli   redis.Cmdable
	selfColl   *mongo.Collection
	followColl *mongo.Collection
)

const (
	colFollowFeed = "follow_feed"
	colSelfFeed   = "self_feed"
)

func InitService(config *conf.Conf) error {
	redisCli = conf.GetRedisCluster(config.RedisCluster.Addr)
	mongoCli, err := conf.GetMongo(config.Mongo.Addr)
	if err != nil {
		return err
	}
	mongoDB := mongoCli.Database(config.Mongo.DBName)
	followColl = mongoDB.Collection(colFollowFeed)
	selfColl = mongoDB.Collection(colSelfFeed)
	return nil
}

func (fs *FeedService) GetSelfFeed(ctx context.Context, req *feed_service.GetFeedRequest, res *feed_service.GetFeedResponse) error {
	ids, nextCur, err := getSelfFeed(ctx, req.Uid, req.Cursor, req.Offset)
	if err != nil {
		return err
	}
	res.ArticleIds = ids
	res.NextCursor = nextCur
	return nil
}

func (fs *FeedService) PushSelfFeed(ctx context.Context, req *feed_service.PushSelfFeedRequest, res *feed_service.EmptyResponse) error {
	return pushSelfFeed(ctx, req.Uid, req.ArticleId)
}

func (fs *FeedService) GetFollowFeed(ctx context.Context, req *feed_service.GetFeedRequest, res *feed_service.GetFeedResponse) error {
	ids, nextCur, err := getFollowFeed(ctx, req.Uid, req.Cursor, req.Offset)
	if err != nil {
		return err
	}
	res.ArticleIds = ids
	res.NextCursor = nextCur
	return nil
}

func (fs *FeedService) PushFollowFeed(ctx context.Context, stream feed_service.FeedServer_PushFollowFeedStream) error {
	defer stream.Close()
	for {
		req, err := stream.Recv()
		if err == nil {
			uids := req.Uids
			articleID := req.ArticleId
			err = pushFollowFeed(ctx, uids, req.Uid, articleID)
			if err != nil {
				return err
			}
			_ = remind.AddBatchUnread(ctx, uids, constant.RemindTypeFollowFeed)
		} else if err == io.EOF {
			break
		} else {
			return err
		}
	}
	return nil
}

func (fs *FeedService) FollowAfterFeed(ctx context.Context, req *feed_service.ActionFeedRequest, res *feed_service.EmptyResponse) error {
	return followAfter(ctx, req.Uid, req.ToUid)
}

func (fs *FeedService) UnfollowAfterFeed(ctx context.Context, req *feed_service.ActionFeedRequest, res *feed_service.EmptyResponse) error {
	return unfollowAfter(ctx, req.Uid, req.ToUid)
}

func (fs *FeedService) Refresh(ctx context.Context, req *feed_service.ReFreshRequest, res *feed_service.EmptyResponse) error {
	return refresh(ctx, req.Uid)
}
