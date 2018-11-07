package db

import (
	"log"
	"os"

	"github.com/blevesearch/bleve"
)

//Analyzer ...
type Analyzer interface {
	Process(path string, doc *Document) error
	Name() string
}

//DB database instance
type DB struct {
	analyzers []Analyzer
	Bleve     bleve.Index
	Queue     Queue
}

//Open opens the database
func Open(path string) (*DB, error) {
	db := new(DB)
	log.Printf("Opening db at '%v'", path)
	_, err := os.Stat(path)
	if err == nil {
		db.Bleve, err = bleve.Open(path)
	} else {
		db.Bleve, err = bleve.New(path, bleve.NewIndexMapping())
	}
	if err != nil {
		return nil, err
	}
	db.analyzers = make([]Analyzer, 0)
	return db, nil
}

//Close closes the database
func (db *DB) Close() error {
	return db.Bleve.Close()
}

//AddDocumentAnalyzer ...
func (db *DB) AddDocumentAnalyzer(a Analyzer) {
	db.analyzers = append(db.analyzers, a)
}

//NewDocument ...
func (db *DB) NewDocument(path string) (*Document, error) {
	doc := new(Document)
	for _, a := range db.analyzers {
		err := a.Process(path, doc)
		if err != nil {
			log.Printf("%v failed on '%v': %v", a.Name(), path, err)
		}
	}
	return doc, nil
}

//Save saves the doc in the DB
func (db *DB) Save(doc *Document) error {
	log.Printf("Saving '%v'", doc)
	return db.Bleve.Index(doc.ID, doc)
}
