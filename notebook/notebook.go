package notebook

import (
	"bytes"
	"fmt"
	"io/ioutil"
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
	return path.Join(DocumentsDir, n.Name, p)
}

// Loads the configuration for the notebook.
//
// Load a global configuration in the home directory before overwriting
// values with the config file currently used, finally loading the
// configuration parameters that are default if not set.
func (n *Notebook) Load(name string) error {
	n.Name = name

	if n.Name == "" {
		return fmt.Errorf("Cannot load notebooks without a name")
	}

	// Load the global configuration
	configFile, err := ioutil.ReadFile(path.Join(DocumentsDir, ".notes"))
	if err == nil {
		err = frontmatter.Unmarshal(configFile, n)
		if err != nil {
			return err
		}
	}

	// Load the configuration for this notebook
	configFile, err = ioutil.ReadFile(n.filePath(".notes"))
	if err != nil {
		return err
	}
	err = frontmatter.Unmarshal(configFile, n)
	if err != nil {
		return err
	}

	// Load any unset values from defaults
	n.loadDefaults()

	return nil
}

// Loads default values into needed fields if they are not set.
func (n *Notebook) loadDefaults() {
	if n.Folder == "" {
		n.Folder = n.Name
	}
	if n.Title == "" {
		n.Title = n.Name
	}
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
