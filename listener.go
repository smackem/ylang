package main

import (
	"bytes"
	"context"
	"fmt"
	pb "github.com/smackem/ylang/internal/listener"
	"github.com/smackem/ylang/internal/program"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)

const (
	listenerPort = ":50051"
)

type server struct {
	pb.UnimplementedImageProcServer
}

func (s *server) ProcessImage(srv pb.ImageProc_ProcessImageServer) error {
	request, err := s.readRequest(srv)
	if err != nil {
		return err
	}
	response, err := s.processImage(srv.Context(), request)
	if err != nil {
		return err
	}
	return s.writeResponse(response, srv)
}

func (s *server) writeResponse(response *pb.ProcessImageResponse, srv pb.ImageProc_ProcessImageServer) error {
	resp := pb.ProcessImageResponse{}
	const chunkSize = 64 * 1024
	index := 0
	isFirstMessage := true
	for remaining := len(response.ImageDataPng); remaining > 0 || isFirstMessage; {
		toWrite := chunkSize
		if toWrite > remaining {
			toWrite = remaining
		}
		resp.ImageDataPng = response.ImageDataPng[index : index+toWrite]
		if isFirstMessage {
			resp.Result = response.Result
			resp.Message = response.Message
			isFirstMessage = false
		}
		if err := srv.Send(&resp); err != nil {
			return err
		}
		index += toWrite
		remaining -= toWrite
	}
	return nil
}

func (s *server) readRequest(srv pb.ImageProc_ProcessImageServer) (*pb.ProcessImageRequest, error) {
	var fullRequest pb.ProcessImageRequest
	first := true

	for {
		request, err := srv.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		if first {
			fullRequest.SourceCode = request.SourceCode
		}
		fullRequest.ImageDataPng = append(fullRequest.ImageDataPng, request.ImageDataPng...)
		first = false
	}

	return &fullRequest, nil
}

func (s *server) processImage(ctx context.Context, in *pb.ProcessImageRequest) (*pb.ProcessImageResponse, error) {
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
