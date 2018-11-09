package db

import (
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
	doc, err := db.NewDocument(path)
	assert.Nil(t, err)
	assert.NotEmpty(t, doc.Mime)
	err = db.Save(doc)
	assert.NoError(t, err)
	assert.True(t, db.DocumentExists(path))

}
