package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/fatih/color"
)

const (
	BOOTSTRAP_ENV = "RIG_BOOTSTRAP"
	RIPGREP_CMD_ENV = "RIG_RIPGREP_CMD"
)

func optionIndex(args []string, option string) int {
	for i := len(args) - 1; i >= 0; i-- {
		if args[i] == option {
			return i
		}
	}
	return -1
}

func isatty(f *os.File) bool {
	stat, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}
	return stat.Mode()&os.ModeCharDevice != 0
}

func extractCmdExitCode(err error) int {
	if err != nil {
		// Extract real exit code
		// Source: https://stackoverflow.com/a/10385867
		if exitErr, ok := err.(*exec.ExitError); ok {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				return status.ExitStatus()
			}
		}
	}
	return 0
}

func passthru(cmd *exec.Cmd) int {
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	return extractCmdExitCode(err)
}

func generateAliases(cmd *exec.Cmd) int {
	cmd.Stderr = os.Stderr

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)

	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	var (
		outputLine  []byte
		aliasIndex  uint
		currentPath string
	)

	aliasFile := NewAliasFile()
	defer aliasFile.WriteFile()

	aliasIndex = 1
	for scanner.Scan() {
		outputLine = scanner.Bytes()
		obj := GrepResult{}
		err := json.Unmarshal(outputLine, &obj)
		if err != nil {
			log.Fatal(err)
		}
		if obj.Type != "match" {
			continue
		}
		data := obj.Data.(*GrepMatch)

		filePath, err := filepath.Abs(data.Path.Text)
		if err != nil {
			log.Fatal(err)
		}
		if filePath != currentPath {
			if currentPath != "" {
				fmt.Println("")
			}
			fmt.Println(filePathPrinter(filePath))
			currentPath = filePath
		}

		for _, submatch := range data.Submatches {
			lineNum, err := data.LineNumber.Int64()
			if err != nil {
				log.Fatal(err)
			}
			colStart, err := submatch.Start.Int64()
			if err != nil {
				log.Fatal(err)
			}
			colEnd, err := submatch.End.Int64()
			if err != nil {
				log.Fatal(err)
			}

			aliasFile.WriteAlias(aliasIndex, filePath, lineNum, colStart)
			fmt.Println(aliasPrinter(aliasIndex, data.Lines.Text, lineNum, colStart, colEnd))
			aliasIndex++
		}
	}

	err = cmd.Wait()
	return extractCmdExitCode(err)
}

func main() {
	userArgs := os.Args[1:]
	rigArgs := []string{"--json", "--no-stats"}

	if shell, ok := os.LookupEnv(BOOTSTRAP_ENV); ok {
		bootstrapper, err := NewBootstrapper(shell)
		if err != nil {
			log.Fatal(err)
		}
		err = bootstrapper.DoBootstrap()
		if err != nil {
			log.Fatal(err)
			panic(err)
		}
		os.Exit(0)
	}

	ripgrepCmd := os.Getenv(RIPGREP_CMD_ENV)
	if ripgrepCmd == "" {
		ripgrepCmd = "rg"
	}

	if len(userArgs) == 0 {
		// Display help message
		os.Exit(passthru(exec.Command(ripgrepCmd, userArgs...)))
	}

	for _, opt := range []string{"--version", "--help", "--files", "--stats", "--type-list"} {
		if idx := optionIndex(userArgs, opt); idx >= 0 {
			if opt == "--help" || opt == "--version" {
				fmt.Fprintln(os.Stdout, VersionString())
			}
			os.Exit(passthru(exec.Command(ripgrepCmd, userArgs...)))
		}
	}

	if !isatty(os.Stdin) || !isatty(os.Stdout) {
		os.Exit(passthru(exec.Command(ripgrepCmd, userArgs...)))
	}

	// Handle auto-coloring
	if idx := optionIndex(userArgs, "--color"); idx >= 0 && userArgs[idx+1] == "never" {
		color.NoColor = true
	} else if idx := optionIndex(userArgs, "--color=never"); idx >= 0 {
		color.NoColor = true
	}

	args := append(rigArgs, userArgs...)
	os.Exit(generateAliases(exec.Command(ripgrepCmd, args...)))
}
