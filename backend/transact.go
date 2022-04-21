package backend

type TransactionLogger interface {
	WritePut(string, string)
	WriteDelete(string)
	Err() <-chan error
	LastSequence() uint64
	Run()
	Wait()
	Close() error
	ReadEvents() (<-chan Event, <-chan error)
}

type Event struct {
	Sequence  uint64
	EventType EventType
	Key       string
	Value     string
}

type EventType byte

const (
	_                  = iota
	EventPut EventType = iota
	EventDelete
)
