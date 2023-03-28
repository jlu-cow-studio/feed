package handler

import (
	"github.com/jlu-cow-studio/common/dal/rpc/feed_service"
)

type Handler struct {
	feed_service.UnimplementedFeedServiceServer
}
