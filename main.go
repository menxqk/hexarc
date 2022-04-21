package main

import (
	"github.com/menxqk/hexarc/backend"
	"github.com/menxqk/hexarc/core"
)

func main() {
	// Create backend
	tl, err := backend.NewTransactionLogger("postgres")
	if err != nil {
		panic(err)
	}

	// Create core
	store := core.NewKeyValueStore().WithTransactionLogger(tl)
	if err := store.Init(); err != nil {
		panic(err)
	}

	// TODO - Create frontend
}
