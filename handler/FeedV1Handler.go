package handler

import (
	"context"
	"encoding/json"
	"log"
	"strconv"

	"github.com/jlu-cow-studio/common/dal/redis"
	"github.com/jlu-cow-studio/common/dal/rpc"
	"github.com/jlu-cow-studio/common/dal/rpc/base"
	"github.com/jlu-cow-studio/common/dal/rpc/feed_service"
	"github.com/jlu-cow-studio/common/dal/rpc/pack"
	redis_model "github.com/jlu-cow-studio/common/model/dao_struct/redis"
	"github.com/jlu-cow-studio/common/model/http_struct/feed"
	"github.com/jlu-cow-studio/feed/biz"
)

const PackServiceName = "cowstudio/pack"

func (h *Handler) FeedV1(ctx context.Context, req *feed_service.GetFeedRequest) (res *feed_service.GetFeedResponse, err error) {
	res = &feed_service.GetFeedResponse{
		Base: &base.BaseRes{
			Message: "",
			Code:    "498",
		},
		Pagination: &feed_service.PaginationInfo{
			CurrentPage:  req.Page,
			ItemsPerPage: req.PageSize,
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
		log.Printf("[AddItem] Unmarshal token error: %v", err)
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

	offset := req.Page*req.PageSize + 1
	cap := req.PageSize
	ids := []string{}
	switch req.Scene {
	case feed.FeedScene_Breeding,
		feed.FeedScene_Cattle_product,
		feed.FeedScene_WholeCattle,
		feed.FeedScene_Service,
		feed.FeedScene_ServiceProduct:
		ids, err = biz.GetToDisplayList(userInfo.Uid, req.Scene, int64(offset), int64(cap))
	case feed.FeedScene_HomePageMix:
		ids, err = biz.GetToDisplayList(userInfo.Uid, feed.FeedScene_WholeCattle, int64(offset), int64(cap))
	default:
		res.Base.Message = "unknowen scene"
		res.Base.Code = "403"
		log.Printf("[Feed] unknowen scene:")
		return
	}

	conn, err := rpc.GetConn(PackServiceName)
	if err != nil {
		res.Base.Message = "get rpc conn failed"
		res.Base.Code = "404"
		log.Printf("[Feed] get rpc conn failed: %v", err.Error())
		return
	}

	cli := pack.NewPackServiceClient(conn)

	itemIdList := []int32{}
	for _, id := range ids {
		if temp, e := strconv.Atoi(id); err != nil {
			res.Base.Message = e.Error()
			res.Base.Code = "405"
			log.Printf("[Feed] parse id failed: %v", e.Error())
			return
		} else {
			itemIdList = append(itemIdList, int32(temp))
		}
	}

	packReq := &pack.PackItemsReq{
		Base: &base.BaseReq{
			Token: req.Base.Token,
			Logid: req.Base.Logid,
		},
		ItemIdList: itemIdList,
	}

	packRes, _ := cli.PackItems(ctx, packReq)

	res.Base = packRes.Base
	if packRes.Base.Code != "200" {
		return
	}

	res.Items = packRes.ItemList
	res.Pagination.TotalItems = int32(biz.GetTotalItems(req.Scene))
	res.Base.Code = "200"
	res.Base.Message = ""
	return
}
