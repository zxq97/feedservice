package main

import (
	"feedservice/client/remind"
	"feedservice/conf"
	"feedservice/global"
	"feedservice/rpc/feed/pb"
	"feedservice/server"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/etcd"
)

var (
	feedConf   *conf.Conf
	remindConf *conf.Conf
	err        error
)

func main() {
	feedConf, err = conf.LoadYaml(conf.FeedConfPath)
	if err != nil {
		panic(err)
	}
	remindConf, err = conf.LoadYaml(conf.RemindConfPath)
	if err != nil {
		panic(err)
	}

	global.InfoLog, err = conf.InitLog(feedConf.LogPath.Info)
	if err != nil {
		panic(err)
	}
	global.ExcLog, err = conf.InitLog(feedConf.LogPath.Exc)
	if err != nil {
		panic(err)
	}
	global.DebugLog, err = conf.InitLog(feedConf.LogPath.Debug)
	if err != nil {
		panic(err)
	}

	remind.InitClient(remindConf)

	err = server.InitService(feedConf)
	if err != nil {
		panic(err)
	}

	etcdRegistry := etcd.NewRegistry(func(options *registry.Options) {
		options.Addrs = feedConf.Etcd.Addr
	})

	service := micro.NewService(
		micro.Name(feedConf.Grpc.Name),
		micro.Address(feedConf.Grpc.Addr),
		micro.Registry(etcdRegistry),
	)
	service.Init()
	err = feed_service.RegisterFeedServerHandler(
		service.Server(),
		new(server.FeedService),
	)
	if err != nil {
		panic(err)
	}
	err = service.Run()
	if err != nil {
		panic(err)
	}
}
