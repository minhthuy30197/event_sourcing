package model

import (
	"log"
	"time"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)


func ConnectDb(user string, password string, database string, address string) (db *pg.DB) {
	db = pg.Connect(&pg.Options{
		User:     user,
		Password: password,
		Database: database,
		Addr:     address,
	})

	return db
}

func MigrationDb(db *pg.DB, schema string) error {
	// Tạo schema theo tên service
	_, err := db.Exec("CREATE SCHEMA IF NOT EXISTS " + schema + ";")
	if err != nil {
		return err
	}

	// Tạo bảng
	var course Course
	var class Class
	var teacher Teacher
	err = createTable(&course, db)
	if err != nil {
		return err
	}
	err = createTable(&class, db)
	if err != nil {
		return err
	}
	err = createTable(&teacher, db)
	if err != nil {
		return err
	}

	return nil
}

func MigrationEventDb(db *pg.DB, schema string) error {
	// Thêm extension timescaledb
	_, err := db.Exec("CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;")
	if err != nil {
		return err
	}

	// Tạo schema theo tên service
	_, err = db.Exec("CREATE SCHEMA IF NOT EXISTS " + schema + ";")
	if err != nil {
		return err
	}

	// Tạo table
	var event EventSource
	err = createTable(&event, db)
	if err != nil {
		return err
	}

	// Tạo hypertable
	_, err = db.Exec("SELECT create_hypertable('es.event_source', 'time', if_not_exists => TRUE);")
	if err != nil {
		return err
	}

	return nil
}

func LogQueryToConsole(db *pg.DB) {
	db.OnQueryProcessed(func(event *pg.QueryProcessedEvent) {
		query, err := event.FormattedQuery()
		if err != nil {
			panic(err)
		}

		log.Printf("%s %s", time.Since(event.StartTime), query)
	})
}

func createTable(model interface{}, db *pg.DB) error {
	err := db.CreateTable(model, &orm.CreateTableOptions{
		Temp:          false,
		FKConstraints: true,
		IfNotExists:   true,
	})

	return err
}
