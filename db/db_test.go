package db

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpen(t *testing.T) {
	db, err := Open("/tmp/connan.db")
	assert.Nil(t, err)
	assert.NotNil(t, db)
	err = db.Close()
	assert.Nil(t, err)
}

func TestQueue(t *testing.T) {
	db, err := Open("/tmp/connan.db")
	assert.Nil(t, err)
	assert.NotNil(t, db)
	defer db.Close()
	db.Queue.Add("/tmp/lala")
	assert.EqualValues(t, "/tmp/lala", db.Queue.items[0])
}

func TestDB_NewDocument(t *testing.T) {
	db, err := Open("/tmp/connan.db")
	assert.NoError(t, err)
	assert.NotNil(t, db)
	defer db.Close()
	path := "db_test.go"
	fileInfo, err := os.Stat(path)
	assert.NoError(t, err)
	doc, err := NewDocument(path, fileInfo)
	assert.Nil(t, err)
	err = db.Save(doc)
	assert.NoError(t, err)
	assert.True(t, db.DocumentExists(path))
}

func TestDB_Save(t *testing.T) {
	os.RemoveAll("/tmp/connand.db")
	db, err := Open("/tmp/connand.db")
	assert.NoError(t, err)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	bc, err := db.Bleve.DocCount()
	assert.NoError(t, err)
	assert.EqualValues(t, 0, bc)
	path := "db_test.go"
	fileInfo, err := os.Stat(path)
	assert.NoError(t, err)
	doc, err := NewDocument(path, fileInfo)
	assert.NoError(t, err)
	err = db.Save(doc)
	assert.NoError(t, err)
	bc, err = db.Bleve.DocCount()
	assert.NoError(t, err)
	assert.EqualValues(t, 1, bc)
	err = db.Save(doc)
	bc, err = db.Bleve.DocCount()
	assert.NoError(t, err)
	assert.EqualValues(t, 1, bc)
	sc, err := db.Storm.Count(&Document{})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, sc)
}
