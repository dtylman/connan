package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewIndexer(t *testing.T) {
	i, err := NewIndexer("/tmp/connan.db")
	if err != nil {
		t.Fatal(err)
	}
	defer i.Close()
	err = i.Start("/tmp")
	assert.Nil(t, err)
	for i.worker.IsRunning() {
		time.Sleep(time.Second)
		t.Log("worker is running...")
	}
}
