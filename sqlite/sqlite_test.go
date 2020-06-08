package sqlite

import (
	assert "github.com/stretchr/testify/require"
	"testing"

	"os"
	"path/filepath"

	"context"
	"time"
)

func TestInit(t *testing.T) {
	sqlPath := filepath.Join("testdata", "test.db")
	defer os.Remove(sqlPath)
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()
	db, err := Init(ctx, sqlPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	assert.NotNil(t, db)
}
