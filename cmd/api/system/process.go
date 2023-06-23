package system

import (
	"context"
	"html/template"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	path         string = "api/system/data"
	file         string = "data.csv"
	HtmlFolder   string = "api/system/html"
	htmlFile     string = "account_info.html"
	StoriLogo    string = "stori_logo.jpeg"
	templateFile string = "template.html"
)

type (
	// HTMLProcessTransactions renders an HTML from the data recieved in the CSV file
	HTMLProcessTransactions func(ctx context.Context) ([]byte, error)
)

// MakeHTMLProcessTransactions creates an HTMLProcessTransactions function
func MakeHTMLProcessTransactions(readCSV ReadCSV, mySQLCreate MySQLCreate) HTMLProcessTransactions {
	return func(ctx context.Context) ([]byte, error) {
		var email Email

		transactions, err := readCSV(ctx, GetFileName(path, file))
		if err != nil {
			return []byte{}, ErrCantGetCsvFile
		}

		err = mySQLCreate(ctx, transactions)
		if err != nil {
			return []byte{}, ErrCantCreateTransactions
		}

		email.Balance, email.AverageDebit, email.AverageCredit = getBalanceInfo(transactions)
		email.WorkingMonths = transactionsPerMonth(transactions)

		templateFile := GetFileName(HtmlFolder, templateFile)
		tmplBytes, err := os.ReadFile(templateFile)
		if err != nil {
			return []byte{}, ErrReadTemplateFile
		}

		var buf strings.Builder
		templateName := "accountInfo"
		tmpl, err := template.New(templateName).Parse(string(tmplBytes))
		if err != nil {
			return []byte{}, ErrTemplateParse
		}

		err = tmpl.Execute(&buf, email)
		if err != nil {
			return []byte{}, ErrTemplateExecute
		}

		htmlBytes := []byte(buf.String())
		return htmlBytes, nil
	}
}

// GetFileName returns the absolute file path of a file
func GetFileName(folder string, file string) string {
	_, filename, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(filename)
	rootDir := filepath.Join(currentDir, "..", "..")

	return filepath.Join(rootDir, folder, file)
}
