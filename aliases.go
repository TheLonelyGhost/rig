package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
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
		"vim": `vim -c "call cursor({{.LineNumber}}, {{.ColumnNumber}})" "{{.Filename}}"`,
		"vscode": "",
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

	if val, ok := EDITOR_MAP[os.Getenv(EDITOR_CMD_ENV)]; ok {
		aliasCommand = val
	} else {
		aliasCommand = EDITOR_MAP[EDITOR_CMD_DEFAULT]
	}

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
	// write the function that cleans up after ourselves
	_, err := fmt.Fprintf(a.writer, "__rig_cleared=0; __rig-clear() { if [ \"${__rig_cleared}\" -gt 0 ]; then return 0; fi; for i in {1..%d}; do unalias \"%s${i}\"; done; __rig_cleared=1; }\n", a.latestAlias, a.aliasPrefix)
	if err != nil {
		log.Fatal(err)
	}

	err = a.writer.Flush()
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(a.filename, a.buf.Bytes(), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
