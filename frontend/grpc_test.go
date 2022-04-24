package frontend

import (
	"context"
	"testing"
	"time"

	pb "github.com/menxqk/hexarc/proto/v1"
)

func TestGrpcPut(t *testing.T) {
	const key = "a-key"
	const value = "a-value"

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resp, err := grpcClient.Put(ctx, &pb.PutRequest{Key: key, Value: value})
	if err != nil {
		t.Error(err)
	}
	if resp == nil {
		t.Error("should have returned a valid PutResponse")
	}
}

func TestGrpcGet(t *testing.T) {
	const key = "a-key"
	const value = "a-value"

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	respPut, err := grpcClient.Put(ctx, &pb.PutRequest{Key: key, Value: value})
	if err != nil {
		t.Error(err)
	}
	if respPut == nil {
		t.Error("should have returned a valid PutResponse")
	}

	respGet, err := grpcClient.Get(ctx, &pb.GetRequest{Key: key})
	if err != nil {
		t.Error(err)
	}
	if respGet == nil {
		t.Error("should have returned a valid GetResponse")
	}
	if respGet.Value != value {
		t.Errorf("respGet.Value/value mismatch, val: %q, value: %q", respGet.Value, value)
	}
}

func TestGrpcDelete(t *testing.T) {
	const key = "a-key"
	const value = "a-value"

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	respPut, err := grpcClient.Put(ctx, &pb.PutRequest{Key: key, Value: value})
	if err != nil {
		t.Error(err)
	}
	if respPut == nil {
		t.Error("should have returned a valid PutResponse")
	}

	respDelete, err := grpcClient.Delete(ctx, &pb.DeleteRequest{Key: key})
	if err != nil {
		t.Error(err)
	}
	if respDelete == nil {
		t.Error("should have returned a valid DeleteResponse")
	}

	respGet, err := grpcClient.Get(ctx, &pb.GetRequest{Key: key})
	if err == nil {
		t.Error("should have returned an error for trying to get a deleted key")
	}
	if respGet != nil {
		t.Error("should have returned a nil GetResponse")
	}
}
