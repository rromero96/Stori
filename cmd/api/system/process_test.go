package system_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rromero96/stori/cmd/api/system"
)

func TestMakeHTMLProcessTransactions_success(t *testing.T) {
	readCSVmock := system.MockReadCSV(system.MockTransactions(), nil)
	mysqlCreateMock := system.MockMySQLCreate(nil)

	got := system.MakeHTMLProcessTransactions(readCSVmock, mysqlCreateMock)

	assert.NotNil(t, got)
}

func TestHTMLProcessTransactions_success(t *testing.T) {
	readCSVmock := system.MockReadCSV(system.MockTransactions(), nil)
	mysqlCreateMock := system.MockMySQLCreate(nil)
	htmlProcessTransactions := system.MakeHTMLProcessTransactions(readCSVmock, mysqlCreateMock)
	ctx := context.Background()

	got, err := htmlProcessTransactions(ctx)

	assert.Nil(t, err)
	assert.NotNil(t, got)
}

func TestHTMLProcessTransactions_failsWhenReadCSVThrowsError(t *testing.T) {
	readCSVmock := system.MockReadCSV(nil, system.ErrOpeningCsv)
	mysqlCreateMock := system.MockMySQLCreate(nil)
	htmlProcessTransactions := system.MakeHTMLProcessTransactions(readCSVmock, mysqlCreateMock)
	ctx := context.Background()

	want := system.ErrCantGetCsvFile
	_, got := htmlProcessTransactions(ctx)

	assert.Equal(t, want, got)
}

func TestHTMLProcessTransactions_failsWhenMySQLCreateThworsError(t *testing.T) {
	readCSVmock := system.MockReadCSV(system.MockTransactions(), nil)
	mysqlCreateMock := system.MockMySQLCreate(system.ErrCantPrepareStatement)
	htmlProcessTransactions := system.MakeHTMLProcessTransactions(readCSVmock, mysqlCreateMock)
	ctx := context.Background()

	want := system.ErrCantCreateTransactions
	_, got := htmlProcessTransactions(ctx)

	assert.Equal(t, want, got)
}
