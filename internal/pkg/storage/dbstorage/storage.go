package dbstorage

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/lib/pq"
)

const (
	dbUser     = "devuser"
	dbPassword = "devuser"
	dbName     = "storage_db"
	port        = 5432
	host        = "postgres"
)
// DbStorage ...
type DbStorage struct {
	db *sql.DB
}
// InitStorage ...
func InitStorage() *DbStorage {

	connection := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", connection)

	if err != nil {
		log.Fatal(err)
	}

	return &DbStorage{db: db}
}
// Get ...
func (d *DbStorage) Get(key string) (string, error) {

	var keyValue string
	selectStmt := `SELECT key_value FROM kv_storage where key_name = $1;`

	row := d.db.QueryRow(selectStmt, key)
	err := row.Scan(&keyValue)

	fmt.Println(row, err, selectStmt, key, keyValue)

	if err != nil {
		if err == sql.ErrNoRows {
			keyValue = fmt.Sprintf("There is no %s key in th db.", key)
			// return value, err
		} else {
			errHandler(err)
		}
	}

	return keyValue, nil
}
// Set ...
func (d *DbStorage) Set(key, value string) error {

	var lastInsertID int
	datetime := time.Now().Format(time.RFC3339)
	insertStmt := `INSERT INTO kv_storage(key_name, key_value, datetime) VALUES($1,$2,$3) returning id;`
	fmt.Println("VALUES to insert:", key, value)
	err := d.db.QueryRow(insertStmt, key, value, datetime).Scan(&lastInsertID)

	if err, ok := err.(*pq.Error); ok {
		fmt.Println("pq error:", err.Code.Name(), err.Code)
		if err.Code == "23505" {
			updateStmt := `UPDATE kv_storage set key_value=$1, datetime=$2 where key_name=$3 returning id;`
			err := d.db.QueryRow(updateStmt, value, datetime, key).Scan(&lastInsertID)
			if err != nil {
				errHandler(err)
			}
		} else {
			errHandler(err)
		}
	}
	fmt.Println("last inserted id =", lastInsertID, err)

	return nil
}
// Delete ...
func (d *DbStorage) Delete(key string) error {

	sqlStatement := `delete from kv_storage where key_name=$1;`
	res, err := d.db.Exec(sqlStatement, key)
	count, err := res.RowsAffected()
	if err != nil {
		errHandler(err)
	}
	fmt.Println("rows changed", err, "count - ", count)
	return nil
}

func errHandler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
// Keys ...
func (d *DbStorage) Keys(pattern string) ([]string, error) {

	selectStmt := `SELECT key_value FROM kv_storage where key_name like '%$1%';`

	rows, err := d.db.Query(selectStmt, pattern)
	//err := row.Scan(&key_value)

	values := make([]string, 0)
	for rows.Next() {
		var value string
		if err := rows.Scan(&value); err != nil {
			log.Fatal(err)
		}
		values = append(values, value)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(values, err)

	return values, nil
}
