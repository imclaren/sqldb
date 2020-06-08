package postgres

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"context"

	"fmt"
)

type Params struct {
	Host     string
	Port     int
	User     string
	Password string
	DbName   string
	Ssl      bool
	Timezone string
}

func TestParams() Params {
	return Params{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "",
		//DbName: "",
		Ssl: false,
		Timezone: "UTC",
	}
}

func TestDB(ctx context.Context, dbName string) (db *sqlx.DB, err error) {
	p := TestParams()
	p.DbName = dbName
	connectString := ConnectString(p.Host, p.Port, p.User, p.Password, p.DbName, p.Ssl, p.Timezone)
	return Init(ctx, connectString)
}

func ConnectString(host string, port int, user, password, dbName string, ssl bool, timezone string) string {
	sslString := "require"
	if !ssl {
		sslString = "disable"
	}
	dbNameString := fmt.Sprintf("dbname='%s'", dbName)
	if dbName == "" {
		dbNameString = ""
	}
	timezoneString := fmt.Sprintf("TimeZone='%s'", timezone)
	if timezone == "" {
		timezoneString = ""
	}
	return fmt.Sprintf("host='%s' port='%d' user='%s' password='%s' %s sslmode='%s' %s", host, port, user, password, dbNameString, sslString, timezoneString)
}

func Init(ctx context.Context, connectString string) (db *sqlx.DB, err error) {
	db, err = sqlx.ConnectContext(ctx, "postgres", connectString)
	return db, err
}
