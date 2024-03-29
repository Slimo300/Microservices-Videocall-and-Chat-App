// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.12.4
// source: tokens.proto

package auth

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	TokenService_NewPairFromUserID_FullMethodName  = "/tokens.TokenService/NewPairFromUserID"
	TokenService_NewPairFromRefresh_FullMethodName = "/tokens.TokenService/NewPairFromRefresh"
	TokenService_DeleteUserToken_FullMethodName    = "/tokens.TokenService/DeleteUserToken"
)

// TokenServiceClient is the client API for TokenService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TokenServiceClient interface {
	NewPairFromUserID(ctx context.Context, in *UserID, opts ...grpc.CallOption) (*TokenPair, error)
	NewPairFromRefresh(ctx context.Context, in *RefreshToken, opts ...grpc.CallOption) (*TokenPair, error)
	DeleteUserToken(ctx context.Context, in *RefreshToken, opts ...grpc.CallOption) (*Msg, error)
}

type tokenServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTokenServiceClient(cc grpc.ClientConnInterface) TokenServiceClient {
	return &tokenServiceClient{cc}
}

func (c *tokenServiceClient) NewPairFromUserID(ctx context.Context, in *UserID, opts ...grpc.CallOption) (*TokenPair, error) {
	out := new(TokenPair)
	err := c.cc.Invoke(ctx, TokenService_NewPairFromUserID_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tokenServiceClient) NewPairFromRefresh(ctx context.Context, in *RefreshToken, opts ...grpc.CallOption) (*TokenPair, error) {
	out := new(TokenPair)
	err := c.cc.Invoke(ctx, TokenService_NewPairFromRefresh_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tokenServiceClient) DeleteUserToken(ctx context.Context, in *RefreshToken, opts ...grpc.CallOption) (*Msg, error) {
	out := new(Msg)
	err := c.cc.Invoke(ctx, TokenService_DeleteUserToken_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TokenServiceServer is the server API for TokenService service.
// All implementations must embed UnimplementedTokenServiceServer
// for forward compatibility
type TokenServiceServer interface {
	NewPairFromUserID(context.Context, *UserID) (*TokenPair, error)
	NewPairFromRefresh(context.Context, *RefreshToken) (*TokenPair, error)
	DeleteUserToken(context.Context, *RefreshToken) (*Msg, error)
	mustEmbedUnimplementedTokenServiceServer()
}

// UnimplementedTokenServiceServer must be embedded to have forward compatible implementations.
type UnimplementedTokenServiceServer struct {
}

func (UnimplementedTokenServiceServer) NewPairFromUserID(context.Context, *UserID) (*TokenPair, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NewPairFromUserID not implemented")
}
func (UnimplementedTokenServiceServer) NewPairFromRefresh(context.Context, *RefreshToken) (*TokenPair, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NewPairFromRefresh not implemented")
}
func (UnimplementedTokenServiceServer) DeleteUserToken(context.Context, *RefreshToken) (*Msg, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteUserToken not implemented")
}
func (UnimplementedTokenServiceServer) mustEmbedUnimplementedTokenServiceServer() {}

// UnsafeTokenServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TokenServiceServer will
// result in compilation errors.
type UnsafeTokenServiceServer interface {
	mustEmbedUnimplementedTokenServiceServer()
}

func RegisterTokenServiceServer(s grpc.ServiceRegistrar, srv TokenServiceServer) {
	s.RegisterService(&TokenService_ServiceDesc, srv)
}

func _TokenService_NewPairFromUserID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokenServiceServer).NewPairFromUserID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TokenService_NewPairFromUserID_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokenServiceServer).NewPairFromUserID(ctx, req.(*UserID))
	}
	return interceptor(ctx, in, info, handler)
}

func _TokenService_NewPairFromRefresh_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RefreshToken)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokenServiceServer).NewPairFromRefresh(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TokenService_NewPairFromRefresh_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokenServiceServer).NewPairFromRefresh(ctx, req.(*RefreshToken))
	}
	return interceptor(ctx, in, info, handler)
}

func _TokenService_DeleteUserToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RefreshToken)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokenServiceServer).DeleteUserToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TokenService_DeleteUserToken_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokenServiceServer).DeleteUserToken(ctx, req.(*RefreshToken))
	}
	return interceptor(ctx, in, info, handler)
}

// TokenService_ServiceDesc is the grpc.ServiceDesc for TokenService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TokenService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "tokens.TokenService",
	HandlerType: (*TokenServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "NewPairFromUserID",
			Handler:    _TokenService_NewPairFromUserID_Handler,
		},
		{
			MethodName: "NewPairFromRefresh",
			Handler:    _TokenService_NewPairFromRefresh_Handler,
		},
		{
			MethodName: "DeleteUserToken",
			Handler:    _TokenService_DeleteUserToken_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "tokens.proto",
}
