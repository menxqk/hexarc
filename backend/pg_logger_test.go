package backend

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"
)

const testTable = "test_transactions"

func TestNewPostgresTransactionLogger(t *testing.T) {
	os.Setenv("PG_HOST", "localhost")
	os.Setenv("PG_DBNAME", "test")
	os.Setenv("PG_USER", "test")
	os.Setenv("PG_PASSWORD", "test")

	tl, err := NewPostgresTransactionLogger(testTable)
	if err != nil {
		t.Error(err)
	}
	_, ok := tl.(*PostgresTransactionLogger)
	if !ok {
		t.Errorf("type mismatch for PostgresTransactionLogger: %v", reflect.TypeOf(tl))
	}
}

func TestPostgresReadEvents(t *testing.T) {
	db, err := getPsqlDb()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	err = dropTestTable(db)
	if err != nil {
		t.Fatal(err)
	}
	err = createTestTable(db)
	if err != nil {
		t.Fatal(err)
	}
	err = write4PostgresEvents(db)
	defer dropTestTable(db)

	tl1, err := NewPostgresTransactionLogger(testTable)
	if err != nil {
		t.Error(err)
	}
	chEv, chErr := tl1.ReadEvents()
	count, ok := 0, true
	for ok && err == nil {
		select {
		case err, ok = <-chErr:
		case _, ok = <-chEv:
			if ok {
				count++
			}
		}
	}
	if count != 4 {
		t.Errorf("4 events were expected, events read: %d", count)
	}
	tl1.Run()
	tl1.WritePut("quatro", "quarta") // Event 5
	tl1.WriteDelete("quatro")        // Event 6
	tl1.WritePut("cinco", "quinta")  // Event 7
	err = tl1.Close()
	if err != nil {
		t.Error(err)
	}

	tl2, err := NewPostgresTransactionLogger(testTable)
	if err != nil {
		t.Fatal(err)
	}
	defer tl2.Close()
	chEv, chErr = tl2.ReadEvents()
	count, ok = 0, true
	for ok && err == nil {
		select {
		case err, ok = <-chErr:
		case _, ok = <-chEv:
			if ok {
				count++
			}
		}
	}
	if count != 7 {
		t.Errorf("7 events were expected, events read: %d", count)
	}
}

func getPsqlDb() (*sql.DB, error) {
	host := "localhost"
	dbName := "test"
	user := "test"
	password := "test"

	connStr := fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=disable", host, dbName, user, password)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to create db value: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to open db connection: %w", err)
	}

	return db, nil
}

func dropTestTable(db *sql.DB) error {
	dropQuery := `DROP TABLE IF EXISTS %s;`
	dropQuery = fmt.Sprintf(dropQuery, testTable)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx, dropQuery)
	if err != nil {
		return err
	}

	return nil
}

func createTestTable(db *sql.DB) error {
	createQuery := `CREATE TABLE %s (
		sequence BIGSERIAL PRIMARY KEY,
		event_type SMALLINT,
		key TEXT, 
		value TEXT
	);
	`
	createQuery = fmt.Sprintf(createQuery, testTable)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx, createQuery)
	if err != nil {
		return err
	}

	return nil
}

func write4PostgresEvents(db *sql.DB) error {
	events := []Event{
		{0, EventPut, "uma", "primeira"},
		{0, EventPut, "duas", "segunda"},
		{0, EventPut, "três", "terceira"},
		{0, EventDelete, "uma", ""},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	insertQuery := `INSERT INTO %s (event_type, key, value) VALUES ($1, $2, $3);`
	insertQuery = fmt.Sprintf(insertQuery, testTable)

	for _, e := range events {
		_, err := db.ExecContext(ctx, insertQuery, e.EventType, e.Key, e.Value)
		if err != nil {
			return err
		}
	}

	return nil
}
