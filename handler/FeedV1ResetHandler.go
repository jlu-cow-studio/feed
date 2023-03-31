package handler

import (
	"context"
	"encoding/json"
	"log"

	"github.com/jlu-cow-studio/common/dal/redis"
	"github.com/jlu-cow-studio/common/dal/rpc/base"
	"github.com/jlu-cow-studio/common/dal/rpc/feed_service"
	redis_model "github.com/jlu-cow-studio/common/model/dao_struct/redis"
	"github.com/jlu-cow-studio/common/model/http_struct/feed"
	"github.com/jlu-cow-studio/feed/biz"
)

func (h *Handler) FeedV1Reset(ctx context.Context, req *feed_service.ResetFeedRequest) (res *feed_service.ResetFeedResponse, err error) {
	res = &feed_service.ResetFeedResponse{
		Base: &base.BaseRes{
			Message: "",
			Code:    "",
		},
	}

	// 获取 token
	token := req.Base.Token
	cmd := redis.DB.Get(redis.GetUserTokenKey(token))
	if cmd.Err() != nil {
		res.Base.Message = cmd.Err().Error()
		res.Base.Code = "400"
		log.Printf("[Feed] Redis get token error: %v", cmd.Err())
		return
	}
	log.Printf("[Feed] Redis get token success, token: %s", token)

	// 解析 token 中的用户信息
	userInfo := new(redis_model.UserInfo)
	if err = json.Unmarshal([]byte(cmd.Val()), userInfo); err != nil {
		res.Base.Message = err.Error()
		res.Base.Code = "401"
		log.Printf("[Feed] Unmarshal token error: %v", err)
		return
	}
	log.Printf("[Feed] Unmarshal token success, userinfo: %v", userInfo)

	//校验用户场景匹配
	if _, ok := feed.Role_FeedScene[userInfo.Role][req.Scene]; !ok {
		res.Base.Message = "role scenen not match"
		res.Base.Code = "402"
		log.Printf("[Feed] role scenen not match: %v:%v", userInfo.Role, req.Scene)
		return
	}

	log.Printf("[Feed] get feed, user role: %v, scene %v", userInfo.Role, req.Scene)

	if err = biz.ResetCache(userInfo.Uid, req.Scene); err != nil {
		res.Base.Message = err.Error()
		res.Base.Code = "403"
		log.Printf("[Feed] reset cache failed: %v", err)
		return
	}

	res.Base.Message = ""
	res.Base.Code = "200"

	return
}
