package main

import (
	"flag"
	"os"
	"os/user"
	"path"

	"github.com/duckbrain/notes/notebook"
)

var Debug bool
var Help bool
var List bool

func init() {
	home := ""
	usr, err := user.Current()
	if err == nil {
		home = path.Join(usr.HomeDir, "Documents")
	}
	var editor string
	os.Getenv("NOTES_EDITOR")
	if editor == "" {
		editor = os.Getenv("EDITOR")
	}
	if editor == "" {
		editor = os.Getenv("VISUAL")
	}
	if editor == "" {
		editor = "vi" // Default to vi if nothing else
	}
	flag.BoolVar(&Debug, "debug", false, "Print out debugging information")
	flag.BoolVar(&List, "list", false, "List the notebooks found")
	flag.StringVar(&notebook.DocumentsDir, "docs", home, "`Directory` where the documents are stored")
	flag.StringVar(&notebook.Defaults.Editor, "editor", editor, "Editor to open documents in")
}
