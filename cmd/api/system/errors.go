package system

import "errors"

var (
	ErrOpeningCsv     = errors.New("error opening csv")
	ErrReadingCsv     = errors.New("error reading csv")
	ErrCantGetCsvFile = errors.New("can't get csv file")
)
