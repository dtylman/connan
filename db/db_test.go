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
	defer db.Close()
	doc, err := db.NewDocument("/tmp/connan.db/index_meta.json")
	assert.Nil(t, err)
	assert.NotEmpty(t, doc.Mime)
}
