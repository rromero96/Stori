package system

import (
	"context"
	"database/sql"
	"strings"
)

const queryCreate = "INSERT INTO stori.transactions (id, date, transaction, type) VALUES "

var lastID *int64

// MySQLCreate is a function that creates a transaction in the database
type MySQLCreate func(ctx context.Context, transactions []Transaction) error

// MakeMySQLCreate creates a new MySQLCreate
func MakeMySQLCreate(db *sql.DB) MySQLCreate {
	return func(ctx context.Context, transactions []Transaction) error {
		if lastID == nil {
			initialValue := int64(-1)
			lastID = &initialValue
		}

		if *lastID < transactions[0].ID {
			var inserts []string
			var params []interface{}

			for _, t := range transactions {
				inserts = append(inserts, "(?, ?, ?, ?)")
				params = append(params, t.ID, t.Date, t.Transaction, t.Type)
			}

			queryVals := strings.Join(inserts, ",")
			query := queryCreate + queryVals

			stmt, err := db.PrepareContext(ctx, query)
			if err != nil {
				return ErrCantPrepareStatement
			}
			defer stmt.Close()

			r, err := stmt.ExecContext(ctx, params...)
			if err != nil {
				return ErrCantRunQuery
			}

			id, err := r.LastInsertId()
			if err != nil {
				return ErrCantGetLastID
			}
			lastID = &id
		}

		return nil
	}
}
