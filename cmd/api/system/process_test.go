package system_test

import (
	"context"
	"testing"

	"github.com/rromero96/stori/cmd/api/system"
	"github.com/stretchr/testify/assert"
)

func TestMakeProcessTransaccion_success(t *testing.T) {
	readCSV := system.MockReadCSV(system.MockTransactions(), nil)

	got := system.MakeProcessTransactions(readCSV)

	assert.NotNil(t, got)
}

func TestProcessTransaccion_success(t *testing.T) {
	readCSV := system.MockReadCSV(system.MockTransactions(), nil)
	processTransactions := system.MakeProcessTransactions(readCSV)
	ctx := context.Background()

	want := system.MockEmail()
	got, err := processTransactions(ctx)

	assert.Nil(t, err)
	assert.Equal(t, want, got)
}

func TestProcessTransaccion_failsWhenReadCSVThrowsError(t *testing.T) {
	readCSV := system.MockReadCSV(nil, system.ErrReadingCsv)
	processTransactions := system.MakeProcessTransactions(readCSV)
	ctx := context.Background()

	want := system.ErrCantGetCsvFile
	_, got := processTransactions(ctx)

	assert.Equal(t, want, got)
}

func TestMakeHTMLProcessTransactions_success(t *testing.T) {
	processTransactions := system.MockProcessTransactions(system.MockEmail(), nil)

	got := system.MakeHTMLProcessTransactions(processTransactions)

	assert.NotNil(t, got)
}

func TestHTMLProcessTransactions_success(t *testing.T) {
	processTransactions := system.MockProcessTransactions(system.MockEmail(), nil)
	htmlProcessTransactions := system.MakeHTMLProcessTransactions(processTransactions)
	ctx := context.Background()

	got, err := htmlProcessTransactions(ctx)

	assert.Nil(t, err)
	assert.NotNil(t, got)
}

func TestHTMLProcessTransactions_failsWhenProcessTransactionsThrowsError(t *testing.T) {
	processTransactions := system.MockProcessTransactions(system.Email{}, system.ErrCantGetCsvFile)
	htmlProcessTransactions := system.MakeHTMLProcessTransactions(processTransactions)
	ctx := context.Background()

	want := system.ErrCantGetTransactionInfo
	_, got := htmlProcessTransactions(ctx)

	assert.Equal(t, want, got)
}
