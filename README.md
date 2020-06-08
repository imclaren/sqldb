# sqldb
sqldb allows access to databases in golang.  It is a thin wrapper around github.com/jmoiron/sqlx.

```
ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
defer cancelFunc()
db, err := sqldb.Init(ctx, cancelFunc, "sqlite", sqlPath)
if err != nil {
	log.Fatal(err)
}
err = db.Close()
if err != nil {
	log.Fatal(err)
}
```