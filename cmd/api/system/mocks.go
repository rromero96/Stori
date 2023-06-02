package system

import (
	"context"
	"time"
)

// MockReadCSV mock
func MockReadCSV(trans []Transaction, err error) ReadCSV {
	return func(context.Context, string) ([]Transaction, error) {
		return trans, err
	}
}

// MockProcessTransactions mock
func MockProcessTransactions(email Email, err error) ProcessTransactions {
	return func(context.Context) (Email, error) {
		return email, err
	}
}

// MockTransaction mock
func MockTransaction(id int64, date time.Time, trType string, amount float64) Transaction {
	return Transaction{
		ID:          id,
		Date:        date,
		Transaction: amount,
		Type:        trType,
	}
}

// MockTranssactions mock
func MockTransactions() []Transaction {
	return []Transaction{
		MockTransaction(0, time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), "credit", +60.5),
		MockTransaction(1, time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC), "debit", -10.3),
		MockTransaction(2, time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC), "debit", -20.46),
		MockTransaction(3, time.Date(2023, 1, 4, 0, 0, 0, 0, time.UTC), "credit", +10),
		MockTransaction(4, time.Date(2023, 1, 5, 0, 0, 0, 0, time.UTC), "credit", +61.5),
		MockTransaction(5, time.Date(2023, 1, 6, 0, 0, 0, 0, time.UTC), "debit", -11.4),
		MockTransaction(6, time.Date(2023, 1, 7, 0, 0, 0, 0, time.UTC), "debit", -21.46),
		MockTransaction(7, time.Date(2023, 1, 8, 0, 0, 0, 0, time.UTC), "credit", +11),
		MockTransaction(8, time.Date(2023, 1, 9, 0, 0, 0, 0, time.UTC), "credit", +62.5),
		MockTransaction(9, time.Date(2023, 1, 10, 0, 0, 0, 0, time.UTC), "debit", -12.4),
		MockTransaction(10, time.Date(2023, 1, 11, 0, 0, 0, 0, time.UTC), "debit", -22.46),
		MockTransaction(11, time.Date(2023, 1, 12, 0, 0, 0, 0, time.UTC), "credit", +12),
		MockTransaction(12, time.Date(2023, 1, 13, 0, 0, 0, 0, time.UTC), "credit", +63.5),
		MockTransaction(13, time.Date(2023, 1, 14, 0, 0, 0, 0, time.UTC), "debit", -13.4),
		MockTransaction(14, time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC), "debit", -23.46),
		MockTransaction(15, time.Date(2023, 1, 16, 0, 0, 0, 0, time.UTC), "credit", +13),
		MockTransaction(16, time.Date(2023, 2, 17, 0, 0, 0, 0, time.UTC), "credit", +64.5),
		MockTransaction(17, time.Date(2023, 2, 18, 0, 0, 0, 0, time.UTC), "debit", -14.5),
		MockTransaction(18, time.Date(2023, 2, 19, 0, 0, 0, 0, time.UTC), "debit", -23.46),
		MockTransaction(19, time.Date(2023, 2, 20, 0, 0, 0, 0, time.UTC), "credit", +14),
		MockTransaction(20, time.Date(2023, 2, 21, 0, 0, 0, 0, time.UTC), "credit", +65.5),
	}
}

// MockEmail mock
func MockEmail() Email {
	monthCount := MockMonthsMap()
	email := Email{
		Balance:       264.69999999999993,
		AverageDebit:  -86.65000000000002,
		AverageCredit: 219,
	}

	email.Body = getEmailBody(email.Balance, email.AverageDebit, email.AverageCredit, monthCount)

	return email
}

// MockMonthsMap mock
func MockMonthsMap() map[time.Month]int {
	return map[time.Month]int{
		time.January:  16,
		time.February: 5,
	}
}
