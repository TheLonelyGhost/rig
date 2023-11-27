package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

func highlightRange(text string, colStart, colStop int64) string {
	if text == "" {
		return ""
	}

	return text[0:colStart] + color.RedString(text[colStart:colStop]) + text[colStop:]
}

func aliasPrinter(aliasIndex uint, text string, lineNum int64, colStart int64, colStop int64) string {
	aliasPrefix := color.BlueString("[") + color.RedString("%d", aliasIndex) + color.BlueString("]")
	coordinates := fmt.Sprintf("%s:%s:", color.GreenString("%d", lineNum), color.CyanString("%d", colStart))
	content := strings.TrimSuffix(highlightRange(text, colStart, colStop), "\n")

	return fmt.Sprintf("%s %s%s", aliasPrefix, coordinates, content)
}

func filePathPrinter(filename string) string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	relativePath, err := filepath.Rel(cwd, filename)
	if err != nil {
		log.Fatal(err)
	}
	if filepath.IsLocal(relativePath) {
		relativePath = "." + string(filepath.Separator) + relativePath
	}

	return color.MagentaString(relativePath)
}
