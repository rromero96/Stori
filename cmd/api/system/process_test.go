package system_test

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/rromero96/stori/cmd/api/system"

	"github.com/stretchr/testify/assert"
)

func TestReadCSV_success(t *testing.T) {
	// Get the current file's directory
	_, filename, _, _ := runtime.Caller(0)
	testDir := filepath.Dir(filename)
	// Construct the absolute file path to the CSV file
	filename = filepath.Join(testDir, "utils", "data.csv")

	readFiles := system.MakeReadCSV()
	ctx := context.Background()

	want := system.MockTransactions()
	got, err := readFiles(ctx, filename)

	assert.Nil(t, err)
	assert.Equal(t, got, want)
}
