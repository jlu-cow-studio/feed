package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/jlu-cow-studio/common/dal/mq"
	"github.com/jlu-cow-studio/common/dal/mysql"
	"github.com/jlu-cow-studio/common/dal/redis"
	"github.com/jlu-cow-studio/common/dal/rpc/feed_service"
	"github.com/jlu-cow-studio/common/discovery"
	"github.com/jlu-cow-studio/feed/handler"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 8080, "The server port")
)

func main() {
	discovery.Init()
	redis.Init()
	mysql.Init()
	mq.Init()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	feed_service.RegisterFeedServiceServer(s, &handler.Handler{})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
