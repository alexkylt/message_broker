package dbstorage

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/lib/pq"
)

const (
	DB_USER     = "devuser"
	DB_PASSWORD = "devuser"
	DB_NAME     = "storage_db"
	PORT        = 5432
	HOST        = "postgres"
)

type dbStorage struct {
	db *sql.DB
}

func InitStorage() *dbStorage {

	connection := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		HOST, PORT, DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", connection)

	if err != nil {
		log.Fatal(err)
	}

	return &dbStorage{db: db}
}

func (d *dbStorage) Get(key string) (string, error) {

	var key_value string
	select_stmt := `SELECT key_value FROM kv_storage where key_name = $1;`

	row := d.db.QueryRow(select_stmt, key)
	err := row.Scan(&key_value)

	fmt.Println(row, err, select_stmt, key, key_value)

	if err != nil {
		if err == sql.ErrNoRows {
			value := fmt.Sprintf("There is no %s key in th db.", key)
			return value, err
		} else {
			errHandler(err)
		}
	}

	return key_value, nil
}

func (d *dbStorage) Set(key, value string) error {

	var lastInsertId int
	datetime := time.Now().Format(time.RFC3339)
	insert_stmt := `INSERT INTO kv_storage(key_name, key_value, datetime) VALUES($1,$2,$3) returning id;`
	fmt.Println("VALUES to insert:", key, value)
	err := d.db.QueryRow(insert_stmt, key, value, datetime).Scan(&lastInsertId)

	if err, ok := err.(*pq.Error); ok {
		fmt.Println("pq error:", err.Code.Name(), err.Code)
		if err.Code == "23505" {
			update_stmt := `UPDATE kv_storage set key_value=$1, datetime=$2 where key_name=$3 returning id;`
			err := d.db.QueryRow(update_stmt, value, datetime, key).Scan(&lastInsertId)
			if err != nil {
				errHandler(err)
			}
		} else {
			errHandler(err)
		}
	}
	fmt.Println("last inserted id =", lastInsertId, err)

	return nil
}

func (d *dbStorage) Delete(key string) error {

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

func (d *dbStorage) Keys(pattern string) ([]string, error) {

	select_stmt := `SELECT key_value FROM kv_storage where key_name like '%$1%';`

	rows, err := d.db.Query(select_stmt, pattern)
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
