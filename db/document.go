package db

import (
	"fmt"
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

func (d *Document) String() string {
	return fmt.Sprintf("%v (%v)", d.Path, d.Type())
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
