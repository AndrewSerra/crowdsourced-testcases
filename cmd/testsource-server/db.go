/*
 * Created on Mon Dec 16 2024
 *
 * Copyright © 2024 Andrew Serra <andy@serra.us>
 */
package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

const (
	FOREIGN_KEY_NO_EXIST_ERROR uint16 = 1452
	DUPLICATE_ENTRY_ERROR      uint16 = 1062
)

func init() {
	cfg := mysql.Config{
		User:      "root",
		Passwd:    "admin",
		Net:       "tcp",
		Addr:      "127.0.0.1:3310",
		DBName:    "csdb",
		ParseTime: true,
	}

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}
	log.Println("Connected to database.")
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
