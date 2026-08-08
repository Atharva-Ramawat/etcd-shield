package main

import (
	"context"
	"log"
	"time"

	pb "go.etcd.io/etcd/api/v3/etcdserverpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewMaintenanceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err = client.Alarm(ctx, &pb.AlarmRequest{
		Action: pb.AlarmRequest_GET,
	})
	
	if err != nil {
		log.Printf("RESULT (SUCCESSFULLY BLOCKED!): %v", err)
	} else {
		log.Println("RESULT: WARNING - Malicious call was allowed!")
	}
}
