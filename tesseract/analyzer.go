package tesseract

import (
	"log"

	"github.com/dtylman/connan/db"
)

//Analyzer is a tesseract OCR analyzer
type Analyzer struct {
}

//NewAnalyzer creates a new analyzer
func NewAnalyzer() *Analyzer {
	return new(Analyzer)
}

//Process ...
func (a *Analyzer) Process(path string, doc *db.Document) error {
	log.Printf("Processing %v", path)
	return nil
}

//Name ...
func (a *Analyzer) Name() string {
	return "tesseract"
}
