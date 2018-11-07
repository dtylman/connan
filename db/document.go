package db

//Document ...
type Document struct {
	ID   string
	Path string `json:"path"`
}

func (d *Document) String() string {
	return d.Path
}
