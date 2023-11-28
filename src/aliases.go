package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

const (
	ALIAS_FILE_ENV = "RIG_ALIAS_FILE"
	ALIAS_FILE_DEFAULT = "/tmp/rig-aliases"

	ALIAS_PREFIX_ENV = "RIG_ALIAS_PREFIX"
	ALIAS_PREFIX_DEFAULT = "e"

	EDITOR_CMD_ENV = "RIG_EDITOR"
	EDITOR_CMD_DEFAULT = "vim"

	EDITOR_CMD_FORMAT_ENV = "RIG_EDIT_COMMAND"
)

var (
	EDITOR_MAP = map[string]string{
		"code": `code --goto "{{ .Filename }}:{{ .LineNumber }}:{{ .ColumnNumber }}"`,
		"emacs": `emacs +{{ .LineNumber }}:{{ .ColumnNumber }} --file="{{ .Filename }}"`,
		"hx": `hx "{{ .Filename }}:{{ .LineNumber }}:{{ .ColumnNumber }}"`,
		"mcedit": `mcedit "{{ .Filename }}:{{ .LineNumber }}"`,
		"micro": `micro +{{ .LineNumber }}:{{ .ColumnNumber }} "{{ .Filename }}"`,
		"nano": `nano +{{.LineNumber}},{{ .ColumnNumber }} "{{ .Filename }}"`,
		"ne": `ne +{{.LineNumber}},{{ .ColumnNumber }} "{{ .Filename }}"`,
		"nvim": `nvim -c "call cursor({{.LineNumber}}, {{.ColumnNumber}})" "{{.Filename}}"`,
		"vim": `vim -c "call cursor({{.LineNumber}}, {{.ColumnNumber}})" "{{.Filename}}"`,
	}
)

type AliasFile struct {
	filename    string
	fmtStr      string
	buf         bytes.Buffer
	writer      *bufio.Writer
	aliasPrefix string
	latestAlias uint
}

func getEditorCommand() string {
	warn := func(envVar string, value string) {
		fmt.Fprintf(os.Stderr, "WARNING: Unknown value %v for %s. To use %v, please configure it using %s instead.\n", envVar, value, envVar, EDITOR_CMD_FORMAT_ENV)
	}

	if val := os.Getenv(EDITOR_CMD_ENV); val != "" {
		if cmd, ok := EDITOR_MAP[filepath.Base(val)]; ok {
			return cmd
		} else {
			warn(EDITOR_CMD_ENV, val)
		}
	}

	if val := os.Getenv("EDITOR"); val != "" {
		if cmd, ok := EDITOR_MAP[filepath.Base(val)]; ok {
			return cmd
		} else {
			warn("EDITOR", val)
		}
	}

	return EDITOR_MAP[EDITOR_CMD_DEFAULT]
}

func NewAliasFile() (out *AliasFile) {
	var (
		aliasFilename string
		aliasCommand string
		aliasPrefix string
	)
	if val := os.Getenv(ALIAS_PREFIX_ENV); val != "" {
		aliasPrefix = val
	} else {
		aliasPrefix = ALIAS_PREFIX_DEFAULT
	}

	if val := os.Getenv(ALIAS_FILE_ENV); val != "" {
		aliasFilename = val
	} else {
		aliasFilename = ALIAS_FILE_DEFAULT
	}

	aliasCommand = getEditorCommand()

	out = &AliasFile{
		fmtStr:      fmt.Sprintf("alias %s{{.MatchIndex}}='%s'\n", aliasPrefix, aliasCommand),
		filename:    aliasFilename,
		aliasPrefix: aliasPrefix,
	}
	out.writer = bufio.NewWriter(&out.buf)

	return
}

type AliasTemplateVars struct {
	MatchIndex   uint
	Filename     string
	LineNumber   string
	ColumnNumber string
}

func (a *AliasFile) WriteAlias(index uint, filename string, linenum, colnum int64) {
	t := template.Must(template.New("alias").Parse(a.fmtStr))

	aliasVars := AliasTemplateVars{
		MatchIndex:   index,
		Filename:     filename,
		LineNumber:   fmt.Sprintf("%d", linenum),
		ColumnNumber: fmt.Sprintf("%d", colnum+1),
	}

	err := t.Execute(a.writer, aliasVars)
	if err != nil {
		log.Fatal(err)
	}
	a.latestAlias = index
}

func (a *AliasFile) WriteFile() {
	err = a.writer.Flush()
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(a.filename, a.buf.Bytes(), 0644)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(fmt.Sprintf("%s.count", a.filename), []byte(a.latestAlias), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
