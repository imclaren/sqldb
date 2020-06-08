package testdb

import (
	"testing"
	"context"
	"time"
)

func TestTestDB(t *testing.T) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	db, err := DB(ctx, cancelFunc, "postgres", "testdbtest1")
	if err != nil {
		t.Fatal(err)
	}
	err = db.Close()
	if err != nil {
		t.Fatal(err)
	}
	err = Drop(ctx, cancelFunc, "postgres", "testdbtest1")
	if err != nil {
		t.Fatal(err)
	}
}
