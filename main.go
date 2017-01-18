package main

import "fmt"
import "os"
import "strings"
import "github.com/olebedev/when"
import "github.com/olebedev/when/rules/en"
import "github.com/olebedev/when/rules/common"

func main() {
	w := when.New(nil)
	w.Add(en.All...)
	w.Add(common.All...)

	if len(os.Args) < 2 {
		fmt.Println("You must have a folder name")
	}
	folderName := os.Args[1]
	var dateString string

	if len(os.Args) > 2 {
		dateString = strings.Join(os.Args[2:], " ")
	} else {
		dateString = "today"
	}

	fmt.Printf("Folder: %v, Date: %v\n", folderName, dateString)

}
