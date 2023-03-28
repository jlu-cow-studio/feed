package handler

import (
	"context"
	"testing"

	"github.com/jlu-cow-studio/common/dal/mq"
	"github.com/jlu-cow-studio/common/dal/mysql"
	"github.com/jlu-cow-studio/common/dal/redis"
	"github.com/jlu-cow-studio/common/dal/rpc/base"
	"github.com/jlu-cow-studio/common/dal/rpc/feed_service"
	"github.com/jlu-cow-studio/common/discovery"
	"github.com/sanity-io/litter"
)

func TestGetFeed(t *testing.T) {

	discovery.Init()
	redis.Init()
	mysql.Init()
	mq.Init()
	req := &feed_service.GetFeedRequest{
		Base: &base.BaseReq{
			Token: "7f9a2748-9cf3-47d4-b69d-cca7d09141c4",
			Logid: "",
		},
		Scene:    "cattle_product",
		Page:     0,
		PageSize: 20,
	}

	litter.Dump(req)
	res, err := new(Handler).GetFeed(context.Background(), req)
	litter.Dump(res)
	litter.Dump(err)
}
