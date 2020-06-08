# sqldb
sqldb allows access to databases in golang.  It is a thin wrapper around github.com/jmoiron/sqlx.

It is designed to allow the reuse of the same code for multiple types of databases (see the insert and all functions below), and for database specific sql where required (see the createTable function below).

```
// Item is a cache item
type Item struct {
	Id              int
	Bucket 			string 		`db:"bucket"`
	Key 			string 		`db:"key"`
	Size 			int64 		`db:"size"`
	AccessCount 	int64  		`db:"access_count"`
	ExpiresAt 		time.Time  	`db:"expires_at"`
	CreatedAt       time.Time 	`db:"created_at"`
	UpdatedAt       time.Time 	`db:"updated_at"`
}

func main() {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()
	dbList := []struct {
		dbType  string
		connectString  string
	}{
		{"sqlite", "/path/to/sqldatabase.db"},
		{"postgresql", postgres.ConnectString("localhost", 5432, "postgres", "", "", false, "UTC")}
	}
	for _, d := range dbList {
		db, err := sqldb.Init(ctx, cancelFunc, d.dbType, d.connectString)
		if err != nil {
			log.Fatal(err)
		}
		err = createTable(&db)
		if err != nil {
			log.Fatal(err)
		}
		err = insert(&db, "mybucket", "mykey", 100, 1, time.Now())
		if err != nil {
			log.Fatal(err)
		}
		items, err = all(db *DB)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(items)
		err = db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func insert(db *sqldb.DB, bucket, key string, size, access_count int64, expires_at time.Time) error {
	db.Mutex.Lock()
	defer db.Mutex.Unlock()

	sqlString := "INSERT INTO cache (bucket, key, size, access_count, expires_at) VALUES (?,?,?,?,?)"
	_, err := db.Exec(db.Rebind(sqlString),
		i.Bucket,
		i.Key,
		i.Size,
		i.AccessCount,
		i.ExpiresAt,
	)
	return err
}

func all(db *sqldb.DB) ([]Item, error) {
	db.Mutex.RLock()
	defer db.Mutex.RUnlock()

	sqlString := "SELECT * FROM cache"
	var items []Item
	err := db.Select(&items, db.Rebind(sqlString))
	return items, err
}

func createTable(db *sqldb.DB) error {
	db.Mutex.Lock()
	defer db.Mutex.Unlock()

	switch db.Type {
	case "sqlite":
		_, err := db.Exec(`
    		CREATE TABLE IF NOT EXISTS cache (
	    		id INTEGER PRIMARY KEY,
	    		bucket TEXT,
	    		key TEXT,
	    		size INT,
	    		access_count INT,
	    		expires_at TIMESTAMP,
				created_at TIMESTAMP NULL DEFAULT(STRFTIME('%Y-%m-%d %H:%M:%f', 'NOW')),
			    updated_at TIMESTAMP NULL DEFAULT(STRFTIME('%Y-%m-%d %H:%M:%f', 'NOW'))	    		
    			)
    		`)
		if err != nil {
			return err
		}
		// Create updated_at trigger for cache table
		_, err = db.Exec(`
			CREATE TRIGGER [update_cache_updated_at]
			    AFTER UPDATE
			    ON cache
			BEGIN
			    UPDATE cache SET updated_at=STRFTIME('%Y-%m-%d %H:%M:%f', 'NOW') WHERE id=NEW.id;
			END;
		`)
		if err != nil && err.Error() != "trigger [update_cache_updated_at] already exists" {
			return err
		}
	case "postgres":
		_, err := db.Exec(`
			CREATE TABLE IF NOT EXISTS cache (
				id BIGSERIAL PRIMARY KEY, 
				bucket TEXT,
				key TEXT, 
				size BIGINT, 
				access_count BIGINT, 
				expires_at TIMESTAMP, 
				created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
			    updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP
			)
		`)
		if err != nil {
			return err
		}
		// Create or replace update_updated_at_column function
		// Note we only need to do this once for all of the tables that we update
		_, err = db.Exec(`
			CREATE OR REPLACE FUNCTION update_updated_at_column()
			RETURNS TRIGGER AS $$
			BEGIN
			   NEW.updated_at = now(); 
			   RETURN NEW;
			END;
			$$ language 'plpgsql';
		`)
		if err != nil {
			return err
		}
		// Create updated_at trigger for cache table
		_, err = db.Exec(`
			CREATE TRIGGER update_cache_updated_at 
				BEFORE UPDATE
				ON cache 
				FOR EACH ROW 
				EXECUTE PROCEDURE update_updated_at_column();
		`)
		if err != nil && err.Error() != "pq: trigger \"update_cacheupdated_at\" for relation \"cache\" already exists" {
			return err
		}
	default:
		return fmt.Errorf("Create table error: database type not implemented: %s", db.Type)
	}
	return nil
}
```