package model

import (
	"bufio"
	"cft/config"
	"cft/log"
	pb "cft/proto"
	"cft/utils"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"os"
	"os/exec"
	"sync/atomic"
	"time"
)

type StateType int

const (
	StateMaster StateType = iota
	StateBackup
	StatusOk          int32 = iota
	HeartbeatTimeout        = 10 * time.Second
	HeartbeatInterval       = 200 * time.Millisecond
	SnapshotInterval        = 10 * time.Second
)

type Container struct {
	Id             string `form:"id" json:"id"`
	Name           string `form:"name" json:"name"`
	BackupHost     string `form:"host" json:"host"`
	HeartbeatCh    chan struct{}
	HeartbeatCount int
	State          StateType
	dead           int32
}

func (c *Container) SendHeartbeat() {
	log.Infof("container %s try to send heartbeat to %s", c.Id, c.BackupHost)
	if c.BackupHost == "" {
		return
	}
	conn, err := grpc.Dial(c.BackupHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Errorf("%v", err)
	}
	defer conn.Close()
	client := pb.NewCftServerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.HeartbeatRequest{
		ContainerId: c.Id,
	}
	reply, err := client.HandleHeartbeat(ctx, req)
	if err != nil {
		log.Errorf("%v", err)
	}
	c.HeartbeatCount++
	log.Infof("container %s, heartbeat reply %+v", c.Id, reply)
}

func (c *Container) SendSnapshot() {
	log.Infof("container %s try to send snapshot to %s", c.Id, c.BackupHost)
	conn, err := grpc.Dial(c.BackupHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Errorf("%v", err)
	}
	defer conn.Close()
	file := c.checkpoint()
	if file == "" {
		return
	}

	f, err := os.Open(file)
	if err != nil {
		log.Fatalf("container %s open checkpoint error %v", c.Id, err)
		return
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	data, err := io.ReadAll(reader)
	if err != nil {
		log.Error(err)
	}

	client := pb.NewCftServerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.SnapshotRequest{
		ContainerId: c.Id,
		Data:        data,
	}
	reply, err := client.HandleSnapshot(ctx, req)
	if err != nil {
		log.Errorf("%v", err)
	}
	log.Infof("container %s, snapshot reply %+v", c.Id, reply)
}

func (c *Container) checkpoint() string {
	log.Infof("container %s try to make checkpoint", c.Id)
	_, err := os.Stat(config.GetConfig().CheckpointDir)
	if err != nil && os.IsNotExist(err) {
		os.Mkdir(config.GetConfig().CheckpointDir, os.ModePerm)
	}
	file := utils.GetCheckpointFile(c.Id)
	cmd := exec.Command("podman", "container", "checkpoint", c.Id, "-e", file, "-R")
	err = cmd.Run()
	if err != nil {
		log.Errorf("container %s make checkpoint error %v", c.Id, err)
		return ""
	}
	return file
}

func (c *Container) becomeMaster() {
	c.State = StateMaster
	c.BackupHost = ""
	file := utils.GetCheckpointFile(c.Id)
	cmd := exec.Command("podman", "container", "restore", "-i", file, "--runtime", "runc")
	err := cmd.Run()
	if err != nil {
		log.Errorf("container %s make checkpoint error %v", c.Id, err)
	}
}

func (c *Container) run() {
	for !c.stopped() {
		switch c.State {
		case StateMaster:
			c.SendHeartbeat()
			c.HeartbeatCount++
			if time.Duration(c.HeartbeatCount)*HeartbeatInterval >= SnapshotInterval {
				c.SendSnapshot()
				c.HeartbeatCount = 0
			}
			time.Sleep(HeartbeatInterval)
		case StateBackup:
			select {
			case <-c.HeartbeatCh:
			case <-time.Tick(HeartbeatTimeout):
				c.becomeMaster()
			}
		}
	}
}

func (c *Container) stopped() bool {
	return atomic.LoadInt32(&c.dead) == 1
}

func (c *Container) Kill() {
	atomic.StoreInt32(&c.dead, 1)
}

func NewContainer(id string, name string, host string, state StateType) *Container {
	c := &Container{
		Id:          id,
		Name:        name,
		BackupHost:  host,
		HeartbeatCh: make(chan struct{}),
		State:       state,
	}
	go c.run()
	return c
}
