package postgres

import (
	assert "github.com/stretchr/testify/require"
	"testing"

	"context"
	"database/sql"
	"time"
)

func TestInit(t *testing.T) {

	// Throw error if db does not exist
	dbName := "notexists"
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	db, err := TestDB(ctx, dbName)
	assert.Equal(t, "pq: database \"notexists\" does not exist", err.Error())

	// Access with no db
	dbName = ""
	db, err = TestDB(ctx, dbName)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Reset the test table
	_, err = db.Exec("DROP TABLE IF EXISTS testtable1")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS testtable1 (id serial PRIMARY KEY, filename TEXT)")
	if err != nil {
		t.Fatal(err)
	}

	// Confirm that there are no rows in the table
	var results int64
	err = db.Get(&results, "SELECT id FROM testtable1")
	assert.Equal(t, sql.ErrNoRows, err)

	// Insert a row in the table
	_, err = db.Exec("INSERT INTO testtable1 (filename) VALUES ($1)", "test.jpg")
	if err != nil {
		t.Fatal(err)
	}

	// Select this new row
	err = db.Get(&results, "SELECT id FROM testtable1")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, int64(1), results)

	// Drop the test table
	_, err = db.Exec("DROP TABLE IF EXISTS testtable1")
	if err != nil {
		t.Fatal(err)
	}
}
