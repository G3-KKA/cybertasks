package main

import (
	"extservice/internal/extservice"
	"extservice/model"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

// Represents some external service with tasks.
// Have its own task storage.
// Can error.

func main() {

	//
	//
	viper.AutomaticEnv()
	viper.MustBindEnv("EXT_SERVICE_ADDR")
	addr := viper.GetString("EXT_SERVICE_ADDR")
	if addr == "" {
		log.Fatal("addr nil")
	}

	//
	//
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	//
	//
	server := grpc.NewServer()
	service := &service{
		tasks:                          make(map[model.TaskID]*model.Task, 1000),
		errcounter:                     atomic.Uint64{},
		mx:                             sync.Mutex{},
		UnimplementedTaskServiceServer: extservice.UnimplementedTaskServiceServer{},
	}
	extservice.RegisterTaskServiceServer(server, service)

	//
	//
	taskcreater := func() {
		ticker := time.NewTicker(time.Second * 6)
		for range ticker.C {
			service.newtask()
		}
	}
	go taskcreater()

	//
	//
	err = server.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}
