package system

import (
	"fmt"
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
		Body          string
	}
)

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

func transactionsPerMonth(transactions []Transaction) map[time.Month]int {
	monthCount := make(map[time.Month]int)

	for _, t := range transactions {
		month := t.Date.Month()
		monthCount[month]++
	}

	return monthCount
}

func getEmailBody(balance float64, avgDebit float64, avgCredit float64, workingMonths map[time.Month]int) string {
	return fmt.Sprintf(`Hello,

	Here is your accounts information.

	Total Balance is: $ %v ,
	Average Debit amount is: $ %v ,
	Average Credit amount is: $ %v ,
	%s

	Thanks,
	Your Bank
	`, balance, avgDebit, avgCredit, getTransactionsPerMonth(workingMonths))
}

func getTransactionsPerMonth(workingMonths map[time.Month]int) string {
	var res string

	for month, count := range workingMonths {
		res += fmt.Sprintf("Number of transactions in %s: %d\n", month.String(), count)
	}

	return res
}
