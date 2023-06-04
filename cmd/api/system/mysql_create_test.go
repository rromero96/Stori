package system_test

import (
	"testing"

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
