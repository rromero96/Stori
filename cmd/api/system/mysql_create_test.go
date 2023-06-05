package system_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/rromero96/stori/cmd/api/system"
)

const queryCreateMock string = "INSERT INTO stori.transactions \\(id, date, transaction, type\\) VALUES \\(\\?, \\?, \\?, \\?\\)"

func TestMakeMySQLCreate_success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	mock.ExpectPrepare(queryCreateMock)
	mock.ExpectExec(queryCreateMock).WillReturnResult(sqlmock.NewResult(1, 2))

	got := system.MakeMySQLCreate(db)

	assert.NotNil(t, got)
}

func TestMySQLCreate_success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	mock.ExpectPrepare(queryCreateMock)
	mock.ExpectExec(queryCreateMock).WillReturnResult(sqlmock.NewResult(1, 2))
	transactions := []system.Transaction{system.MockTransaction(0, time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), "credit", +60.5)}

	mysqlCreate := system.MakeMySQLCreate(db)
	ctx := context.Background()

	got := mysqlCreate(ctx, transactions)

	assert.Nil(t, got)
}

func TestMySQLCreate_failsWhenCantPrepareStatement(t *testing.T) {
	db, mock, _ := sqlmock.New()
	mock.ExpectPrepare("invalid statement")
	mock.ExpectExec(queryCreateMock).WillReturnResult(sqlmock.NewResult(1, 2))
	transactions := []system.Transaction{system.MockTransaction(0, time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), "credit", +60.5)}

	mysqlCreate := system.MakeMySQLCreate(db)
	ctx := context.Background()

	want := system.ErrCantPrepareStatement
	got := mysqlCreate(ctx, transactions)

	assert.Equal(t, want, got)
}

func TestMySQLCreate_failsWhenCantRunQuery(t *testing.T) {
	db, mock, _ := sqlmock.New()
	mock.ExpectPrepare(queryCreateMock)
	mock.ExpectExec(queryCreateMock).WillReturnError(errors.New("some error"))
	transactions := []system.Transaction{system.MockTransaction(0, time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), "credit", +60.5)}

	mysqlCreate := system.MakeMySQLCreate(db)
	ctx := context.Background()

	want := system.ErrCantRunQuery
	got := mysqlCreate(ctx, transactions)

	assert.Equal(t, want, got)
}
