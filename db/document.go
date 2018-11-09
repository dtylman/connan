package db

import (
	"fmt"
	"time"
)

//Document ...
type Document struct {
	Path     string    `json:"path" storm:"id"`
	Content  string    `json:"content"`
	Mime     string    `json:"mime"`
	Analyzed time.Time `json:"analyzed"`
	Created  time.Time `json:"created"`
}

func (d *Document) String() string {
	return fmt.Sprintf("%v (%v)", d.Path, d.Mime)
}
