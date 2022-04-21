package backend

import "fmt"

func NewTransactionLogger(s string) (TransactionLogger, error) {
	switch s {
	case "file":
		return NewFileTransactionLogger("./transactions.txt")
	case "postgres":
		return NewPostgresTransactionLogger()
	default:
		return nil, fmt.Errorf("no such transaction logger %q", s)
	}
}
