
# Notes

This is a simple launcher that opens an editor for editing notes (currently markdown only). It handles opening multiple notebooks by date or week with a template for each new entry. I made this to make it easier to quickly take notes in my classes in markdown with a header in place.

My project [atom-journal](https://github.com/duckbrain/atom-journal) is roughly the same thing, but that one is an atom plugin, and this one is more convenient for editing with other editors like vim and typora.

## Getting Started

### Install

You can install this by installing [Go](https://golang.org/) and running

```
go get -u github.com/duckbrain/notes
```

This will install a binary to `$GOPATH/bin/notes`

### Setup

Before you can use Notes, you will need to setup your documents directory. By default, notes uses `~/Documents/` as your documents directory. It assumes that you have different "notebooks" sorted into subdirectories of the documents directory. You may have as many other files as you like in that directory and it's subdirectory.

Each subdirectory that you would like to use as a notebook should contain a `.notes` file. This file is a template with YAML frontmatter. You can define values to setup several varieties of note taking systems.

Below is a simplified example of a documents directory setup with the notes.

```bash
~/Documents $ tree -a
.
├── church
│   ├── church-2015-03-08.md
│   ├── church-2017-02-26.md
│   ├── church-2017-02-27.md
│   ├── church-2017-03-05.md
│   ├── church-2017-03-12.md
│   ├── general-conference-2015-april.md
│   ├── general-conference-2015-october.md
│   ├── general-conference-2016-april.md
│   ├── general-conference-2016-october.md
│   ├── .notes
├── geol102
│   ├── geol102-00.md
│   ├── geol102-01.md
│   ├── geol102-02.md
│   ├── geol102-03.md
│   ├── geol102-04.md
│   ├── geol102-05.md
│   ├── geol102-06.md
│   ├── geol102-07.md
│   ├── .notes
│   └── radiometric-dating.md
├── journal
│   ├── dream-journal-2017-02-23.md
│   ├── journal-2013-03-31.md
│   ├── journal-2016-11-13.md.gpg
│   ├── journal-2016-11-19.md
│   ├── journal-2016-11-20.md
│   ├── journal-2016-11-20.md.gpg
│   ├── journal-2017-02-26.md
│   ├── journal-2017-03-14.md
│   ├── .notes
│   ├── study-journal-2016-10-09.md
│   ├── study-journal-2016-10-19.md
├── stat411
│   ├── .notes
│   ├── stat411-00.md
│   ├── stat411-01.md
│   ├── stat411-02.md
│   ├── stat411-03-figure1.png
│   ├── stat411-03.md
│   ├── stat411-04.md
│   ├── stat411-05.md
│   ├── stat411-06.md
│   ├── stat411-07.md
│   ├── stat411-08-figure1.png
│   ├── stat411-08-figure2.png
│   ├── stat411-08.md
│   └── stat411-formulas.md
└── todo.md
```

The `.notes` file contains a template and data for the notebook. A directory can contain more than one category of notes as seen in the "journal" directory up above. The `.notes` file will for the "stat411" directory looks like this

```markdown
---
title: STAT 411
weekStart: 3
---

# {{.Title}} - {{.Tuesday.Format "Monday 2 January 2006"}}

# {{.Title}} - {{.Thursday.Format "Monday 2 January 2006"}}
```

This renders a template like this `stat411-08.md`

```markdown

# STAT 411 - Tuesday 14 March 2017

# STAT 411 - Thursday 16 March 2017
```

The frontmatter at the top of the file defines properties that Notes uses to determine how to name the file, give data to the template, and other things. Everything below that is used as a [Go template](https://golang.org/pkg/text/template/) to generate the template contents.

A more complicated frontmatter is used in "journal" to get several notebooks in one folder

```markdown
---
title: Journal
notebooks:
  - title: Dream Journal
    name: dream-journal
  - title: Study Journal
    name: study-journal
---

# {{.Title}} - {{.Date.Format "Monday 2 January 2006"}}

```

### Usage

Notes can be run without any arguments to list the notebooks that it detects and can use.

```bash
~ $ notes
church
geol102
journal
dream-journal
study-journal
stat41
```

By supplying a notebook name, it will open the notebook for today. You may supply the first unique characters instead of the full name, similar to Git hashes. In this example, I could supply "d" for dream-journal, but I would need "sta" for stat411. Passing extra parameters are used to parse a date to use instead of today.

```bash
~ $ notes journal
~ $ notes stat
~ $ notes geo next wednesday
```

The dates are parsed by [github.com/olebedev/when](https://github.com/olebedev/when), but it does not work as well as I would like for this situation, so I will likely be augmenting it with some custom parsing before passing it to when.

## Configuration

The following configuration options are available to use in the frontmatter

```go
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
}
```

The template is then rendered with the following data

```go
struct {
	Notebook // This means that it inherits all the fields above
	Date time.Time
	Tag  string

	Sunday, Monday, Tuesday, Wednesday, Thursday, Friday, Saturday time.Time
}
```

