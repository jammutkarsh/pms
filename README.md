# PMS

Project Management ~~System~~ Sucks!

A TUI utility to open your projects in your favorite editor.

## Commands

Since this is going to be a TUI. Some of the subcommands needs a different DX.

### Commands in active development

- `pms add [arg]` : adds the `[arg]` directory to your projects.
- `pms ls` or `pms`: lists all the projects you have added.

### Improve DX wrt TUI

- `pms rm [arg]` : removes the  `[arg]` directory.
- `pms go [arg]` : changes  to `[arg]` directory.
- `pms op [arg]` : opens  `[arg]` directory in default editor.
- `pms config -e "vscode"` : sets the default editor.

> `[arg]` would be a string directory. Like `.` *or* `~/code/go/cShare` *or* `/root/go/src/grec` *or* `$GOPATH/ugit`
