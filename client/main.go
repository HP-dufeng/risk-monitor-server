package main

import (
	"context"
	"flag"
	"io"
	"log"
	"time"

	pb "github.com/fengdu/risk-monitor-server/pb"
	"google.golang.org/grpc"
)

func main() {
	flag.Parse()
	conn, err := grpc.Dial(*serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewRiskMonitorServerClient(conn)

	printContracts(client)
}

var (
	serverAddr = flag.String("server_addr", "10.1.7.127:8080", "The server address in the format of host:port")
)

// printContracts lists all the contracts.
func printContracts(client pb.RiskMonitorServerClient) {
	log.Printf("Looking for contracts within ")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.QryContract(ctx, &pb.QryContractReq{})
	if err != nil {
		log.Fatalf("%v.QryContract(_) = _, %v", client, err)
	}
	for {
		contract, err := stream.Recv()
		if err == io.EOF {
			log.Println("EOF")
			break
		}
		if err != nil {
			log.Fatalf("%v.QryContract(_) = _, %v", client, err)
		}
		log.Println(contract)
	}
}
