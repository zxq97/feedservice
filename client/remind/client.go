package remind

import (
	"context"
	"feedservice/conf"
	"feedservice/global"
	"feedservice/rpc/remind/pb"
	"github.com/micro/go-micro"
)

var (
	client remind_service.RemindServerService
)

func InitClient(config *conf.Conf) {
	service := micro.NewService(micro.Name(
		config.Grpc.Name),
	)
	client = remind_service.NewRemindServerService(
		config.Grpc.Name,
		service.Client(),
	)
}

func AddUnread(ctx context.Context, uid int64, rType int32) error {
	_, err := client.AddUnread(ctx, toRemindInfo(uid, rType))
	if err != nil {
		global.ExcLog.Printf("ctx %v AddUnread uid %v rtype %v err %v", ctx, uid, rType, err)
	}
	return err
}

func AddBatchUnread(ctx context.Context, uids []int64, rType int32) error {
	_, err := client.AddBatchUnread(ctx, toBatchRemindRequest(uids, rType))
	if err != nil {
		global.ExcLog.Printf("ctx %v AddBatchUnread uids %v rtype %v err %v", ctx, uids, rType, err)
	}
	return err
}

func DeleteUnread(ctx context.Context, uid int64, rType int32) error {
	_, err := client.DeleteUnread(ctx, toRemindInfo(uid, rType))
	if err != nil {
		global.ExcLog.Printf("ctx %v DeleteUnread uid %v rtype %v err %v", ctx, uid, rType, err)
	}
	return err
}

func CheckUnread(ctx context.Context, uid int64, rType int32) (bool, error) {
	res, err := client.CheckUnread(ctx, toRemindInfo(uid, rType))
	if err != nil {
		global.ExcLog.Printf("ctx %v CheckUnread uid %v rtype %v err %v", ctx, uid, rType, err)
		return false, err
	}
	return res.Unread, nil
}
