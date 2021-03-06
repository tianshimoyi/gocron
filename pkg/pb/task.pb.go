// Code generated by protoc-gen-go. DO NOT EDIT.
// source: task.proto

package pb

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

type TaskRequest struct {
	Command              string   `protobuf:"bytes,2,opt,name=command,proto3" json:"command,omitempty"`
	Timeout              int32    `protobuf:"varint,3,opt,name=timeout,proto3" json:"timeout,omitempty"`
	Id                   int64    `protobuf:"varint,4,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TaskRequest) Reset()         { *m = TaskRequest{} }
func (m *TaskRequest) String() string { return proto.CompactTextString(m) }
func (*TaskRequest) ProtoMessage()    {}
func (*TaskRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce5d8dd45b4a91ff, []int{0}
}

func (m *TaskRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TaskRequest.Unmarshal(m, b)
}
func (m *TaskRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TaskRequest.Marshal(b, m, deterministic)
}
func (m *TaskRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TaskRequest.Merge(m, src)
}
func (m *TaskRequest) XXX_Size() int {
	return xxx_messageInfo_TaskRequest.Size(m)
}
func (m *TaskRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_TaskRequest.DiscardUnknown(m)
}

var xxx_messageInfo_TaskRequest proto.InternalMessageInfo

func (m *TaskRequest) GetCommand() string {
	if m != nil {
		return m.Command
	}
	return ""
}

func (m *TaskRequest) GetTimeout() int32 {
	if m != nil {
		return m.Timeout
	}
	return 0
}

func (m *TaskRequest) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

type TaskResponse struct {
	Output               string   `protobuf:"bytes,1,opt,name=output,proto3" json:"output,omitempty"`
	Error                string   `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TaskResponse) Reset()         { *m = TaskResponse{} }
func (m *TaskResponse) String() string { return proto.CompactTextString(m) }
func (*TaskResponse) ProtoMessage()    {}
func (*TaskResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce5d8dd45b4a91ff, []int{1}
}

func (m *TaskResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TaskResponse.Unmarshal(m, b)
}
func (m *TaskResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TaskResponse.Marshal(b, m, deterministic)
}
func (m *TaskResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TaskResponse.Merge(m, src)
}
func (m *TaskResponse) XXX_Size() int {
	return xxx_messageInfo_TaskResponse.Size(m)
}
func (m *TaskResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_TaskResponse.DiscardUnknown(m)
}

var xxx_messageInfo_TaskResponse proto.InternalMessageInfo

func (m *TaskResponse) GetOutput() string {
	if m != nil {
		return m.Output
	}
	return ""
}

func (m *TaskResponse) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

func init() {
	proto.RegisterType((*TaskRequest)(nil), "rpc.TaskRequest")
	proto.RegisterType((*TaskResponse)(nil), "rpc.TaskResponse")
}

func init() { proto.RegisterFile("task.proto", fileDescriptor_ce5d8dd45b4a91ff) }

var fileDescriptor_ce5d8dd45b4a91ff = []byte{
	// 188 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x8f, 0x3d, 0x0b, 0xc2, 0x30,
	0x10, 0x86, 0x4d, 0xd2, 0x56, 0x3c, 0x45, 0x34, 0x88, 0x04, 0xa7, 0xd2, 0xa9, 0x83, 0x74, 0x50,
	0x47, 0x27, 0xff, 0x81, 0xc1, 0xc9, 0xad, 0x1f, 0x19, 0x4a, 0x69, 0x13, 0x93, 0xcb, 0xff, 0x97,
	0x7e, 0x41, 0xc7, 0xe7, 0x8e, 0x7b, 0x9f, 0x7b, 0x01, 0x30, 0x77, 0x4d, 0x66, 0xac, 0x46, 0xcd,
	0x99, 0x35, 0x65, 0xf2, 0x86, 0xed, 0x27, 0x77, 0x8d, 0x54, 0x3f, 0xaf, 0x1c, 0x72, 0x01, 0xeb,
	0x52, 0xb7, 0x6d, 0xde, 0x55, 0x82, 0xc6, 0x24, 0xdd, 0xc8, 0x19, 0xfb, 0x0d, 0xd6, 0xad, 0xd2,
	0x1e, 0x05, 0x8b, 0x49, 0x1a, 0xca, 0x19, 0xf9, 0x1e, 0x68, 0x5d, 0x89, 0x20, 0x26, 0x29, 0x93,
	0xb4, 0xae, 0x92, 0x27, 0xec, 0xc6, 0x48, 0x67, 0x74, 0xe7, 0x14, 0x3f, 0x43, 0xa4, 0x3d, 0x1a,
	0x8f, 0x82, 0x0c, 0x91, 0x13, 0xf1, 0x13, 0x84, 0xca, 0x5a, 0x6d, 0x27, 0xd3, 0x08, 0xb7, 0x07,
	0x04, 0xfd, 0x35, 0xbf, 0x02, 0x93, 0xbe, 0xe3, 0x87, 0xcc, 0x9a, 0x32, 0x5b, 0xbc, 0x78, 0x39,
	0x2e, 0x26, 0xa3, 0x21, 0x59, 0xbd, 0x82, 0x2f, 0x35, 0x45, 0x11, 0x0d, 0xc5, 0xee, 0xff, 0x00,
	0x00, 0x00, 0xff, 0xff, 0x5c, 0xda, 0x51, 0x9e, 0xe6, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// TaskClient is the client API for Task service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type TaskClient interface {
	Run(ctx context.Context, in *TaskRequest, opts ...grpc.CallOption) (*TaskResponse, error)
}

type taskClient struct {
	cc grpc.ClientConnInterface
}

func NewTaskClient(cc grpc.ClientConnInterface) TaskClient {
	return &taskClient{cc}
}

func (c *taskClient) Run(ctx context.Context, in *TaskRequest, opts ...grpc.CallOption) (*TaskResponse, error) {
	out := new(TaskResponse)
	err := c.cc.Invoke(ctx, "/rpc.Task/Run", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TaskServer is the server API for Task service.
type TaskServer interface {
	Run(context.Context, *TaskRequest) (*TaskResponse, error)
}

// UnimplementedTaskServer can be embedded to have forward compatible implementations.
type UnimplementedTaskServer struct {
}

func (*UnimplementedTaskServer) Run(ctx context.Context, req *TaskRequest) (*TaskResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Run not implemented")
}

func RegisterTaskServer(s *grpc.Server, srv TaskServer) {
	s.RegisterService(&_Task_serviceDesc, srv)
}

func _Task_Run_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskServer).Run(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.Task/Run",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskServer).Run(ctx, req.(*TaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Task_serviceDesc = grpc.ServiceDesc{
	ServiceName: "rpc.Task",
	HandlerType: (*TaskServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Run",
			Handler:    _Task_Run_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "task.proto",
}
