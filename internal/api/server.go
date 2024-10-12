package server

import (
	"fmt"
	"net"

	pb "github.com/anandu86130/Hospital-user-service/internal/pb"
	"google.golang.org/grpc"
)

type Server struct {
	server  *grpc.Server
	listner net.Listener
}

func NewGRPCServer(server pb.UserServer) (*Server, error) {
	lis, err := net.Listen("tcp", ":8082")
	if err != nil {
		return nil, err
	}
	newServer := grpc.NewServer()
	pb.RegisterUserServer(newServer, server)
	return &Server{
		server:  newServer,
		listner: lis,
	}, nil
}
func (c *Server) Start() error {
	fmt.Println("grpc server listening on port :50051")
	return c.server.Serve(c.listner)
}
