package notebook

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// Finds a notebook based on the name. Returns an error if the name is not
// found or note specific enough to limit to one.
//
// Search matches on the first few characters, similar to how git matches
// commit hashs, if the full name is not provided, the first characters
// are allowed, as long as there are no duplicates
func Search(text string) (Notebook, error) {
	n := Notebook{}
	text = strings.ToLower(text)
	notebooks, err := All()
	if err != nil {
		return n, err
	}
	matchCount := 0
	for _, notebook := range notebooks {
		name := strings.ToLower(notebook.Name)
		title := strings.ToLower(notebook.Title)
		if strings.Index(name, text) == 0 || strings.Index(title, text) == 0 {
			n = notebook
			matchCount++
			if name == text {
				return n, nil
			}
		}
	}
	if matchCount > 1 {
		return n, fmt.Errorf("Name not specific enough. Matches %v notebooks", matchCount)
	}
	if matchCount < 1 {
		return n, fmt.Errorf("Did not find matching notebook")
	}
	return n, nil
}

// Finds all notebooks that can be used and returns them
func All() ([]Notebook, error) {
	log.Printf("ls %v", DocumentsDir)
	files, err := ioutil.ReadDir(DocumentsDir)
	if err != nil {
		return nil, err
	}

	notebooks := make([]Notebook, 0)
	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		name := file.Name()
		n := Defaults
		log.Printf("open notebook %v", name)
		err := n.Load(name)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return nil, err
		}
		notebooks = append(notebooks, n)
		notebooks = append(notebooks, n.Notebooks...)
	}
	return notebooks, nil
}
