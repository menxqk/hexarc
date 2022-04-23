package backend

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

type PostgresTransactionLogger struct {
	events    chan<- Event
	errors    <-chan error
	db        *sql.DB
	wg        *sync.WaitGroup
	tableName string
}

func NewPostgresTransactionLogger(table string) (TransactionLogger, error) {
	host := os.Getenv("PG_HOST")
	dbName := os.Getenv("PG_DBNAME")
	user := os.Getenv("PG_USER")
	password := os.Getenv("PG_PASSWORD")

	connStr := fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=disable", host, dbName, user, password)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to create db value: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to open db connection: %w", err)
	}

	tl := &PostgresTransactionLogger{db: db, wg: &sync.WaitGroup{}, tableName: table}

	exists, err := tl.verifyTableExists()
	if err != nil {
		return nil, fmt.Errorf("failed to verify table exists: %w", err)
	}
	if !exists {
		if err := tl.createTable(); err != nil {
			return nil, fmt.Errorf("failed to create table: %w", err)
		}
	}

	return tl, nil
}

func (p *PostgresTransactionLogger) WritePut(key string, value string) {
	p.wg.Add(1)
	p.events <- Event{EventType: EventPut, Key: key, Value: value}
}

func (p *PostgresTransactionLogger) WriteDelete(key string) {
	p.wg.Add(1)
	p.events <- Event{EventType: EventDelete, Key: key}
}

func (p *PostgresTransactionLogger) Err() <-chan error {
	return p.errors
}

func (p *PostgresTransactionLogger) LastSequence() uint64 {
	return 0
}

func (p *PostgresTransactionLogger) Run() {
	events := make(chan Event, 16)
	p.events = events

	errors := make(chan error, 1)
	p.errors = errors

	// events goroutine
	go func() {
		query := `INSERT INTO %s (event_type, key, value) VALUES ($1, $2, $3)`
		query = fmt.Sprintf(query, p.tableName)

		for e := range events {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			_, err := p.db.ExecContext(ctx, query, e.EventType, e.Key, e.Value)
			if err != nil {
				errors <- err
			}

			p.wg.Done()
		}
	}()

	// errors goroutine
	go func() {
		for err := range errors {
			log.Println(err)
		}
	}()

}

func (p *PostgresTransactionLogger) Wait() {
	p.wg.Wait()
}

func (p *PostgresTransactionLogger) Close() error {
	p.wg.Wait()

	if p.events != nil {
		close(p.events)
	}

	return p.db.Close()
}

func (p *PostgresTransactionLogger) ReadEvents() (<-chan Event, <-chan error) {
	outEvent := make(chan Event)
	outError := make(chan error, 1)

	query := `SELECT sequence, event_type, key, value FROM %s`
	query = fmt.Sprintf(query, p.tableName)

	go func() {
		defer close(outEvent)
		defer close(outError)

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		rows, err := p.db.QueryContext(ctx, query)
		if err != nil {
			outError <- fmt.Errorf("sql query error: %w", err)
			return
		}

		defer rows.Close()

		var e Event

		for rows.Next() {
			err = rows.Scan(&e.Sequence, &e.EventType, &e.Key, &e.Value)
			if err != nil {
				outError <- err
				return
			}

			outEvent <- e
		}

		err = rows.Err()
		if err != nil {
			outError <- fmt.Errorf("transaction log read failure: %w", err)
		}

	}()

	return outEvent, outError
}

func (p *PostgresTransactionLogger) verifyTableExists() (bool, error) {
	var result string

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := p.db.QueryContext(ctx, fmt.Sprintf("SELECT to_regclass('public.%s');", p.tableName))
	defer rows.Close()
	if err != nil {
		return false, err
	}

	for rows.Next() && result != p.tableName {
		rows.Scan(&result)
	}

	return result == p.tableName, rows.Err()
}

func (p *PostgresTransactionLogger) createTable() error {
	var err error

	createQuery := `CREATE TABLE %s (
		sequence BIGSERIAL PRIMARY KEY,
		event_type SMALLINT,
		key TEXT, 
		value TEXT
	);
	`
	createQuery = fmt.Sprintf(createQuery, p.tableName)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err = p.db.ExecContext(ctx, createQuery)
	if err != nil {
		return err
	}

	return nil
}
