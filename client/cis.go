package client

import (
	"context"

	pb "github.com/ekspand/trusty/api/v1/pb"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
)

// CertInfoClient client interface
type CertInfoClient interface {
	// Roots returns the root CAs
	Roots(ctx context.Context, in *empty.Empty) (*pb.RootsResponse, error)
}

type cisClient struct {
	remote   pb.CertInfoServiceClient
	callOpts []grpc.CallOption
}

// NewCertInfo returns instance of CertInfoService client
func NewCertInfo(conn *grpc.ClientConn, callOpts []grpc.CallOption) CertInfoClient {
	return &cisClient{
		remote:   RetryCertInfoClient(conn),
		callOpts: callOpts,
	}
}

// NewCertInfoFromProxy returns instance of CertInfoService client
func NewCertInfoFromProxy(proxy pb.CertInfoServiceClient) CertInfoClient {
	return &cisClient{
		remote: proxy,
	}
}

// Roots returns the root CAs
func (c *cisClient) Roots(ctx context.Context, in *empty.Empty) (*pb.RootsResponse, error) {
	return c.remote.Roots(ctx, in, c.callOpts...)
}

type retryCertInfoClient struct {
	cis pb.CertInfoServiceClient
}

// TODO: implement retry for gRPC client interceptor

// RetryCertInfoClient implements a CertInfoServiceClient.
func RetryCertInfoClient(conn *grpc.ClientConn) pb.CertInfoServiceClient {
	return &retryCertInfoClient{
		cis: pb.NewCertInfoServiceClient(conn),
	}
}

// Roots returns the root CAs
func (c *retryCertInfoClient) Roots(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*pb.RootsResponse, error) {
	return c.cis.Roots(ctx, in, opts...)
}
