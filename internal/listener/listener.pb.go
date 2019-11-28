// Code generated by protoc-gen-go. DO NOT EDIT.
// source: listener.proto

package listener

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type ProcessImageResponse_CompilationResult int32

const (
	ProcessImageResponse_UNDEFINED ProcessImageResponse_CompilationResult = 0
	ProcessImageResponse_OK        ProcessImageResponse_CompilationResult = 1
	ProcessImageResponse_ERROR     ProcessImageResponse_CompilationResult = 2
)

var ProcessImageResponse_CompilationResult_name = map[int32]string{
	0: "UNDEFINED",
	1: "OK",
	2: "ERROR",
}

var ProcessImageResponse_CompilationResult_value = map[string]int32{
	"UNDEFINED": 0,
	"OK":        1,
	"ERROR":     2,
}

func (x ProcessImageResponse_CompilationResult) String() string {
	return proto.EnumName(ProcessImageResponse_CompilationResult_name, int32(x))
}

func (ProcessImageResponse_CompilationResult) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_f75aade3a9f7de9c, []int{1, 0}
}

type ProcessImageRequest struct {
	SourceCode           string   `protobuf:"bytes,1,opt,name=sourceCode,proto3" json:"sourceCode,omitempty"`
	ImageDataPng         []byte   `protobuf:"bytes,2,opt,name=imageDataPng,proto3" json:"imageDataPng,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ProcessImageRequest) Reset()         { *m = ProcessImageRequest{} }
func (m *ProcessImageRequest) String() string { return proto.CompactTextString(m) }
func (*ProcessImageRequest) ProtoMessage()    {}
func (*ProcessImageRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_f75aade3a9f7de9c, []int{0}
}

func (m *ProcessImageRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProcessImageRequest.Unmarshal(m, b)
}
func (m *ProcessImageRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProcessImageRequest.Marshal(b, m, deterministic)
}
func (m *ProcessImageRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProcessImageRequest.Merge(m, src)
}
func (m *ProcessImageRequest) XXX_Size() int {
	return xxx_messageInfo_ProcessImageRequest.Size(m)
}
func (m *ProcessImageRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ProcessImageRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ProcessImageRequest proto.InternalMessageInfo

func (m *ProcessImageRequest) GetSourceCode() string {
	if m != nil {
		return m.SourceCode
	}
	return ""
}

func (m *ProcessImageRequest) GetImageDataPng() []byte {
	if m != nil {
		return m.ImageDataPng
	}
	return nil
}

type ProcessImageResponse struct {
	Result               ProcessImageResponse_CompilationResult `protobuf:"varint,1,opt,name=result,proto3,enum=listener.ProcessImageResponse_CompilationResult" json:"result,omitempty"`
	Message              string                                 `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	ImageDataPng         []byte                                 `protobuf:"bytes,3,opt,name=imageDataPng,proto3" json:"imageDataPng,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                               `json:"-"`
	XXX_unrecognized     []byte                                 `json:"-"`
	XXX_sizecache        int32                                  `json:"-"`
}

func (m *ProcessImageResponse) Reset()         { *m = ProcessImageResponse{} }
func (m *ProcessImageResponse) String() string { return proto.CompactTextString(m) }
func (*ProcessImageResponse) ProtoMessage()    {}
func (*ProcessImageResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_f75aade3a9f7de9c, []int{1}
}

func (m *ProcessImageResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProcessImageResponse.Unmarshal(m, b)
}
func (m *ProcessImageResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProcessImageResponse.Marshal(b, m, deterministic)
}
func (m *ProcessImageResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProcessImageResponse.Merge(m, src)
}
func (m *ProcessImageResponse) XXX_Size() int {
	return xxx_messageInfo_ProcessImageResponse.Size(m)
}
func (m *ProcessImageResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ProcessImageResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ProcessImageResponse proto.InternalMessageInfo

func (m *ProcessImageResponse) GetResult() ProcessImageResponse_CompilationResult {
	if m != nil {
		return m.Result
	}
	return ProcessImageResponse_UNDEFINED
}

func (m *ProcessImageResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *ProcessImageResponse) GetImageDataPng() []byte {
	if m != nil {
		return m.ImageDataPng
	}
	return nil
}

func init() {
	proto.RegisterEnum("listener.ProcessImageResponse_CompilationResult", ProcessImageResponse_CompilationResult_name, ProcessImageResponse_CompilationResult_value)
	proto.RegisterType((*ProcessImageRequest)(nil), "listener.ProcessImageRequest")
	proto.RegisterType((*ProcessImageResponse)(nil), "listener.ProcessImageResponse")
}

func init() { proto.RegisterFile("listener.proto", fileDescriptor_f75aade3a9f7de9c) }

var fileDescriptor_f75aade3a9f7de9c = []byte{
	// 288 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x91, 0x41, 0x4b, 0xc3, 0x40,
	0x10, 0x85, 0xbb, 0x15, 0xab, 0x19, 0x6b, 0xa9, 0xab, 0x87, 0x50, 0xb0, 0x94, 0x9c, 0xe2, 0x25,
	0x94, 0x8a, 0x7f, 0xa0, 0x4d, 0xc4, 0xa2, 0xb4, 0x61, 0xc5, 0x43, 0x6f, 0xae, 0x71, 0x08, 0xc1,
	0x64, 0x37, 0x66, 0x36, 0x07, 0x7f, 0xa8, 0xff, 0x47, 0xb2, 0xd8, 0xd2, 0xd2, 0xea, 0x71, 0x86,
	0x79, 0x8f, 0xef, 0xcd, 0x83, 0x5e, 0x9e, 0x91, 0x41, 0x85, 0x55, 0x50, 0x56, 0xda, 0x68, 0x7e,
	0xba, 0x9e, 0xbd, 0x15, 0x5c, 0xc6, 0x95, 0x4e, 0x90, 0x68, 0x5e, 0xc8, 0x14, 0x05, 0x7e, 0xd6,
	0x48, 0x86, 0x0f, 0x01, 0x48, 0xd7, 0x55, 0x82, 0x33, 0xfd, 0x8e, 0x2e, 0x1b, 0x31, 0xdf, 0x11,
	0x5b, 0x1b, 0xee, 0x41, 0x37, 0x6b, 0xee, 0x43, 0x69, 0x64, 0xac, 0x52, 0xb7, 0x3d, 0x62, 0x7e,
	0x57, 0xec, 0xec, 0xbc, 0x6f, 0x06, 0x57, 0xbb, 0xde, 0x54, 0x6a, 0x45, 0xc8, 0x1f, 0xa0, 0x53,
	0x21, 0xd5, 0xb9, 0xb1, 0xc6, 0xbd, 0xc9, 0x38, 0xd8, 0xe0, 0x1d, 0xba, 0x0f, 0x66, 0xba, 0x28,
	0xb3, 0x5c, 0x9a, 0x4c, 0x2b, 0x61, 0x75, 0xe2, 0x57, 0xcf, 0x5d, 0x38, 0x29, 0x90, 0x48, 0xa6,
	0x68, 0x09, 0x1c, 0xb1, 0x1e, 0xf7, 0x00, 0x8f, 0x0e, 0x00, 0xde, 0xc1, 0xc5, 0x9e, 0x35, 0x3f,
	0x07, 0xe7, 0x65, 0x11, 0x46, 0xf7, 0xf3, 0x45, 0x14, 0xf6, 0x5b, 0xbc, 0x03, 0xed, 0xe5, 0x63,
	0x9f, 0x71, 0x07, 0x8e, 0x23, 0x21, 0x96, 0xa2, 0xdf, 0x9e, 0xbc, 0x82, 0x63, 0xf9, 0x1a, 0x56,
	0xfe, 0x0c, 0xdd, 0x6d, 0x66, 0x7e, 0xfd, 0x57, 0x16, 0xfb, 0xd7, 0xc1, 0xf0, 0xff, 0xa8, 0x5e,
	0xcb, 0x67, 0x63, 0x36, 0xbd, 0x81, 0x81, 0x42, 0x13, 0x50, 0x21, 0x93, 0x0f, 0x2c, 0x82, 0xaf,
	0x5c, 0xaa, 0x74, 0x23, 0x9c, 0x9e, 0xad, 0x9e, 0xa4, 0x4a, 0xe3, 0xa6, 0x48, 0x7a, 0xeb, 0xd8,
	0x42, 0x6f, 0x7f, 0x02, 0x00, 0x00, 0xff, 0xff, 0x4a, 0x08, 0x3f, 0x47, 0xe2, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ImageProcClient is the client API for ImageProc service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ImageProcClient interface {
	ProcessImage(ctx context.Context, opts ...grpc.CallOption) (ImageProc_ProcessImageClient, error)
}

type imageProcClient struct {
	cc *grpc.ClientConn
}

func NewImageProcClient(cc *grpc.ClientConn) ImageProcClient {
	return &imageProcClient{cc}
}

func (c *imageProcClient) ProcessImage(ctx context.Context, opts ...grpc.CallOption) (ImageProc_ProcessImageClient, error) {
	stream, err := c.cc.NewStream(ctx, &_ImageProc_serviceDesc.Streams[0], "/listener.ImageProc/ProcessImage", opts...)
	if err != nil {
		return nil, err
	}
	x := &imageProcProcessImageClient{stream}
	return x, nil
}

type ImageProc_ProcessImageClient interface {
	Send(*ProcessImageRequest) error
	Recv() (*ProcessImageResponse, error)
	grpc.ClientStream
}

type imageProcProcessImageClient struct {
	grpc.ClientStream
}

func (x *imageProcProcessImageClient) Send(m *ProcessImageRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *imageProcProcessImageClient) Recv() (*ProcessImageResponse, error) {
	m := new(ProcessImageResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ImageProcServer is the server API for ImageProc service.
type ImageProcServer interface {
	ProcessImage(ImageProc_ProcessImageServer) error
}

// UnimplementedImageProcServer can be embedded to have forward compatible implementations.
type UnimplementedImageProcServer struct {
}

func (*UnimplementedImageProcServer) ProcessImage(srv ImageProc_ProcessImageServer) error {
	return status.Errorf(codes.Unimplemented, "method ProcessImage not implemented")
}

func RegisterImageProcServer(s *grpc.Server, srv ImageProcServer) {
	s.RegisterService(&_ImageProc_serviceDesc, srv)
}

func _ImageProc_ProcessImage_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ImageProcServer).ProcessImage(&imageProcProcessImageServer{stream})
}

type ImageProc_ProcessImageServer interface {
	Send(*ProcessImageResponse) error
	Recv() (*ProcessImageRequest, error)
	grpc.ServerStream
}

type imageProcProcessImageServer struct {
	grpc.ServerStream
}

func (x *imageProcProcessImageServer) Send(m *ProcessImageResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *imageProcProcessImageServer) Recv() (*ProcessImageRequest, error) {
	m := new(ProcessImageRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _ImageProc_serviceDesc = grpc.ServiceDesc{
	ServiceName: "listener.ImageProc",
	HandlerType: (*ImageProcServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ProcessImage",
			Handler:       _ImageProc_ProcessImage_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "listener.proto",
}