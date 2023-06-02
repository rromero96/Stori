package system_test

import (
	"context"
	"testing"

	"github.com/rromero96/stori/cmd/api/system"

	"github.com/stretchr/testify/assert"
)

func TestReadCSV_success(t *testing.T) {
	filename := system.GetFileName()

	readFiles := system.MakeReadCSV()
	ctx := context.Background()

	want := system.MockTransactions()
	got, err := readFiles(ctx, filename)

	assert.Nil(t, err)
	assert.Equal(t, got, want)
}

func TestReadCSV_fails(t *testing.T) {
	readFiles := system.MakeReadCSV()
	ctx := context.Background()

	want := system.ErrOpeningCsv
	_, got := readFiles(ctx, "")

	assert.Equal(t, got, want)
}
