package main

import (
	"cft/api"
	"cft/config"
	"cft/log"
	pb "cft/proto"
	"cft/rpc"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"net"
)

func main() {
	config.ParseConfig("./config/application.json")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", config.GetConfig().RpcConfig.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterCftServerServer(s, &rpc.CftServer{})
	log.Infof("rpc server listening at %v", lis.Addr())
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	r := gin.Default()
	r.POST("/container", api.AddContainer)

	log.Infof("http server listening at %v", lis.Addr())
	r.Run(fmt.Sprintf(":%s", config.GetConfig().HttpConfig.Port))
}
