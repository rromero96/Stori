package system_test

import (
	"context"
	"errors"
	"testing"

	"github.com/rromero96/stori/cmd/api/system"

	"github.com/stretchr/testify/assert"
)

func TestReadCSV_success(t *testing.T) {
	filename := system.GetFileName("data", "data.csv")

	mySQLCreateMock := system.MockMySQLCreate(nil)
	readFiles := system.MakeReadCSV(mySQLCreateMock)
	ctx := context.Background()

	want := system.MockTransactions()
	got, err := readFiles(ctx, filename)

	assert.Nil(t, err)
	assert.Equal(t, got, want)
}

func TestReadCSV_failsWhenCantOpenCsvFile(t *testing.T) {
	mySQLCreateMock := system.MockMySQLCreate(nil)
	readFiles := system.MakeReadCSV(mySQLCreateMock)
	ctx := context.Background()

	want := system.ErrOpeningCsv
	_, got := readFiles(ctx, "")

	assert.Equal(t, got, want)
}

func TestReadCSV_failsWhenMySQLCreateThrowsError(t *testing.T) {
	filename := system.GetFileName("data", "data.csv")
	err := errors.New("error")

	mySQLCreateMock := system.MockMySQLCreate(err)
	readFiles := system.MakeReadCSV(mySQLCreateMock)
	ctx := context.Background()

	want := system.ErrCantCreateTransactions
	_, got := readFiles(ctx, filename)

	assert.Equal(t, got, want)
}
