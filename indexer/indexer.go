package indexer

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/mapping"
)

var context struct {
	running    bool
	opened     bool
	libFolder  string
	blevePath  string
	mapping    mapping.IndexMapping
	bleveIndex bleve.Index
}

//Open opens indexer
func Open(libFolder string) error {
	if context.opened {
		return errors.New("Indexer already opened, close first")
	}
	context.blevePath = filepath.Join(libFolder, "connan.db")
	log.Printf("Indexing to %v", context.blevePath)
	context.mapping = bleve.NewIndexMapping()
	var err error
	context.bleveIndex, err = bleve.New(context.blevePath, context.mapping)
	return err
}

//Close closes indexer
func Close() error {
	if !context.opened {
		return nil
	}
	return context.bleveIndex.Close()
}

//Start starts indexing
func Start() error {
	if !context.opened {
		return errors.New("Indexer Not Opened")
	}
	if context.running {
		return errors.New("Indexer Allready Running")
	}
	context.running = true
	defer func() {
		context.running = false
		log.Println("Indexing ended")
	}()
	return filepath.Walk(context.libFolder, walk)
}

func walk(path string, info os.FileInfo, err error) error {
	return nil
}
