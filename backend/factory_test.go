package backend

import (
	"testing"
)

func TestFactory(t *testing.T) {
	tl, err := NewTransactionLogger("")
	if err == nil {
		t.Error("should have returned an error for empty TransactionLogger")
	}
	if tl != nil {
		t.Error("should have returned nil")
	}
}
