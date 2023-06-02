package system

import (
	"context"
	"path/filepath"
	"runtime"
)

const (
	folder string = "data"
	file   string = "data.csv"
)

type ProcessTransactions func(ctx context.Context) (Email, error)

func MakeProcessTransactions(readCSV ReadCSV) ProcessTransactions {
	return func(ctx context.Context) (Email, error) {
		var email Email

		transactions, err := readCSV(ctx, GetFileName())
		if err != nil {
			return Email{}, ErrCantGetCsvFile
		}

		email.Balance, email.AverageDebit, email.AverageCredit = getBalanceInfo(transactions)
		monthCount := transactionsPerMonth(transactions)

		email.Body = getEmailBody(email.Balance, email.AverageDebit, email.AverageCredit, monthCount)

		return email, nil
	}
}

func GetFileName() string {
	// Get the current file's directory
	_, filename, _, _ := runtime.Caller(0)
	testDir := filepath.Dir(filename)
	// Construct the absolute file path to the CSV file
	return filepath.Join(testDir, folder, file)
}
