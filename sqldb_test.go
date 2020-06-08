package sqldb

import (
	assert "github.com/stretchr/testify/require"
	"testing"

	"context"
	"os"
	"path/filepath"
	"time"
)

func TestInit(t *testing.T) {
	sqlPath := filepath.Join("testdata", "test.db")
	defer os.Remove(sqlPath)
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()
	db, err := Init(ctx, cancelFunc, "sqlite", sqlPath)
	if err != nil {
		t.Fatal(err)
	}
	err = db.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestWriteSlice(t *testing.T) {
	sqlPath := filepath.Join("testdata", "test.db")
	defer os.Remove(sqlPath)
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()
	db, err := Init(ctx, cancelFunc, "sqlite", sqlPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	err = db.WriteSlice([]string{"CREATE TABLE IF NOT EXISTS table1 (id INTEGER PRIMARY KEY, filename TEXT)"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestRowExists(t *testing.T) {
	sqlPath := filepath.Join("testdata", "test.db")
	defer os.Remove(sqlPath)
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()
	db, err := Init(ctx, cancelFunc, "sqlite", sqlPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS table1 (id INTEGER PRIMARY KEY, filename TEXT)")
	if err != nil {
		t.Fatal(err)
	}

	exists, err := db.RowExists("SELECT id FROM table1 WHERE filename=?", "test.jpg")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, false, exists)

	_, err = db.Exec("INSERT INTO table1 (filename) VALUES (?)", "test.jpg")
	if err != nil {
		t.Fatal(err)
	}
	exists, err = db.RowExists("SELECT id FROM table1 WHERE filename=?", "test.jpg")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, true, exists)
}
