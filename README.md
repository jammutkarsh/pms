# PMS - Project Management ~~System~~ Sucks

![Sample Image](./pms.gif)

**A Raycast-style project launcher: fuzzy search your projects, then open one in your editor, IDE or AI agent, or `cd` into it.**

## Install

```bash
curl -sSfL https://goblin.run/github.com/JammUtkarsh/pms | sh
# or
go install github.com/JammUtkarsh/pms@latest
```

Then add this to `~/.zshrc` (or `~/.bashrc`) so `cd` works in your current shell. Without it, `cd` opens a subshell.

```bash
eval "$(pms init zsh)"
```

## Usage

```bash
pms                          # open the launcher
pms add [dir...]             # add directories (current one by default)
pms <git-url> [dir]          # clone and add: into ./, inside existing [dir], or as [dir]
pms user/repo                # same, from GitHub
pms gitlab.example.com/g/r   # same, from any other host over https
```

Each project shows its git branch, changed file count and path. Pinned projects stay on top; the rest are ordered by how often and how recently you open them.

| Key | Action |
| --- | --- |
| `enter` / `ctrl+e` | Open in `$EDITOR` (falls back to `vi`) |
| `ctrl+shift+e` | Open in IDE (`code .`) |
| `alt+enter` / `tab` | `cd` into the project |
| `↑` `↓` / `ctrl+p` `ctrl+n` | Move |
| `ctrl+t` | Pin / unpin |
| `ctrl+d` | Remove from list |
| `?` | Help |
| `esc` / `ctrl+c` | Quit |

## Config

`~/.config/pms/config.json` (respects `$XDG_CONFIG_HOME`; an old `~/.pms.json` is moved there). The first three keys above are actions you can rebind or extend, e.g. with AI agents:

```json
"actions": [
    { "name": "Editor", "keys": ["ctrl+e"], "run": "${EDITOR:-vi} ." },
    { "name": "IDE", "keys": ["ctrl+shift+e"], "run": "code ." },
    { "name": "cd", "keys": ["alt+enter", "tab"], "run": "" },
    { "name": "Claude Code", "keys": ["enter"], "run": "claude" },
    { "name": "Codex", "keys": ["shift+enter", "ctrl+x"], "run": "codex" }
]
```

`run` is a shell command run inside the project; empty `run` just `cd`s. If two actions share a key, the first wins. `shift+enter`, `ctrl+enter` and `ctrl+shift+…` need a terminal with the kitty keyboard protocol (Ghostty, kitty, WezTerm, iTerm2 with CSI u).
