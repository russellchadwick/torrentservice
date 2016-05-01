// Code generated by protoc-gen-go.
// source: torrent.proto
// DO NOT EDIT!

/*
Package torrent is a generated protocol buffer package.

It is generated from these files:
	torrent.proto

It has these top-level messages:
	AddTorrentRequest
	AddTorrentResponse
*/
package torrent

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
const _ = proto.ProtoPackageIsVersion1

type AddTorrentRequest struct {
	Url string `protobuf:"bytes,1,opt,name=url" json:"url,omitempty"`
}

func (m *AddTorrentRequest) Reset()                    { *m = AddTorrentRequest{} }
func (m *AddTorrentRequest) String() string            { return proto.CompactTextString(m) }
func (*AddTorrentRequest) ProtoMessage()               {}
func (*AddTorrentRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type AddTorrentResponse struct {
	Id   uint32 `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	Name string `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
	Hash string `protobuf:"bytes,3,opt,name=hash" json:"hash,omitempty"`
}

func (m *AddTorrentResponse) Reset()                    { *m = AddTorrentResponse{} }
func (m *AddTorrentResponse) String() string            { return proto.CompactTextString(m) }
func (*AddTorrentResponse) ProtoMessage()               {}
func (*AddTorrentResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func init() {
	proto.RegisterType((*AddTorrentRequest)(nil), "torrent.AddTorrentRequest")
	proto.RegisterType((*AddTorrentResponse)(nil), "torrent.AddTorrentResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion1

// Client API for Torrent service

type TorrentClient interface {
	AddTorrent(ctx context.Context, in *AddTorrentRequest, opts ...grpc.CallOption) (*AddTorrentResponse, error)
}

type torrentClient struct {
	cc *grpc.ClientConn
}

func NewTorrentClient(cc *grpc.ClientConn) TorrentClient {
	return &torrentClient{cc}
}

func (c *torrentClient) AddTorrent(ctx context.Context, in *AddTorrentRequest, opts ...grpc.CallOption) (*AddTorrentResponse, error) {
	out := new(AddTorrentResponse)
	err := grpc.Invoke(ctx, "/torrent.Torrent/AddTorrent", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Torrent service

type TorrentServer interface {
	AddTorrent(context.Context, *AddTorrentRequest) (*AddTorrentResponse, error)
}

func RegisterTorrentServer(s *grpc.Server, srv TorrentServer) {
	s.RegisterService(&_Torrent_serviceDesc, srv)
}

func _Torrent_AddTorrent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(AddTorrentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(TorrentServer).AddTorrent(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

var _Torrent_serviceDesc = grpc.ServiceDesc{
	ServiceName: "torrent.Torrent",
	HandlerType: (*TorrentServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddTorrent",
			Handler:    _Torrent_AddTorrent_Handler,
		},
	},
	Streams: []grpc.StreamDesc{},
}

var fileDescriptor0 = []byte{
	// 158 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0x2d, 0xc9, 0x2f, 0x2a,
	0x4a, 0xcd, 0x2b, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x87, 0x72, 0x95, 0x54, 0xb9,
	0x04, 0x1d, 0x53, 0x52, 0x42, 0x20, 0xbc, 0xa0, 0xd4, 0xc2, 0xd2, 0xd4, 0xe2, 0x12, 0x21, 0x01,
	0x2e, 0xe6, 0xd2, 0xa2, 0x1c, 0x09, 0x46, 0x05, 0x46, 0x0d, 0xce, 0x20, 0x10, 0x53, 0xc9, 0x87,
	0x4b, 0x08, 0x59, 0x59, 0x71, 0x41, 0x7e, 0x5e, 0x71, 0xaa, 0x10, 0x1f, 0x17, 0x53, 0x66, 0x0a,
	0x58, 0x19, 0x6f, 0x10, 0x90, 0x25, 0x24, 0xc4, 0xc5, 0x92, 0x97, 0x98, 0x9b, 0x2a, 0xc1, 0x04,
	0xd6, 0x08, 0x66, 0x83, 0xc4, 0x32, 0x12, 0x8b, 0x33, 0x24, 0x98, 0x21, 0x62, 0x20, 0xb6, 0x51,
	0x10, 0x17, 0x3b, 0xd4, 0x28, 0x21, 0x77, 0x2e, 0x2e, 0x84, 0xc1, 0x42, 0x52, 0x7a, 0x30, 0x67,
	0x62, 0x38, 0x4a, 0x4a, 0x1a, 0xab, 0x1c, 0xc4, 0x25, 0x4a, 0x0c, 0x49, 0x6c, 0x60, 0x8f, 0x19,
	0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0x2b, 0x99, 0x62, 0x4b, 0xe9, 0x00, 0x00, 0x00,
}
