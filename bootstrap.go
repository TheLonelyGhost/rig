package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
)

type BootstrapperVars struct {
	AliasFileEnv       string
	AliasFileDefault   string
	AliasPrefixEnv     string
	AliasPrefixDefault string
}

const (
	BEGIN_MARKER = `

#================================#
# === BEGIN RIG BOOTSTRAPPER === #
`
	BASH_BOOTSTRAPPER = `
function rig() {
  __rig-clear
  command rig "$@"
  source "{{ "${" }}{{ .AliasFileEnv }}:-{{ .AliasFileDefault }}{{ "}" }}"
}
function __rig-clear() {
  local alias_count
  if [ -e "{{ "${" }}{{ .AliasFileEnv }}{{ "}" }}.count" ]; then
	alias_count="$(cat "{{ "${" }}{{ .AliasFileEnv }}{{ "}" }}.count")"
  fi
  if [ "$alias_count" -eq 0 ]; then return 0; fi

  for i in $(seq $alias_count); do
    unalias "{{ "${" }}{{ .AliasPrefixEnv }}{{ "}" }}${i}"
  done

  echo 0 > "{{ "${" }}{{ .AliasFileEnv }}{{ "}" }}.count"
}
: "{{ "${" }}{{ .AliasPrefixEnv }}:={{ .AliasPrefixDefault }}{{ "}" }}"
: "{{ "${" }}{{ .AliasFileEnv }}:={{ .AliasFileDefault }}{{ "}" }}"
`
	ZSH_BOOTSTRAPPER = BASH_BOOTSTRAPPER
	FISH_BOOTSTRAPPER = `
function rig
  __rig-clear
  command rig $argv

  source "${{ .AliasFileEnv }}"
end
function __rig-clear
  if not test -e "${{ .AliasFileEnv }}.count"
    return 0
  end

  set -f alias_count (cat "${{ .AliasFileEnv }}.count")

  if test -z $alias_count
    return 1
  else if test $alias_count -eq 0
    return 0
  end

  for i in (seq $alias_count)
    functions -e "${{ .AliasPrefixEnv }}$i"
  end

  echo 0 > "${{ .AliasFileEnv }}.count"
end
if test -z "${{ .AliasPrefixEnv }}"
  set -u {{ .AliasPrefixEnv }} '{{ .AliasPrefixDefault }}'
end
if test -z "${{ .AliasFileEnv }}"
  set -u {{ .AliasFileEnv }} '{{ .AliasFileDefault }}'
end
`
	END_MARKER = `
# ==== END RIG BOOTSTRAPPER ==== #
#================================#

`
)

type ShellBootstrapper interface {
	GetRcFile() string
	IsBootstrapped() bool
	DoBootstrap() error
}

type BashBootstrapper struct {
	rcFile string
}

func (b *BashBootstrapper) GetRcFile() string {
	if b.rcFile == "" {
		if runtime.GOOS == "darwin" {
			b.rcFile = os.ExpandEnv("${HOME}/.bash_profile")
		} else {
			b.rcFile = os.ExpandEnv("${HOME}/.bashrc")
		}
	}
	return b.rcFile
}
func (b *BashBootstrapper) IsBootstrapped() bool {
	_, err := os.Stat(b.GetRcFile())
	if err != nil {
		return false
	}
	handle, err := os.OpenFile(b.GetRcFile(), os.O_RDONLY, os.ModePerm)
	defer func() { _ = handle.Close() }()
	if err != nil {
		return false
	}
	scanner := bufio.NewScanner(handle)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), " BEGIN RIG BOOTSTRAPPER ") {
			return true
		}
	}

	return false
}
func (b *BashBootstrapper) DoBootstrap() (err error) {
	if b.IsBootstrapped() {
		log.Printf("%s is already bootstrapped\n", b.GetRcFile())
		return
	}
	buffer := bytes.Buffer{}
	buf := bufio.NewWriter(&buffer)

	_, err = fmt.Fprint(buf, BEGIN_MARKER)
	if err != nil {
		return
	}

	aliasVars := BootstrapperVars{
		AliasFileEnv:       ALIAS_FILE_ENV,
		AliasFileDefault:   ALIAS_FILE_DEFAULT,
		AliasPrefixEnv:     ALIAS_PREFIX_ENV,
		AliasPrefixDefault: ALIAS_PREFIX_DEFAULT,
	}
	t := template.Must(template.New("bootstrapper").Parse(BASH_BOOTSTRAPPER))
	err = t.Execute(buf, aliasVars)
	if err != nil {
		log.Fatal(err)
		return
	}

	_, err = fmt.Fprint(buf, END_MARKER)
	if err != nil {
		return
	}

	err = buf.Flush()
	if err != nil {
		log.Fatal(err)
		return
	}

	err = os.MkdirAll(filepath.Dir(b.GetRcFile()), 0755)
	if err != nil {
		return
	}

	handle, err := os.OpenFile(b.GetRcFile(), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	defer func() {
		err = handle.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	if err != nil {
		return
	}

	_, err = handle.Write(buffer.Bytes())

	log.Printf("Finished bootstrapping %s with shell hooks. Please restart your terminal session.\n", b.GetRcFile())
	return
}

type FishBootstrapper struct {
	rcFile string
}

func (f *FishBootstrapper) GetRcFile() string {
	if f.rcFile == "" {
		f.rcFile = os.ExpandEnv("${HOME}/.config/fish/conf.d/rig.fish")
	}
	return f.rcFile
}
func (f *FishBootstrapper) IsBootstrapped() bool {
	_, err := os.Stat(f.GetRcFile())
	if err != nil {
		return false
	}
	handle, err := os.OpenFile(f.GetRcFile(), os.O_RDONLY, os.ModePerm)
	defer func() { _ = handle.Close() }()
	if err != nil {
		return false
	}
	scanner := bufio.NewScanner(handle)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), " BEGIN RIG BOOTSTRAPPER ") {
			return true
		}
	}

	return false
}
func (f *FishBootstrapper) DoBootstrap() (err error) {
	if f.IsBootstrapped() {
		log.Printf("%s is already bootstrapped\n", f.GetRcFile())
		return
	}
	buffer := bytes.Buffer{}
	buf := bufio.NewWriter(&buffer)

	_, err = fmt.Fprint(buf, BEGIN_MARKER)
	if err != nil {
		return
	}

	aliasVars := BootstrapperVars{
		AliasFileEnv:       ALIAS_FILE_ENV,
		AliasFileDefault:   ALIAS_FILE_DEFAULT,
		AliasPrefixEnv:     ALIAS_PREFIX_ENV,
		AliasPrefixDefault: ALIAS_PREFIX_DEFAULT,
	}
	t := template.Must(template.New("bootstrapper").Parse(FISH_BOOTSTRAPPER))
	err = t.Execute(buf, aliasVars)
	if err != nil {
		log.Fatal(err)
		return
	}

	_, err = fmt.Fprint(buf, END_MARKER)
	if err != nil {
		return
	}

	err = buf.Flush()
	if err != nil {
		log.Fatal(err)
		return
	}

	err = os.MkdirAll(filepath.Dir(f.GetRcFile()), 0755)
	if err != nil {
		return
	}

	handle, err := os.OpenFile(f.GetRcFile(), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	defer func() {
		err = handle.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	if err != nil {
		return
	}

	_, err = handle.Write(buffer.Bytes())

	log.Printf("Finished bootstrapping %s with shell hooks. Please restart your terminal session.\n", f.GetRcFile())
	return
}

type ZshBootstrapper struct {
	rcFile string
}

func (z *ZshBootstrapper) GetRcFile() string {
	if z.rcFile == "" {
		z.rcFile = os.ExpandEnv("${HOME}/.zshrc")
	}
	return z.rcFile
}
func (z *ZshBootstrapper) IsBootstrapped() bool {
	_, err := os.Stat(z.GetRcFile())
	if err != nil {
		return false
	}
	handle, err := os.OpenFile(z.GetRcFile(), os.O_RDONLY, os.ModePerm)
	defer func() { _ = handle.Close() }()
	if err != nil {
		return false
	}
	scanner := bufio.NewScanner(handle)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), " BEGIN RIG BOOTSTRAPPER ") {
			return true
		}
	}

	return false
}
func (z *ZshBootstrapper) DoBootstrap() (err error) {
	if z.IsBootstrapped() {
		log.Printf("%s is already bootstrapped\n", z.GetRcFile())
		return
	}
	buffer := bytes.Buffer{}
	buf := bufio.NewWriter(&buffer)

	_, err = fmt.Fprint(buf, BEGIN_MARKER)
	if err != nil {
		return
	}

	aliasVars := BootstrapperVars{
		AliasFileEnv:       ALIAS_FILE_ENV,
		AliasFileDefault:   ALIAS_FILE_DEFAULT,
		AliasPrefixEnv:     ALIAS_PREFIX_ENV,
		AliasPrefixDefault: ALIAS_PREFIX_DEFAULT,
	}
	t := template.Must(template.New("bootstrapper").Parse(ZSH_BOOTSTRAPPER))
	err = t.Execute(buf, aliasVars)
	if err != nil {
		log.Fatal(err)
		return
	}

	_, err = fmt.Fprint(buf, END_MARKER)
	if err != nil {
		return
	}

	err = buf.Flush()
	if err != nil {
		log.Fatal(err)
		return
	}

	err = os.MkdirAll(filepath.Dir(z.GetRcFile()), 0755)
	if err != nil {
		log.Fatal(err)
		return
	}

	handle, err := os.OpenFile(z.GetRcFile(), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	defer func() {
		err = handle.Close()
		if err != nil {
			panic(err)
			// log.Fatal(err)
		}
	}()
	if err != nil {
		panic(err)
		// return
	}

	_, err = handle.Write(buffer.Bytes())

	if err != nil {
		panic(err)
	}

	log.Printf("Finished bootstrapping %s with shell hooks. Please restart your terminal session.\n", z.GetRcFile())
	return
}

func NewBootstrapper(shell string) (out ShellBootstrapper, err error) {
	switch filepath.Base(shell) {
	case "bash":
		log.Println("Detected bash bootstrap")
		out = &BashBootstrapper{}

	case "fish":
		log.Println("Detected fish bootstrap")
		out = &FishBootstrapper{}

	case "zsh":
		log.Println("Detected zsh bootstrap")
		out = &ZshBootstrapper{}

	default:
		err = fmt.Errorf("unsupported shell: %s", shell)
	}
	return
}
