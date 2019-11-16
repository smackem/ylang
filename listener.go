package main

import (
	"bytes"
	"context"
	"fmt"
	pb "github.com/smackem/ylang/internal/listener"
	"github.com/smackem/ylang/internal/program"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	listenerPort = ":50051"
)

type server struct {
	pb.UnimplementedImageProcServer
}

func (s *server) SayHello(ctx context.Context, in *pb.SimpleMsg) (*pb.SimpleMsg, error) {
	log.Print("Received SayHello")
	return &pb.SimpleMsg{
		Text: in.Text + " jupp",
	}, nil
}

func (s *server) ProcessImage(ctx context.Context, in *pb.ProcessImageRequest) (*pb.ProcessImageResponse, error) {
	surf, err := loadSurface(bytes.NewBuffer(in.ImageDataPng))
	if err != nil {
		return nil, fmt.Errorf("error decoding imageData: %s", err)
	}

	prog, err := program.Compile(string(in.SourceCode))
	if err != nil {
		return &pb.ProcessImageResponse{
			Result:       pb.ProcessImageResponse_ERROR,
			Message:      fmt.Sprintf("compilation error: %s", err),
			ImageDataPng: nil,
		}, nil
	}

	err = program.Execute(prog, surf)
	if err != nil {
		return &pb.ProcessImageResponse{
			Result:       pb.ProcessImageResponse_ERROR,
			Message:      fmt.Sprintf("execution error: %s", err),
			ImageDataPng: nil,
		}, nil
	}

	buf := bytes.Buffer{}
	err = writeImage(surf.target, &buf)
	if err != nil {
		return nil, fmt.Errorf("error encoding imageData : %s", err)
	}

	return &pb.ProcessImageResponse{
		Result:       pb.ProcessImageResponse_OK,
		Message:      "",
		ImageDataPng: buf.Bytes(),
	}, nil
}

func listenerMain() {
	lis, err := net.Listen("tcp", listenerPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterImageProcServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
