// This program provides a quick way to open notes for me. You are free to use
// it too.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
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

var (
	Debug   bool
	Help    bool
	List    bool
	GPGPath string
)

func init() {
	home := ""
	usr, err := user.Current()
	if err == nil {
		home = path.Join(usr.HomeDir, "Documents")
	}
	editor := os.Getenv("NOTES_EDITOR")
	if len(editor) == 0 {
		editor = os.Getenv("EDITOR")
	}
	if len(editor) == 0 {
		editor = os.Getenv("VISUAL")
	}
	if len(editor) == 0 {
		editor = "vi" // Default to vi if nothing else
	}
	flag.BoolVar(&Debug, "debug", false, "Print out debugging information")
	flag.BoolVar(&List, "list", false, "List the notebooks found")
	flag.StringVar(&notebook.DocumentsDir, "docs", home, "`Directory` where the documents are stored")
	flag.StringVar(&notebook.Defaults.Editor, "editor", editor, "Editor to open documents in")
	flag.StringVar(&GPGPath, "pgp", "gpg", "Path to GnuPG installed on the system")
}

func main() {
	flag.Parse()
	log.Println("default editor", notebook.Defaults.Editor)

	if !Debug {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}

	switch {
	case Help:
		flag.PrintDefaults()
	case List:
		printNotebooks()
	default:
		openDoc()
	}
}

func printNotebooks() {
	notebooks, err := notebook.All()
	handleErr(err)
	for _, n := range notebooks {
		fmt.Println(n)
	}
}

func openDoc() {
	args := flag.Args()

	if len(args) < 1 {
		printNotebooks()
		return
	}
	var date time.Time

	if len(args) > 1 {
		w := when.New(nil)
		w.Add(en.All...)
		w.Add(common.All...)
		dateString := strings.Join(args[1:], " ")
		result, err := w.Parse(dateString, time.Now())
		handleErr(err)
		if result == nil {
			handleErr(fmt.Errorf("Could not understand the date"))
		}
		date = result.Time
	} else {
		date = time.Now()
	}

	n, err := notebook.Search(args[0])
	handleErr(err)

	date, err = n.AllowedDate(date)
	handleErr(err)

	file, err := n.FileName(date)
	handleErr(err)

	if len(n.PGPID) > 0 {
		gpgfile := fmt.Sprintf("%v.gpg", file)
		if _, err := os.Stat(gpgfile); err == nil {
			flags := os.O_WRONLY | os.O_CREATE | os.O_EXCL
			f, err := os.OpenFile(file, flags, 0622)
			handleErr(err)

			cmd := exec.Command(GPGPath, "--decrypt", gpgfile)
			cmd.Stdin = os.Stdin
			cmd.Stdout = f
			cmd.Stderr = os.Stderr
			err = cmd.Run()
			if err != nil {
				f.Close()
				os.Remove(file)
			}
		}
	}

	if Debug {
		njson, _ := json.MarshalIndent(n, "", "\t")
		fmt.Println(date)
		fmt.Println(string(njson))
	}

	os.MkdirAll(path.Dir(file), os.ModePerm)
	res, err := n.Render(date)
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

	log.Printf("launch %v %v", n.Editor, file)
	cmd := exec.Command(n.Editor, file)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	handleErr(err)

	if len(n.PGPID) > 0 {
		cmd := exec.Command(GPGPath, "--encrypt", "--yes", "--recipient", n.PGPID, file)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		err = cmd.Run()
		handleErr(err)

		err = os.Remove(file)
		handleErr(err)
	}
}

func handleErr(err error) {
	if err == nil {
		return
	}
	if Debug {
		panic(err)
	} else {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
