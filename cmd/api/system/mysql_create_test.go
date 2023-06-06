package system_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rromero96/stori/cmd/api/system"
	"github.com/stretchr/testify/assert"
)

const (
	queryCreateMock string = "INSERT INTO stori.transactions \\(id, date, transaction, type\\) VALUES \\(\\?, \\?, \\?, \\?\\)"
	queryFindMock   string = "SELECT MAX(id) FROM stori.transactions"
)

func TestMakeMySQLCreate_success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	mysqlFindMock := system.MockMySQLFind(-1, nil)
	mock.ExpectPrepare(queryCreateMock)
	mock.ExpectExec(queryCreateMock).WillReturnResult(sqlmock.NewResult(1, 2))

	got := system.MakeMySQLCreate(db, mysqlFindMock)

	assert.NotNil(t, got)
}

func TestMySQLCreate_success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	mysqlFindMock := system.MockMySQLFind(-1, nil)
	mock.ExpectPrepare(queryCreateMock)
	mock.ExpectExec(queryCreateMock).WillReturnResult(sqlmock.NewResult(1, 2))
	transactions := []system.Transaction{system.MockTransaction(0, time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), "credit", +60.5)}

	mysqlCreate := system.MakeMySQLCreate(db, mysqlFindMock)
	ctx := context.Background()

	got := mysqlCreate(ctx, transactions)

	assert.Nil(t, got)
}

func TestMySQLCreate_failsWhenMySQLFindThrowsError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	mysqlFindMock := system.MockMySQLFind(0, system.ErrCantRunQuery)
	mock.ExpectPrepare(queryCreateMock)
	mock.ExpectExec(queryCreateMock).WillReturnResult(sqlmock.NewResult(1, 2))
	transactions := []system.Transaction{system.MockTransaction(0, time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), "credit", +60.5)}

	mysqlCreate := system.MakeMySQLCreate(db, mysqlFindMock)
	ctx := context.Background()

	want := system.ErrCantGetLastID
	got := mysqlCreate(ctx, transactions)

	assert.Equal(t, want, got)
}

func TestMySQLCreate_failsWhenCantPrepareStatement(t *testing.T) {
	db, mock, _ := sqlmock.New()
	mysqlFindMock := system.MockMySQLFind(-1, nil)
	mock.ExpectPrepare("invalid statement")
	mock.ExpectExec(queryCreateMock).WillReturnResult(sqlmock.NewResult(1, 2))
	transactions := []system.Transaction{system.MockTransaction(0, time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), "credit", +60.5)}

	mysqlCreate := system.MakeMySQLCreate(db, mysqlFindMock)
	ctx := context.Background()

	want := system.ErrCantPrepareStatement
	got := mysqlCreate(ctx, transactions)

	assert.Equal(t, want, got)
}

func TestMySQLCreate_failsWhenCantRunQuery(t *testing.T) {
	db, mock, _ := sqlmock.New()
	mysqlFindMock := system.MockMySQLFind(-1, nil)
	mock.ExpectPrepare(queryCreateMock)
	mock.ExpectExec(queryCreateMock).WillReturnError(errors.New("some error"))
	transactions := []system.Transaction{system.MockTransaction(0, time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), "credit", +60.5)}

	mysqlCreate := system.MakeMySQLCreate(db, mysqlFindMock)
	ctx := context.Background()

	want := system.ErrCantRunQuery
	got := mysqlCreate(ctx, transactions)

	assert.Equal(t, want, got)
}

func TestMakeMySQLFind_success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	mock.ExpectPrepare(queryFindMock)
	mock.ExpectExec(queryFindMock).WillReturnResult(sqlmock.NewResult(1, 2))

	got := system.MakeMySQLFind(db)

	assert.NotNil(t, got)
}

func TestMySQLFind_success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	rows := mock.NewRows([]string{"MAX(id)"}).AddRow(10)
	mock.ExpectQuery(queryFindMock).WillReturnRows(rows)
	ctx := context.Background()

	mysqlFind := system.MakeMySQLFind(db)

	want := int64(-1)
	got, err := mysqlFind(ctx)

	assert.Nil(t, err)
	assert.Equal(t, got, want)
}
