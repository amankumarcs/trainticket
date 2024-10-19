package api

import (
	"log"
	"net"

	model "github.com/amankumarcs/trainticket/pkg/model/ticketing"
	"google.golang.org/grpc"
)

func StartServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	model.RegisterTicketServiceServer(grpcServer, NewTicketServiceServer())
	log.Println("Server is running on port 50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
