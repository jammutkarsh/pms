# PMS - Project Management ~~System~~ Sucks

![Sample Image](./pms.gif)

**A Raycast-style launcher for your projects: search, then open it in your editor, IDE, any AI agent you configure, or `cd` into it.**

## Installation

### Using Curl

```
 curl -sSfL https://goblin.run/github.com/JammUtkarsh/pms | sh
```

### With Go Installed 

```bash
go install github.com/JammUtkarsh/pms@latest
```

## Shell integration

A program can't change its parent shell's directory, so add this to `~/.zshrc` (or `~/.bashrc`):

```bash
eval "$(pms init zsh)"
```

Now `cd` lands in your current shell, and agents run in your shell inside the project.
Without it, `cd` opens a subshell in the project instead.

## TUI

Just start typing to fuzzy search. Like Raycast, each kind of enter does something different.

Default actions (all rebindable):

| Key | Action |
| --- | --- |
| `enter` / `ctrl+e` | `$EDITOR` (falls back to `vi`) |
| `ctrl+shift+e` | IDE (`code .`) |
| `alt+enter` / `tab` | `cd` into the project |

Fixed keys:

| Key | Action |
| --- | --- |
| `↑` `↓` / `ctrl+p` `ctrl+n` | Move |
| `?` | Help: every action, its keys and the config path |
| `ctrl+d` | Remove project from the list |
| `esc` / `ctrl+c` | Quit |

`shift+enter`, `ctrl+enter` and `ctrl+shift+…` need a terminal with the kitty keyboard protocol (Ghostty, kitty, WezTerm, iTerm2 with CSI u). Elsewhere bind a plain `ctrl+…` key too.

### Actions

Actions live in `~/.pms.json` (defaults are written on first run). `run` is a shell command run inside the project, so `$VARS` expand; an empty `run` just `cd`s. Add your AI agents here:

```json
"actions": [
    { "name": "Editor", "keys": ["ctrl+e"], "run": "${EDITOR:-vi} ." },
    { "name": "IDE", "keys": ["ctrl+shift+e"], "run": "code ." },
    { "name": "cd", "keys": ["alt+enter", "tab"], "run": "" },
    { "name": "Claude Code", "keys": ["enter"], "run": "claude" },
    { "name": "Codex", "keys": ["shift+enter", "ctrl+x"], "run": "codex" },
    { "name": "OpenCode", "keys": ["ctrl+enter", "ctrl+o"], "run": "opencode" }
]
```

If two actions share a key, the first one wins.

## CLI

- [x] `pms add [dir...]` : adds one or more directories to your projects.
<!--
- [ ] `pms ls`: lists all the projects you have added.
- [ ] `pms rm [arg]` : removes the  `[arg]` directory.
- [ ] `pms op [arg]` : opens  `[arg]` directory in default editor.
- [ ] `pms c [arg]` : changes  to `[arg]` directory.
- [ ] `pms config -e "vscode"` : sets the default editor.
-->

> `[arg]` would be a string directory. Like `.` *or* `~/code/go/cShare` *or* `/root/go/src/grec` *or* `$GOPATH/ugit`
