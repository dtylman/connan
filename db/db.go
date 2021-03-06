package db

import (
	"log"
	"os"
	"path/filepath"

	"github.com/asdine/storm"
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
	Storm     *storm.DB
	Queue     Queue
}

//Open opens the database
func Open(path string) (*DB, error) {
	db := new(DB)
	log.Printf("Opening db at '%v'", path)
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return nil, err
	}
	db.Storm, err = storm.Open(filepath.Join(path, "storm.db"))
	if err != nil {
		return nil, err
	}
	blevepath := filepath.Join(path, "bleve")
	_, err = os.Stat(blevepath)
	if err == nil {
		db.Bleve, err = bleve.Open(blevepath)
	} else {
		d := Document{}
		im := bleve.NewIndexMapping()
		im.AddDocumentMapping(d.BleveType(), d.bleveMapping())
		db.Bleve, err = bleve.New(blevepath, im)
	}
	if err != nil {
		return nil, err
	}
	db.analyzers = make([]Analyzer, 0)
	return db, nil
}

//Close closes the database
func (db *DB) Close() error {
	log.Printf("Closing Storm %v", db.Storm.Bolt.Path())
	err := db.Storm.Close()
	if err != nil {
		return err
	}
	log.Println("Closing Bleve")
	return db.Bleve.Close()
}

//AddDocumentAnalyzer ...
func (db *DB) AddDocumentAnalyzer(a Analyzer) {
	db.analyzers = append(db.analyzers, a)
}

//Save saves the doc in the DB
func (db *DB) Save(doc *Document) error {
	log.Printf("Saving '%v'", doc)
	err := db.Storm.Save(doc)
	if err != nil {
		return nil
	}
	return db.Bleve.Index(doc.Path, doc)
}

//DocumentExists is true if document already exists
func (db *DB) DocumentExists(path string) bool {
	var doc Document
	err := db.Storm.One("Path", path, &doc)
	return err == nil
}

//Document returns a document for path or nil if not found
func (db *DB) Document(path string) (*Document, error) {
	var doc Document
	err := db.Storm.One("Path", path, &doc)
	if err != nil {
		return nil, err
	}
	return &doc, nil
}
