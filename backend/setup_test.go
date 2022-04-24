package backend

import (
	"log"
	"testing"

	"github.com/joho/godotenv"
)

const testTable = "test_transactions"

func TestMain(m *testing.M) {
	err := godotenv.Load("./tests.env")
	if err != nil {
		log.Fatal(err)
	}

	m.Run()
}
