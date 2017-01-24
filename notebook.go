package main

import (
	"path"
	"time"

	"github.com/ericaro/frontmatter"
)

type Notebook struct {
	name     string
	folder   string
	template string
}

func (n Notebook) FilePath(path string) string {
	return path.Join(DocumentsDir, n.name, path)
}

func (n *Notebook) Load() {

}

func (n Notebook) FileTag(date time.Time) string {

}
