package indexer

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/blevesearch/bleve"
)

var context struct {
	running    bool
	opened     bool
	libFolder  string
	blevePath  string
	bleveIndex bleve.Index
}

//Open opens indexer
func Open(libFolder string) error {
	if context.opened {
		return errors.New("Indexer already opened, close first")
	}
	context.blevePath = filepath.Join(libFolder, "connan.db")
	_, err := os.Stat(context.blevePath)
	if err == nil {
		log.Printf("Opening index at '%v'", context.blevePath)
		context.bleveIndex, err = bleve.Open(context.blevePath)
	} else {
		log.Printf("Creating index at '%v'", context.blevePath)
		context.bleveIndex, err = bleve.New(context.blevePath, bleve.NewIndexMapping())
	}
	if err == nil {
		context.opened = true
	}
	return err
}

//Close closes indexer
func Close() error {
	if !context.opened {
		return nil
	}
	context.opened = false
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
	log.Println(path)
	return nil
}
