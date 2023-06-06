package system

import (
	"context"
	"database/sql"
	"strings"
)

const (
	queryCreate = "INSERT INTO stori.transactions (id, date, transaction, type) VALUES "
	queryFind   = "SELECT MAX(id) FROM stori.transactions"
)

type (
	// MySQLCreate is a function that creates a transaction in the database
	MySQLCreate func(ctx context.Context, transactions []Transaction) error

	// MySQLFind is a function that finds the last id in the database
	MySQLFind func(ctx context.Context) (int64, error)
)

// MakeMySQLCreate creates a new MySQLCreate
func MakeMySQLCreate(db *sql.DB, mySQLFind MySQLFind) MySQLCreate {
	return func(ctx context.Context, transactions []Transaction) error {
		lastID, err := mySQLFind(ctx)
		if err != nil {
			return ErrCantGetLastID
		}

		if lastID < transactions[0].ID {
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

			_, err = stmt.ExecContext(ctx, params...)
			if err != nil {
				return ErrCantRunQuery
			}
		}

		return nil
	}
}

// MakeMySQLFind creates a new MySQLFind
func MakeMySQLFind(db *sql.DB) MySQLFind {
	return func(ctx context.Context) (int64, error) {
		var lastID sql.NullInt64
		err := db.QueryRow(queryFind).Scan(&lastID)
		if err != nil {
			if err == sql.ErrNoRows {
				return -1, nil
			}

			return 0, ErrCantRunQuery
		}
		if lastID.Valid {
			return lastID.Int64, nil
		}

		return -1, nil
	}
}
