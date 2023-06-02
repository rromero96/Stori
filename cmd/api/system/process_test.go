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
