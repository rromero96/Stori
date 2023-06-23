package system

import "errors"

var (
	ErrOpeningCsv             = errors.New("error opening csv")
	ErrReadingCsv             = errors.New("error reading csv")
	ErrCantGetCsvFile         = errors.New("can't get csv file")
	ErrCantGetTransactionInfo = errors.New("can't get transaction info")
	ErrReadTemplateFile       = errors.New("can't read template file")
	ErrTemplateParse          = errors.New("can't parse template")
	ErrTemplateExecute        = errors.New("can't execute template")
	ErrCantPrepareStatement   = errors.New("can't prepare statement")
	ErrCantRunQuery           = errors.New("can't run query")
	ErrCantGetLastID          = errors.New("can't get last id")
	ErrCantCreateTransactions = errors.New("can't create transactions")
)

const (
	CantGetInfo         string = "can't get info"
	CantWriteHtml       string = "can't write html"
	CantWriteSwaggerYML string = "can't write swagger yml"
)
