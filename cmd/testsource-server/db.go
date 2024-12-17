/*
 * Created on Mon Dec 16 2024
 *
 * Copyright © 2024 Andrew Serra <andy@serra.us>
 */
package main

import (
	"database/sql"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3310",
		DBName: "csdb",
	}

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		panic(err)
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
}

func GetDB() *sql.DB {
	return db
}

func SafelyCloseDB() {
	db.Close()
}
