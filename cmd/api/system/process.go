package system

import (
	"context"
	"encoding/csv"
	"os"
	"strconv"
	"time"
)

type (
	ReadCSV func(ctx context.Context, filename string) ([]Transaction, error)
)

func MakeReadCSV() ReadCSV {
	return func(ctx context.Context, filename string) ([]Transaction, error) {
		file, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		reader := csv.NewReader(file)
		records, err := reader.ReadAll()
		if err != nil {
			return nil, err
		}

		var transactions []Transaction
		for i, record := range records {
			if i == 0 {
				continue
			}

			id, _ := strconv.ParseInt(record[0], 10, 64)
			date, _ := time.Parse("2/1", record[1])
			currentYear := time.Now().Year()
			date = time.Date(currentYear, date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
			amount, _ := strconv.ParseFloat(record[2], 64)

			transaction := Transaction{
				ID:          int64(id),
				Date:        date,
				Transaction: float64(amount),
			}

			transaction.Type = "credit"
			if amount < 0 {
				transaction.Type = "debit"
			}

			transactions = append(transactions, transaction)
		}

		return transactions, nil
	}
}
