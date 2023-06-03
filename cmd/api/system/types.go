package system

import (
	"time"
)

type (
	Transaction struct {
		ID          int64
		Date        time.Time
		Transaction float64
		Type        string
	}

	Email struct {
		Balance       float64
		AverageDebit  float64
		AverageCredit float64
		WorkingMonths map[string]int
	}
)

func (e Email) toDTO() EmailDTO {
	return EmailDTO(e)
}

func getBalanceInfo(transactions []Transaction) (float64, float64, float64) {
	var total, debit, credit float64
	for _, t := range transactions {
		total += t.Transaction

		if t.Type == "debit" {
			debit += t.Transaction
		}
		if t.Type == "credit" {
			credit += t.Transaction
		}
	}

	return total, debit / 2, credit / 2
}

func transactionsPerMonth(transactions []Transaction) map[string]int {
	monthCount := make(map[string]int)

	for _, t := range transactions {
		month := t.Date.Month().String()
		monthCount[month]++
	}

	return monthCount
}
