package main

import "fmt"
import "os"
import "os/exec"
import "os/user"
import "time"
import "strings"
import "github.com/olebedev/when"
import "github.com/olebedev/when/rules/en"
import "github.com/olebedev/when/rules/common"

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

	dateString := date.Format("2006-01-02")
	fileName := fmt.Sprintf("%v/Documents/%v/%v-%v.md", usr.HomeDir, folder, folder, dateString)

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = os.Getenv("VISUAL")
	}
	if editor == "" {
		editor = "vim"
	}

	cmd := exec.Command(editor, fileName)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
}
