package storage

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
)

func (m *MetricStorage) saveToDB(dbConnString string) {
	db, err := sql.Open("pgx",
		dbConnString)
	if err != nil {
		log.Println(err.Error())
	}
	defer func() {
		if err = db.Close(); err != nil {
			log.Print(err.Error())
		}
	}()

	writeMetric(db, m)
}

func writeMetric(db *sql.DB, m *MetricStorage) {
	query := `INSERT INTO metrics(type, name, value)
		VALUES($1, $2, $3)
		ON CONFLICT (type, name) DO UPDATE
		SET value = $3;`
	var err error
	tx, err := db.Begin()
	if err != nil {
		log.Println(err.Error())
		return
	}

	defer func() {
		if err != nil {
			log.Println(err.Error())
			if err = tx.Rollback(); err != nil {
				log.Println(err.Error())
			}
		} else {
			if err = tx.Commit(); err != nil {
				log.Println(err.Error())
			}
		}
	}()
	ctx := context.Background()
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer func() {
		if err = stmt.Close(); err != nil {
			log.Println(err.Error())
		}
	}()

	for k, v := range m.Gauge {
		if _, err = stmt.ExecContext(ctx, "gauge", k, v); err != nil {
			log.Println(err.Error())
		}
	}
	for k, v := range m.Counter {
		if _, err := stmt.ExecContext(ctx, "counter", k, v); err != nil {
			log.Println(err.Error())
		}
	}

}
