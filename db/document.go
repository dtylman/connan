package db

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/iancoleman/strcase"
)

//Document ...
type Document struct {
	Path     string            `json:"path" storm:"id"`
	Analyzed time.Time         `json:"analyzed"`
	Modified time.Time         `json:"modified"`
	Size     int64             `json:"size"`
	Fields   map[string]string `json:"fields"`
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
	return true
	//return strings.HasPrefix(d.GetField("mime"), "image")
}
