.PHONY: run
run:
	go run main.go

.PHONY: test
test:
	go test -v ./...

.PHONY: protoc
protoc:
	protoc proto/v1/*.proto --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative --proto_path=.