package system

import "errors"

var (
	ErrOpeningCsv             = errors.New("error opening csv")
	ErrReadingCsv             = errors.New("error reading csv")
	ErrCantGetCsvFile         = errors.New("can't get csv file")
	ErrCantGetTransactionInfo = errors.New("can't get transaction info")
	ErrReadTemplateFile       = errors.New("can't read template file")
	ErrTemplateParse          = errors.New("can't parse template")
	ErrCreateOutputFile       = errors.New("can't create output file")
	ErrTemplateExecute        = errors.New("can't execute template")
	ErrReadFile               = errors.New("can't read file")
)

const (
	CantGetInfo   string = "can't get info"
	CantWriteHtml string = "can't write html"
)
