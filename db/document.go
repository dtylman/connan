package db

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/mapping"
	"github.com/iancoleman/strcase"
)

//Document ...
type Document struct {
	Path     string               `json:"path" storm:"id"`
	Analysis map[string]time.Time `json:"analysis"`
	Modified time.Time            `json:"modified"`
	Size     int64                `json:"size"`
	Fields   map[string]string    `json:"fields"`
}

//NewDocument ...
func NewDocument(path string, fileInfo os.FileInfo) (*Document, error) {
	if fileInfo.IsDir() {
		return nil, fmt.Errorf("'%v' is a folder", path)
	}
	doc := new(Document)
	doc.Fields = make(map[string]string)
	doc.Analysis = make(map[string]time.Time)
	doc.Path = path
	doc.UpdateFileInfo(fileInfo)
	return doc, nil
}

//UpdateFileInfo updates properties from file info
func (d *Document) UpdateFileInfo(fileInfo os.FileInfo) {
	d.Modified = fileInfo.ModTime()
	d.Size = fileInfo.Size()
}

//Analyze performs analysis
func (d *Document) Analyze(analyzers []Analyzer) bool {
	dirty := false
	for _, a := range analyzers {
		last, ok := d.Analysis[a.Name()]
		if !ok || (d.Modified.After(last)) {
			err := a.Process(d.Path, d)
			if err != nil {
				log.Printf("%v failed on '%v': %v", a.Name(), d.Path, err)
			}
			d.Analysis[a.Name()] = time.Now()
			dirty = true
		} else {
			log.Printf("Skipping %v: no change since last analysis on '%v'", a.Name(), d.Path)
		}
	}
	return dirty
}

func (d *Document) String() string {
	return fmt.Sprintf("%v (%v)", d.Path, d.Mime())
}

//SetField sets document field
func (d *Document) SetField(name string, value string) {
	val := reflect.ValueOf(d).Elem().FieldByName(strcase.ToCamel(name))
	if val.IsValid() {
		if val.Kind() != reflect.String {
			panic(fmt.Sprintf("cannot set %v=%v it is not string and already exists in struct", name, value))
		}
		val.SetString(strings.TrimSpace(value))
		return
	}
	d.Fields[strings.ToLower(name)] = strings.TrimSpace(value)
}

//GetField gets a field value or "" if no such field
func (d *Document) GetField(name string) string {
	n := strings.ToLower(name)
	val := reflect.ValueOf(d).Elem().FieldByName(name)
	if val.IsValid() {
		return val.String()
	}
	return d.Fields[n]
}

//Content returns the doc content
func (d *Document) Content() string {
	return d.GetField("content")
}

//SetContent sets content
func (d *Document) SetContent(content string) {
	d.SetField("content", content)
}

//Mime returns the mime-typeo of the doc
func (d *Document) Mime() string {
	return d.GetField("mime")
}

//Type returns doc type
func (d *Document) Type() string {
	return d.GetField("type")
}

//IsImage returns true if document is an image
func (d *Document) IsImage() bool {
	return strings.HasPrefix(d.GetField("mime"), "image")
}

//Name returns the document name
func (d *Document) Name() string {
	return filepath.Base(d.Path)
}

//BleveType returns bleve type
func (d Document) BleveType() string {
	return "document"
}

func (d *Document) bleveMapping() *mapping.DocumentMapping {
	dm := bleve.NewDocumentStaticMapping()
	dm.AddFieldMappingsAt("fields", bleve.NewTextFieldMapping())
	dm.AddFieldMappingsAt("modified", bleve.NewDateTimeFieldMapping())
	dm.AddFieldMappingsAt("path", bleve.NewTextFieldMapping())
	dm.AddFieldMappingsAt("size", bleve.NewNumericFieldMapping())
	return dm
}
