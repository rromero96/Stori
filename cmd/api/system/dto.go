package system

type EmailDTO struct {
	Balance       float64        `json:"balance"`
	AverageDebit  float64        `json:"average_debit"`
	AverageCredit float64        `json:"average_credit"`
	WorkingMonths map[string]int `json:"working_months"`
}
