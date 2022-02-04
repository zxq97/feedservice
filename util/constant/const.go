package constant

import "time"

const (
	BatchSize = 1000

	FollowTypePerson = 1
	FollowTypeTopic  = 2

	GenderUndefined = 0
	GenderBody      = 1
	GenderGirl      = 2

	RemindTypeFollowFeed = 1
	RemindTypeComment    = 2
	RemindTypeChat       = 3

	RPCTimeOut = 3 * time.Second
)
