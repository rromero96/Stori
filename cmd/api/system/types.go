package system

import "time"

type Transaction struct {
	ID          int64
	Date        time.Time
	Transaction float64
	Type        string
}
