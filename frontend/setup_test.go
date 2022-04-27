package frontend

import (
	"context"
	"log"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/menxqk/hexarc/core"
	pb "github.com/menxqk/hexarc/proto/v1"
	"google.golang.org/grpc"
)

var (
	testKeyValueStore *core.KeyValueStore = core.NewKeyValueStore()

	testGrpcFrontEnd *grpcFrontEnd
	grpcPort         string
	grpcClient       pb.KeyValueClient

	testRestFrontend *restFrontEnd
	restPort         string

	testWebserverFrontend *webserverFrontEnd
	webserverPort         string
)

func TestMain(m *testing.M) {
	err := godotenv.Load("tests.env")
	if err != nil {
		log.Fatal(err)
	}

	// GRPC server and client setup
	grpcPort = os.Getenv("GRPC_PORT")
	feGrpc, err := NewFrontEnd("grpc")
	if err != nil {
		log.Fatal(err)
	}
	var ok bool
	testGrpcFrontEnd, ok = feGrpc.(*grpcFrontEnd)
	if !ok {
		log.Fatalf("type mismatch for grpcFrontEnd: %v", reflect.TypeOf(feGrpc))
	}
	// Start grpc server on its own goroutine
	go func() {
		testGrpcFrontEnd.Start(testKeyValueStore)
	}()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	opts := []grpc.DialOption{grpc.WithInsecure(), grpc.WithBlock()}
	conn, err := grpc.DialContext(ctx, "localhost:"+grpcPort, opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	grpcClient = pb.NewKeyValueClient(conn)

	// REST server setup
	restPort = os.Getenv("REST_PORT")
	feRest, err := NewFrontEnd("rest")
	if err != nil {
		log.Fatal(err)
	}
	testRestFrontend, ok = feRest.(*restFrontEnd)
	if !ok {
		log.Fatalf("type mismatch for restFrontEnd: %v", reflect.TypeOf(feRest))
	}
	// Start rest server on its own goroutine
	go func() {
		testRestFrontend.Start(testKeyValueStore)
	}()

	// WEBSERVER setup
	webserverPort = os.Getenv("WEBSERVER_PORT")
	feWebserver, err := NewFrontEnd("webserver")
	if err != nil {
		log.Fatal(err)
	}
	testWebserverFrontend, ok = feWebserver.(*webserverFrontEnd)
	if !ok {
		log.Fatalf("type mismatch for webserverFrontEnd: %v", reflect.TypeOf(feWebserver))
	}
	// Start webserver on its own goroutine
	go func() {
		testWebserverFrontend.Start(testKeyValueStore)
	}()

	m.Run()
}
