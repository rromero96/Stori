package system

import (
	"context"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

const (
	folder       string = "data"
	file         string = "data.csv"
	htmlFolder   string = "html"
	htmlFile     string = "account_info.html"
	storiLogo    string = "stori_logo.jpeg"
	templateFile string = "template.html"
)

type (
	// ProcessTransactions transforms the data recieved in the CSV file into an Email struct
	ProcessTransactions func(ctx context.Context) (Email, error)

	// HTMLProcessTransactions renders an HTML from the data recieved in the CSV file
	HTMLProcessTransactions func(ctx context.Context) ([]byte, error)
)

// MakeHTMLProcessTransactions creates an HTMLProcessTransactions function
func MakeHTMLProcessTransactions(processTransactions ProcessTransactions) HTMLProcessTransactions {
	return func(ctx context.Context) ([]byte, error) {
		email, err := processTransactions(ctx)
		if err != nil {
			return []byte{}, ErrCantGetTransactionInfo
		}

		templateFile := GetFileName(htmlFolder, templateFile)
		outputFile := GetFileName(htmlFolder, htmlFile)

		tmplBytes, err := ioutil.ReadFile(templateFile)
		if err != nil {
			return []byte{}, ErrReadTemplateFile
		}

		templateName := "accountInfo"
		template, err := template.New(templateName).Parse(string(tmplBytes))
		if err != nil {
			return []byte{}, ErrTemplateParse
		}

		output, err := os.Create(outputFile)
		if err != nil {
			return []byte{}, ErrCreateOutputFile
		}
		defer output.Close()

		err = template.Execute(output, email)
		if err != nil {
			return []byte{}, ErrTemplateExecute
		}

		htmlBytes, err := ioutil.ReadFile(outputFile)
		if err != nil {
			return []byte{}, ErrReadFile
		}

		return htmlBytes, nil
	}
}

// MakeProcessTransactions  creates a ProcessTransactions function
func MakeProcessTransactions(readCSV ReadCSV) ProcessTransactions {
	return func(ctx context.Context) (Email, error) {
		var email Email

		transactions, err := readCSV(ctx, GetFileName(folder, file))
		if err != nil {
			return Email{}, ErrCantGetCsvFile
		}

		email.Balance, email.AverageDebit, email.AverageCredit = getBalanceInfo(transactions)
		email.WorkingMonths = transactionsPerMonth(transactions)

		return email, nil
	}
}

// GetFileName returns the absolute file path of a file
func GetFileName(folder string, file string) string {
	_, filename, _, _ := runtime.Caller(0)
	testDir := filepath.Dir(filename)

	return filepath.Join(testDir, folder, file)
}
