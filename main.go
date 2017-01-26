package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"time"
	"path"

	"github.com/olebedev/when"
	"github.com/olebedev/when/rules/common"
	"github.com/olebedev/when/rules/en"
)

var DocumentsDir string

func main() {
	if len(os.Args) < 2 {
		fmt.Println("You must have a folder name")
		return
	}
	folder := os.Args[1]
	var date time.Time

	if len(os.Args) > 2 {
		w := when.New(nil)
		w.Add(en.All...)
		w.Add(common.All...)
		dateString := strings.Join(os.Args[2:], " ")
		result, err := w.Parse(dateString, time.Now())
		if err != nil {
			fmt.Println(err)
			return
		}
		date = result.Time
	} else {
		date = time.Now()
	}

	usr, err := user.Current()
	if err != nil {
		fmt.Println("Error finding home directory")
		return
	}
	DocumentsDir = path.Join(usr.HomeDir, "Documents")

	notebook := Notebook{Name:folder}
	err = notebook.Load()
  if err != nil {
    panic(err)
  }
	tag := notebook.FileTag(date)
	file := notebook.FilePath(fmt.Sprintf("%v-%v.md", folder, tag))

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = os.Getenv("VISUAL")
	}
	if editor == "" {
		editor = "vim"
	}

  //TODO: Mkdir and template file

	cmd := exec.Command(editor, file)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
}
