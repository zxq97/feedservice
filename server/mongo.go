package server

import (
	"context"
	"feedservice/global"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sort"
	"time"
)

type BoxItem struct {
	OwnerUID  int64 `bson:"owner_uid"`
	ArticleID int64 `bson:"article_id"`
	Ctime     int64 `bson:"ctime"`
}

type BoxItemSlice []BoxItem

func (b BoxItemSlice) Len() int {
	return len(b)
}

func (b BoxItemSlice) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b BoxItemSlice) Less(i, j int) bool {
	return b[i].Ctime < b[j].Ctime
}

func moPushSelfFeed(ctx context.Context, uid, articleID int64) error {
	inBox := BoxItem{
		OwnerUID:  uid,
		ArticleID: articleID,
		Ctime:     time.Now().Unix(),
	}
	var flag = true
	_, err := selfColl.UpdateOne(
		ctx,
		bson.M{"uid": uid},
		bson.M{"$push": bson.M{"in_box": inBox}, "$set": bson.M{"mtime": time.Now()}},
		&options.UpdateOptions{
			Upsert: &flag,
		},
	)
	if err != nil {
		global.ExcLog.Printf("ctx %v moPushSelfFeed uid %v articleid %v err %v", ctx, uid, articleID, err)
	}
	return err
}

func moPushFollowFeed(ctx context.Context, uids []int64, uid, articleID int64) error {
	outBox := BoxItem{
		OwnerUID:  uid,
		ArticleID: articleID,
		Ctime:     time.Now().Unix(),
	}
	var flag = true
	_, err := followColl.UpdateMany(
		ctx,
		bson.M{"uid": bson.M{"$in": uids}},
		bson.M{"$push": bson.M{"out_box": outBox}, "$set": bson.M{"ctime": time.Now()}},
		&options.UpdateOptions{
			Upsert: &flag,
		},
	)
	if err != nil {
		global.ExcLog.Printf("ctx %v moPushFollowFeed uids %v articleid %v err %v", ctx, uids, articleID, err)
	}
	return err
}

func moGetSelfFeed(ctx context.Context, uid int64) ([]int64, map[int64]int64, error) {
	var inBox struct {
		InBox []BoxItem `bson:"in_box"`
	}
	err := selfColl.FindOne(
		ctx,
		bson.M{"uid": uid},
	).Decode(&inBox)
	if err != nil {
		global.ExcLog.Printf("ctx %v moGetSelfFeed uid %v err %v", ctx, uid, err)
		return nil, nil, err
	}
	sort.Sort(BoxItemSlice(inBox.InBox))
	ids := make([]int64, 0, len(inBox.InBox))
	feedMap := make(map[int64]int64, len(inBox.InBox))
	for _, v := range inBox.InBox {
		feedMap[v.ArticleID] = v.Ctime
		ids = append(ids, v.ArticleID)
	}
	return ids, feedMap, nil
}

func moGetFollowFeed(ctx context.Context, uid int64) ([]int64, map[int64]int64, error) {
	var outBox struct {
		OutBox []BoxItem `bson:"out_box"`
	}
	err := followColl.FindOne(
		ctx,
		bson.M{"uid": uid},
	).Decode(&outBox)
	if err != nil {
		global.ExcLog.Printf("ctx %v moGetFollowFeed uid %v err %v", ctx, uid, err)
		return nil, nil, err
	}
	sort.Sort(BoxItemSlice(outBox.OutBox))
	ids := make([]int64, 0, len(outBox.OutBox))
	feedMap := make(map[int64]int64, len(outBox.OutBox))
	for _, v := range outBox.OutBox {
		feedMap[v.ArticleID] = v.Ctime
		ids = append(ids, v.ArticleID)
	}
	return ids, feedMap, nil
}

func moInsertFollowFeed(ctx context.Context, uid, toUID int64, feedMap map[int64]int64) error {
	outBox := make([]BoxItem, 0, len(feedMap))
	for k, v := range feedMap {
		outBox = append(outBox, BoxItem{
			OwnerUID:  toUID,
			ArticleID: k,
			Ctime:     v,
		})
	}
	var flag = true
	_, err := followColl.UpdateOne(
		ctx,
		bson.M{"uid": uid},
		bson.M{"$push": bson.M{"out_box": bson.M{"$each": outBox}}, "$set": bson.M{"mtime": time.Now()}},
		&options.UpdateOptions{
			Upsert: &flag,
		},
	)
	if err != nil {
		global.ExcLog.Printf("ctx %v moInsertFollowFeed uid %v touid %v feedmap %v err %v", ctx, uid, toUID, feedMap, err)
	}
	return err
}

func moPullFollowFeed(ctx context.Context, uid, toUID int64) error {
	_, err := followColl.UpdateMany(
		ctx,
		bson.M{"uid": uid},
		bson.M{"$pull": bson.M{"out_box": bson.M{"owner_uid": toUID}}},
	)
	if err != nil {
		global.ExcLog.Printf("ctx %v moPullFollowFeed uid %v touid %v err %v", ctx, uid, toUID, err)
	}
	return err
}
