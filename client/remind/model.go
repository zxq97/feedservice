package remind

import (
	"feedservice/rpc/remind/pb"
)

func toRemindInfo(uid int64, rType int32) *remind_service.RemindInfo {
	return &remind_service.RemindInfo{
		Uid:        uid,
		RemindType: rType,
	}
}

func toBatchRemindRequest(uids []int64, rType int32) *remind_service.RemindBatchRequest {
	return &remind_service.RemindBatchRequest{
		Uids:       uids,
		RemindType: rType,
	}
}
