# Rig

Rigging up [`ripgrep`](https://github.com/burntsushi/ripgrep) to open results in an editor. Unlike regular `ripgrep`, `rig` hooks into your shell to provide shell aliases that open up your editor of choice.

Inspired by [`tag`](https://github.com/aykamko/tag).

## Usage

```bash
~/workspace $ rig 'my pattern here' ./path/to/search/under
```

## Configuration

Configuration is done entirely with environment variables:

| Variable           | Default            | Description                                                      |
|:-------------------|:-------------------|:-----------------------------------------------------------------|
| `RIG_BOOTSTRAP`    | _N/a_              | Bootstraps your shell init script with the required shell hooks. |
| `RIG_ALIAS_FILE`   | `/tmp/rig-aliases` | A temp file containing shell aliases for each match.             |
| `RIG_ALIAS_PREFIX` | `e`                | The prefix for each shell alias. Keep this to 1 or 2 letters.    |
| `RIG_EDITOR`       | `vim`              | Which editor you want the shell aliases to open. See below.      |
| `RIG_EDIT_COMMAND` | _N/a_              | Command used to open your editor to the exact line and column.   |
| `RIG_RIPGREP_CMD`  | `rg`               | Fully-qualified path to the ripgrep executable for `rig` to use. |

These environment variables may be placed in your shell init script (e.g., `~/.bashrc`).

### Editors

The following are a list of supported editors for which `RIG_EDITOR` can specify the command name. To support any other editors, the command must be specified with `RIG_EDIT_COMMAND`. (e.g., `mcedit` for Midnight Commander) to choose your editor

| Editor               | `RIG_EDITOR` | `RIG_EDIT_COMMAND` |
|:---------------------|:-------------|:-------------------|
| [Emacs]              | `emacs`      | `emacs +{{ .LineNumber }}:{{ .ColumnNumber }} --file="{{ .Filename }}"` |
| [Helix]              | `hx`         | `hx "{{ .Filename }}:{{ .LineNumber }}:{{ .ColumnNumber }}"` |
| [Micro]              | `micro`      | `micro +{{ .LineNumber }}:{{ .ColumnNumber }} "{{ .Filename }}"` |
| [Midnight Commander] | `mcedit`     | `mcedit "{{ .Filename }}:{{ .LineNumber }}"` |
| [Nano]               | `nano`       | `nano +{{.LineNumber}},{{ .ColumnNumber }} "{{ .Filename }}"` |
| [NeoVim]             | `nvim`       | `nvim -c "call cursor({{.LineNumber}}, {{.ColumnNumber}})" "{{.Filename}}"` |
| [Nice Editor]        | `ne`         | `ne +{{.LineNumber}},{{ .ColumnNumber }} "{{ .Filename }}"` |
| [Vim]                | `vim`        | `vim -c "call cursor({{.LineNumber}}, {{.ColumnNumber}})" "{{.Filename}}"` |
| [Visual Studio Code] | `code`       | `code --goto "{{ .Filename }}:{{ .LineNumber }}:{{ .ColumnNumber }}"` |

[Emacs]: https://www.gnu.org/software/emacs/
[Helix]: https://helix-editor.com/
[Micro]: https://micro-editor.github.io/
[Midnight Commander]: https://midnight-commander.org/
[Nano]: https://www.nano-editor.org/
[NeoVim]: https://neovim.io/
[Nice Editor]: https://ne.di.unimi.it/
[Vim]: https://www.vim.org/
[Visual Studio Code]: https://code.visualstudio.com/

#### Custom editor support

The `RIG_EDIT_COMMAND` variable is a rendered go template that exposes the following variables:

| Variable              | Example               | Description                                              |
|:----------------------|:----------------------|:---------------------------------------------------------|
| `{{ .MatchIndex }}`   | `23`                  | The number appended to the shell alias                   |
| `{{ .Filename }}`     | `/tmp/fizz/buzz.html` | The absolute path to the file containing the match       |
| `{{ .LineNumber }}`   | `590`                 | The line on which the match occurs                       |
| `{{ .ColumnNumber }}` | `14`                  | The column on which the first matching character appears |

### Shell environments

- Bourne Again Shell (`bash`)
- Zsh (`zsh`)
- Fish (`fish`)

## Install

Compile it yourself using Nix, or download the appropriate binary from the [latest release](https://github.com/TheLonelyGhost/rig/releases/latest).

The resulting binary must be named `rig` and be present somewhere in your PATH (e.g., `/usr/local/bin/rig`).

On first run, execute as follows:

```bash
~/workspace $ env RIG_BOOTSTRAP="$SHELL" rig
```

Then restart your shell session.

**OPTIONAL:** Run `alias rg=rig` to use `rig` as a drop-in replacement for most of the common uses of ripgrep.
