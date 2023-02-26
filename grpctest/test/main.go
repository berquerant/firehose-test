package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	grpchealth "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func main() {
	log.Fatal(run())
}

func run() error {
	ls, err := net.Listen("tcp", fmt.Sprintf(":%d", getPort()))
	if err != nil {
		return err
	}
	srv := grpc.NewServer()
	grpchealth.RegisterHealthServer(srv, health.NewServer())
	reflection.Register(srv)
	log.Printf("listening on %v", ls.Addr())
	return srv.Serve(ls)
}

func getPort() int {
	if p := os.Getenv("PORT"); p != "" {
		if i, err := strconv.Atoi(p); err == nil {
			return i
		}
	}
	return 17001
}
