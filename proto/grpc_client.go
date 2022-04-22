package main

import (
	"context"
	"log"
	"os"
	"time"

	pb "github.com/menxqk/hexarc/proto/v1"
	"google.golang.org/grpc"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatalln("Usage: [GET|PUT|DELETE] key value")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	opts := []grpc.DialOption{grpc.WithInsecure(), grpc.WithBlock()}
	conn, err := grpc.DialContext(ctx, "localhost:50051", opts...)
	if err != nil {
		log.Fatalf("did not connect: %w", err)
	}
	defer conn.Close()

	client := pb.NewKeyValueClient(conn)

	action, key := os.Args[1], os.Args[2]
	var value string
	if len(os.Args) > 3 {
		value = os.Args[3]
	}

	switch action {
	case "get":
		r, err := client.Get(ctx, &pb.GetRequest{Key: key})
		if err != nil {
			log.Fatalf("could not get value for key %s: %v\n", key, err)
		}
		log.Printf("Get %s returns: %s\n", key, r.Value)
	case "put":
		_, err := client.Put(ctx, &pb.PutRequest{Key: key, Value: value})
		if err != nil {
			log.Fatalf("could not put key %s: %v\n", key, err)
		}
		log.Printf("Put %s\n", key)
	case "delete":
		_, err := client.Delete(ctx, &pb.DeleteRequest{Key: key})
		if err != nil {
			log.Fatalf("could not delete key %s: %v\n", key, err)
		}
		log.Printf("Delete %s\n", key)
	default:
		log.Fatalln("Usage: [GET|PUT|DELETE] key value")
	}

}
