package backend

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"sync"
)

type FileTransactionLogger struct {
	events       chan<- Event
	errors       <-chan error
	lastSequence uint64
	file         *os.File
	wg           *sync.WaitGroup
}

func NewFileTransactionLogger(filename string) (TransactionLogger, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		return nil, fmt.Errorf("cannot open transaction log file: %w", err)
	}

	return &FileTransactionLogger{file: file, wg: &sync.WaitGroup{}}, nil
}

func (f *FileTransactionLogger) WritePut(key string, value string) {
	f.wg.Add(1)
	f.events <- Event{EventType: EventPut, Key: key, Value: value}
}

func (f *FileTransactionLogger) WriteDelete(key string) {
	f.wg.Add(1)
	f.events <- Event{EventType: EventDelete, Key: key}
}

func (f *FileTransactionLogger) Err() <-chan error {
	return f.errors
}

func (f *FileTransactionLogger) LastSequence() uint64 {
	return f.lastSequence
}

func (f *FileTransactionLogger) Run() {
	events := make(chan Event, 16)
	f.events = events

	errors := make(chan error, 1)
	f.errors = errors

	// events goroutine
	go func() {
		for e := range events {
			f.lastSequence++

			_, err := fmt.Fprintf(f.file, "%d\t%d\t%s\t%s\n", f.lastSequence, e.EventType, e.Key, e.Value)
			if err != nil {
				errors <- fmt.Errorf("cannot write to log file: %w", err)
			}

			f.wg.Done()
		}
	}()

	// errors goroutine
	go func() {
		for err := range errors {
			log.Println(err)
		}
	}()

}

func (f *FileTransactionLogger) Wait() {
	f.wg.Wait()
}

func (f *FileTransactionLogger) Close() error {
	f.wg.Wait()

	if f.events != nil {
		close(f.events)
	}

	return f.file.Close()
}

func (f *FileTransactionLogger) ReadEvents() (<-chan Event, <-chan error) {
	outEvent := make(chan Event)
	outError := make(chan error, 1)
	scanner := bufio.NewScanner(f.file)

	go func() {
		var e Event

		defer close(outEvent)
		defer close(outError)

		for scanner.Scan() {
			line := scanner.Text()

			fmt.Sscanf(line, "%d\t%d\t%s\t%s", &e.Sequence, &e.EventType, &e.Key, &e.Value)

			if f.lastSequence >= e.Sequence {
				outError <- fmt.Errorf("transaction numbers out of sequence")
				return
			}

			uv, err := url.QueryUnescape(e.Value)
			if err != nil {
				outError <- fmt.Errorf("value decoding failure: %w", err)
				return
			}

			e.Value = uv
			f.lastSequence = e.Sequence

			outEvent <- e
		}

		if err := scanner.Err(); err != nil {
			outError <- fmt.Errorf("transaction log read failure: %w", err)
		}
	}()

	return outEvent, outError
}
