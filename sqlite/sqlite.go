package sqlite

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"context"
)

func ConnectString(path string, timezone string) string {
	timezoneString := fmt.Sprintf("?_loc=%s", timezone)
	if timezone == "" {
		timezoneString = ""
	}
	return fmt.Sprintf("file:%s%s", path, timezoneString)
}

func Init(ctx context.Context, connectString string) (db *sqlx.DB, err error) {
	return sqlx.ConnectContext(ctx, "sqlite3", connectString)
}
