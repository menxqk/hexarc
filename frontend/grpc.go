package frontend

import (
	"context"
	"fmt"
	"net"

	"github.com/menxqk/hexarc/core"
	"google.golang.org/grpc"

	pb "github.com/menxqk/hexarc/proto/v1"
)

type grpcFrontEnd struct {
	pb.UnimplementedKeyValueServer
	store *core.KeyValueStore
}

func (g *grpcFrontEnd) Start(kv *core.KeyValueStore) error {
	g.store = kv

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	s := grpc.NewServer()

	pb.RegisterKeyValueServer(s, g)

	return s.Serve(lis)
}

func (g *grpcFrontEnd) Get(ctx context.Context, r *pb.GetRequest) (*pb.GetResponse, error) {
	value, err := g.store.Get(r.Key)
	if err != nil {
		return nil, err
	}

	return &pb.GetResponse{Value: value}, nil
}

func (g *grpcFrontEnd) Put(ctx context.Context, r *pb.PutRequest) (*pb.PutResponse, error) {
	err := g.store.Put(r.Key, r.Value)
	if err != nil {
		return nil, err
	}

	return &pb.PutResponse{}, nil
}

func (g *grpcFrontEnd) Delete(ctx context.Context, r *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	err := g.store.Delete(r.Key)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteResponse{}, nil
}
