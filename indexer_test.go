package main

import (
	"testing"
	"time"

	"github.com/dtylman/connan/db"
	"github.com/stretchr/testify/assert"
)

func TestNewIndexer(t *testing.T) {
	db, err := db.Open("/tmp/connan.db")
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, db)
	defer func() {
		err := db.Close()
		assert.NoError(t, err)
	}()
	i, err := NewIndexer(db)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, i)
	defer func() {
		err := i.Close()
		assert.NoError(t, err)
	}()
	err = i.Start("/tmp")
	assert.Nil(t, err)
	for i.worker.IsRunning() {
		time.Sleep(time.Second)
		t.Log("worker is running...")
	}
}
