package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path"
	"strings"
	"time"

	"github.com/olebedev/when"
	"github.com/olebedev/when/rules/common"
	"github.com/olebedev/when/rules/en"

	"github.com/duckbrain/notes/notebook"
)

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
		if result == nil {
			fmt.Println("Could not understand the date")
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
	notebook.DocumentsDir = path.Join(usr.HomeDir, "Documents")

	editor := os.Getenv("NOTES_EDITOR")
	if editor == "" {
		editor = os.Getenv("EDITOR")
	}
	if editor == "" {
		editor = os.Getenv("VISUAL")
	}
	if editor == "" {
		editor = "vi" // Default to vi if nothing else
	}

	var debug bool
	switch os.Getenv("DEBUG") {
	case "1", "true", "TRUE":
		debug = true
	default:
		debug = false
	}

	n := notebook.Notebook{Name: folder, Folder: folder, Editor: editor}
	err = n.Load()
	if err != nil {
		panic(err)
	}
	tag := n.FileTag(date)
	file := n.FilePath(fmt.Sprintf("%v-%v.md", folder, tag))

	if debug {
		njson, _ := json.MarshalIndent(n, "", "\t")
		fmt.Println(string(njson))
	}

	os.MkdirAll(n.FilePath(""), os.ModePerm)
	res, err := n.TemplateResult(date)
	if err == nil {
		flags := os.O_WRONLY | os.O_CREATE | os.O_EXCL
		f, err := os.OpenFile(file, flags, 0622)
		if err != nil {
			fmt.Println(err)
		} else {
			_, err = f.Write(res)
			if err != nil {
				fmt.Println(err)
			}
			f.Close()
		}
	} else {
		fmt.Println(err)
	}

	cmd := exec.Command(n.Editor, file)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
}
