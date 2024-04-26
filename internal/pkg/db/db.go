package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func InitDB(DBUrl string) (*sql.DB, error) {
	db, err := sql.Open("postgres", DBUrl)

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	fmt.Println("Database Conneted")
	return db, nil
}

func SetDeductionValue(db *sql.DB, key string, value float64) float64 {
	db.QueryRow("INSERT INTO \"deductions\" (\"name\", maxAmount) VALUES ($1, $2) ON CONFLICT (\"name\") DO UPDATE SET maxAmount = EXCLUDED.maxAmount RETURNING maxAmount;", key, value).Scan(&value)
	return value
}

func GetDeductionValue(db *sql.DB, key string) float64 {
	var value float64
	db.QueryRow("SELECT 'maxAmount' FROM \"dedictions\" WHERE \"name\" = $1;", key).Scan(&value)
	return value
}
