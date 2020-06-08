package sqldb

import (
	"database/sql"
	//_ "github.com/mattn/go-sqlite3"
	"github.com/jmoiron/sqlx"

	"github.com/imclaren/sqldb/postgres"
	"github.com/imclaren/sqldb/sqlite"

	"context"
	"fmt"
	"strings"
	"sync"
)

type DB struct {
	Ctx        context.Context
	CancelFunc context.CancelFunc
	*sqlx.DB
	Mutex sync.RWMutex
	Type  string
}

func WithDB(ctx context.Context, cancelFunc context.CancelFunc, dbType string, db *sqlx.DB) DB {
	return DB{
		ctx,
		cancelFunc,
		db,
		sync.RWMutex{},
		dbType,
	}
}

func Init(ctx context.Context, cancelFunc context.CancelFunc, dbType, connectString string) (DB, error) {

	var db *sqlx.DB
	var err error

	switch dbType {
	case "sqlite":
		db, err = sqlite.Init(ctx, connectString)
		if err != nil {
			return DB{}, err
		}
	case "postgres":
		db, err = postgres.Init(ctx, connectString)
		if err != nil {
			return DB{}, err
		}
	default:
		return DB{}, fmt.Errorf("Error: database type not implemented: %s", dbType)
	}

	//defer db.Stop()
	return DB{
		ctx,
		cancelFunc,
		db,
		sync.RWMutex{},
		dbType,
	}, nil
}

func (db *DB) Stop() error {
	defer db.CancelFunc()
	return db.Close()
}

func (db *DB) CreateDB(dbName string) error {
	db.Mutex.Lock()
	defer db.Mutex.Unlock()

	if dbName == "" {
		return fmt.Errorf("no database name provided")
	}

	sqlString := fmt.Sprintf(`CREATE DATABASE %s`, dbName)
	_, err := db.Exec(sqlString)
	if err != nil {
		if strings.HasSuffix(err.Error(), fmt.Sprintf("database \"%s\" already exists", dbName)) {
			return nil
		}
		return err
	}
	return nil
}

func (db *DB) DropDB(dbName string) error {
	db.Mutex.Lock()
	defer db.Mutex.Unlock()

	sqlString := fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName)
	_, err := db.Exec(sqlString)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) WriteSlice(sqlStrings []string) error {

	// Block if already writing to database
	db.Mutex.Lock()
	defer db.Mutex.Unlock()

	for _, sql := range sqlStrings {
		_, err := db.Exec(sql)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *DB) RowExists(query string, args ...interface{}) (bool, error) {
	var exists bool
	query = fmt.Sprintf("SELECT exists (%s)", query)
	err := db.QueryRow(query, args...).Scan(&exists)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return false, nil
		default:
			// continue
		}
		return false, err
	}
	return exists, nil
}
