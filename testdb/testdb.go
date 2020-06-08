package testdb

import (
	"github.com/jmoiron/sqlx"

	"github.com/imclaren/sqldb"
	"github.com/imclaren/sqldb/postgres"

	"fmt"

	"context"
)

func DB(ctx context.Context, cancelFunc context.CancelFunc, dbType, dbName string) (sqldb.DB, error) {
	var db *sqlx.DB

	// Create DB if required
	err := Create(ctx, cancelFunc, dbType, dbName)
	if err != nil {
		return sqldb.DB{}, err
	}

	switch dbType {
	case "postgres":
		db, err = postgres.TestDB(ctx, dbName)
		if err != nil {
			return sqldb.DB{}, err
		}
	default:
		return sqldb.DB{}, fmt.Errorf("Error: test database type not implemented: %s", dbType)
	}

	//defer db.Close()
	return sqldb.WithDB(ctx, cancelFunc, dbType, db), nil
}

func Create(ctx context.Context, cancelFunc context.CancelFunc, dbType, dbName string) error {
	switch dbType {
	case "postgres":
		sqlDb, err := postgres.TestDB(ctx, "")
		if err != nil {
			return err
		}
		db := sqldb.WithDB(ctx, cancelFunc, dbType, sqlDb)
		err = db.CreateDB(dbName)
		if err != nil {
			return err
		}
		defer db.Close()
	default:
		return fmt.Errorf("Error: test database type not implemented: %s", dbType)
	}
	return nil
}

func Drop(ctx context.Context, cancelFunc context.CancelFunc, dbType, dbName string) error {
	switch dbType {
	case "postgres":
		sqlDb, err := postgres.TestDB(ctx, "")
		if err != nil {
			return err
		}
		db := sqldb.WithDB(ctx, cancelFunc, dbType, sqlDb)
		err = db.DropDB(dbName)
		if err != nil {
			return err
		}
		defer db.Close()
	default:
		return fmt.Errorf("Error: test database type not implemented: %s", dbType)
	}
	return nil
}
