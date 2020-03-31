package main

import (
	"io"
	"log"
	"net"

	"github.com/grpc-playground/proto/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct{}

func (s *server) SayHello(gs pb.Greeter_SayHelloServer) error {
	for {
		in, err := gs.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Printf("failed to recv: %v", err)
			return err
		}

		gs.Send(&pb.HelloReply{Message: "Hello " + in.Name})
	}

	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
