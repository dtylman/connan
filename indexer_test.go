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
	i, err := NewIndexer(db)
	if err != nil {
		t.Fatal(err)
	}
	err = i.Start("/tmp")
	assert.Nil(t, err)
	for i.worker.IsRunning() {
		time.Sleep(5000)
		t.Log("worker is running...")
	}
}
