package rpc

import (
	"cft/log"
	"cft/model"
	"cft/monitor"
	pb "cft/proto"
	"cft/utils"
	"context"
	"os"
	"sync"
)

type StateType int

const (
	StatusOk int32 = iota
	StatusNotFound
)

type CftServer struct {
	pb.UnimplementedCftServerServer
	Id   int
	dead int32
	lock sync.RWMutex
}

func (s *CftServer) HandleHeartbeat(ctx context.Context, req *pb.HeartbeatRequest) (*pb.HeartBeatReply, error) {
	container := monitor.GetContainer(req.ContainerId)
	reply := &pb.HeartBeatReply{}
	if container == nil {
		monitor.AddContainer(req.ContainerId, req.ContainerName, "", model.StateBackup)
		container = monitor.GetContainer(req.ContainerId)
	}
	reply.Status = StatusOk
	container.HeartbeatCh <- struct{}{}
	log.Infof("container %s receive heartbeat", container.Id)
	return reply, nil
}

func (s *CftServer) HandleSnapshot(ctx context.Context, req *pb.SnapshotRequest) (*pb.SnapshotReply, error) {
	file := utils.GetCheckpointFile(req.ContainerId)
	f, err := os.Create(file)
	if err != nil {
		log.Error(err)
	}
	defer f.Close()
	_, err = f.Write(req.Data)
	reply := &pb.SnapshotReply{
		ContainerId: req.ContainerId,
		StatusId:    StatusOk,
	}
	log.Infof("container %s receive snapshot", req.ContainerId)
	return reply, nil
}
