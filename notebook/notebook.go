package notebook

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"text/template"
	"time"

	"github.com/ericaro/frontmatter"
)

// A notebook to create entries in.
//
// The fields in this struct are loaded from the front matter in the template files.
type Notebook struct {
	// The name in the file name. Defaults to the notebook name by search.
	Name string

	// Purely for the templating, represents the title to display at the top
	// of a template. Defaults to the value in Name.
	Title string

	// The name of the directory the notebook is located in. Often multiple
	// notebooks are grouped locally in one directory, but seperated by having
	// different names.
	Folder string

	// Template content
	Template string `fm:"content"`

	// The week number of the first week of the weekly iteration
	WeekStart int `yaml:"weekStart"`

	// The program to edit the notebook in
	Editor string

	// Sub-notebooks that inherit this notebook's folder and other properties
	Notebooks []Notebook

	// The days of the week this notebook can be used
	Weekdays []time.Weekday

	// If set, causes the notes to be encrypted using PGP with the specified ID
	PGPID string `yaml:"pgpid"`

	FileNameTemplate string `yaml:"filename"`
}

func (n Notebook) String() string {
	return n.Name
}

// Returns a file path of a document given by name.
func (n Notebook) filePath(p string) string {
	return path.Join(DocumentsDir, n.Folder, p)
}

// Loads the configuration for the notebook.
//
// Load a global configuration in the home directory before overwriting
// values with the config file currently used, finally loading the
// configuration parameters that are default if not set.
func (n *Notebook) Load(name string) error {
	log.Printf("editor from default %v", n.Editor)
	n.Name = name
	n.Folder = name
	n.Title = name

	if n.Name == "" {
		return fmt.Errorf("Cannot load notebooks without a name")
	}

	// Load the global configuration
	configFile, err := ioutil.ReadFile(path.Join(DocumentsDir, ".notes"))
	log.Printf("open main config %v", path.Join(DocumentsDir, ".notes"))
	if err == nil {
		err = frontmatter.Unmarshal(configFile, n)
		log.Printf("editor after main config %v", n.Editor)
		if err != nil {
			return err
		}
	}

	log.Printf("open directory config %v", n.filePath(".notes"))
	// Load the configuration for this notebook
	configFile, err = ioutil.ReadFile(n.filePath(".notes"))
	log.Printf("editor after directory config %v", n.Editor)
	if err != nil {
		return err
	}
	err = frontmatter.Unmarshal(configFile, n)
	if err != nil {
		return err
	}

	for i, c := range n.Notebooks {
		if len(c.Name) == 0 {
			return fmt.Errorf("Sub notebook needs a name")
		}
		if c.Name == n.Name {
			return fmt.Errorf("Sub notebook needs a name different from the parent")
		}
		if len(c.Title) == 0 {
			c.Title = c.Name
		}
		c.Folder = n.Folder
		c.Template = n.Template
		c.FileNameTemplate = n.FileNameTemplate
		c.WeekStart = n.WeekStart
		c.Editor = n.Editor
		if len(c.Weekdays) == 0 {
			c.Weekdays = n.Weekdays
		}
		if len(c.PGPID) == 0 {
			c.PGPID = n.PGPID
		}

		n.Notebooks[i] = c
	}

	return nil
}

func (n Notebook) runTmp(nameTmp, tmp string, date time.Time) ([]byte, error) {
	t := template.New(n.Name)
	t, err := t.Parse(tmp)
	if err != nil {
		return nil, err
	}
	buffer := &bytes.Buffer{}

	data := struct {
		Notebook
		Date time.Time
		Week int

		Sunday, Monday, Tuesday, Wednesday, Thursday, Friday, Saturday time.Time
	}{}

	_, week := date.ISOWeek()
	week -= n.WeekStart

	data.Notebook = n
	data.Date = date
	data.Week = week
	data.Sunday = date.AddDate(0, 0, -int(date.Weekday()))
	data.Monday = data.Sunday.AddDate(0, 0, 1)
	data.Tuesday = data.Sunday.AddDate(0, 0, 2)
	data.Wednesday = data.Sunday.AddDate(0, 0, 3)
	data.Thursday = data.Sunday.AddDate(0, 0, 4)
	data.Friday = data.Sunday.AddDate(0, 0, 5)
	data.Saturday = data.Sunday.AddDate(0, 0, 6)

	err = t.Execute(buffer, data)
	return buffer.Bytes(), err
}

// Gets a tag to distinguish a document for one day from another. Typically a
// variation of the date or the number of weeks that have passed since a
// specific week.
func (n Notebook) FileName(date time.Time) (string, error) {
	tmp := n.FileNameTemplate
	if len(tmp) == 0 {
		if n.WeekStart != 0 {
			tmp = "{{.Name}}-{{.Week}}.md"
		} else {
			tmp = `{{.Name}}-{{.Date.Format "2006-01-02"}}.md`
		}
	}
	res, err := n.runTmp("%v Filename template", tmp, date)
	return n.filePath(string(res)), err
}

// Renders the template with the values given
func (n Notebook) Render(date time.Time) ([]byte, error) {
	return n.runTmp("%v", n.Template, date)
}
