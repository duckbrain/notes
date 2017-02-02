package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path"
	"text/template"
	"time"

	"github.com/ericaro/frontmatter"
)

type Notebook struct {
	Name   string
	Folder string

	// Template content
	Template string `fm:"content"`

	// The week number of the first week of the weekly iteration
	WeekStart int `yaml:"weekStart"`

	Editor string
}

func (n Notebook) FilePath(p string) string {
	return path.Join(DocumentsDir, n.Name, p)
}

func (n *Notebook) Load() error {
	if n.Name == "" {
		return fmt.Errorf("Cannot load notebooks without a name")
	}
	if n.Folder == "" {
		n.Folder = n.Name
	}

	configFile, err := ioutil.ReadFile(n.FilePath(".notes"))
	if err == nil {
		err = frontmatter.Unmarshal(configFile, n)
		if err != nil {
			return err
		}
	}

	return nil
}

func (n Notebook) FileTag(date time.Time) string {
	if n.WeekStart != 0 {
		_, week := date.ISOWeek()
		week -= n.WeekStart
		return fmt.Sprintf("%02v", week)
	}

	return date.Format("2006-01-02")
}

func (n Notebook) TemplateResult(date time.Time) ([]byte, error) {
	t := template.New(n.Name)
	t, err := t.Parse(n.Template)
	if err != nil {
		return nil, err
	}
	buffer := &bytes.Buffer{}
	err = t.Execute(buffer, struct {
		Notebook
		Date time.Time
	}{
		n,
		date,
	})
	return buffer.Bytes(), err
}
