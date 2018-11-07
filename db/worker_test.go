package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorker_Stop(t *testing.T) {
	db, err := Open("/tmp/connan.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	assert.NoError(t, err)
	w := NewWorker(db)
	w.db.Queue.Add("/tmp")
	w.Start()
}
