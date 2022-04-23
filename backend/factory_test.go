package backend

import (
	"os"
	"reflect"
	"testing"
)

func TestFactory(t *testing.T) {
	tl, err := NewTransactionLogger("file")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove("./transactions.txt")
	_, ok := tl.(*FileTransactionLogger)
	if !ok {
		t.Errorf("type mismatch for FileTransactionLogger: %v", reflect.TypeOf(tl))
	}

	os.Setenv("PG_HOST", "localhost")
	os.Setenv("PG_DBNAME", "test")
	os.Setenv("PG_USER", "test")
	os.Setenv("PG_PASSWORD", "test")

	tl, err = NewTransactionLogger("postgres")
	if err != nil {
		t.Error(err)
	}
	_, ok = tl.(*PostgresTransactionLogger)
	if !ok {
		t.Errorf("type mismatch for PostgresTransactionLogger: %v", reflect.TypeOf(tl))
	}

	tl, err = NewTransactionLogger("")
	if err == nil {
		t.Error("should have returned an error for empty TransactionLogger")
	}
	if tl != nil {
		t.Error("should have returned nil")
	}

}
