package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/menxqk/hexarc/backend"
	"github.com/menxqk/hexarc/core"
	"github.com/menxqk/hexarc/frontend"
)

func main() {
	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	chError := make(chan error)

	// Create backend
	tl, err := backend.NewTransactionLogger("file")
	if err != nil {
		panic(err)
	}

	// Create core
	store := core.NewKeyValueStore().WithTransactionLogger(tl)
	if err := store.Init(); err != nil {
		panic(err)
	}

	// Create frontends
	feRest, err := frontend.NewFrontEnd("rest")
	if err != nil {
		panic(err)
	}
	feGrpc, err := frontend.NewFrontEnd("grpc")
	if err != nil {
		panic(err)
	}
	feWebserver, err := frontend.NewFrontEnd("webserver")
	if err != nil {
		panic(err)
	}

	// Start frontends
	startFrontEnds(store, chError, feRest, feGrpc, feWebserver)

	// Block on error channel waiting for and error coming from the frontends
	log.Fatal(<-chError)
}

// startFrontEnds starts each frontend on its own goroutine
func startFrontEnds(store *core.KeyValueStore, errorChannel chan<- error, frontEnds ...frontend.FrontEnd) {
	for _, frontEnd := range frontEnds {
		go func(fe frontend.FrontEnd) {
			errorChannel <- fe.Start(store)
		}(frontEnd)
	}
}
