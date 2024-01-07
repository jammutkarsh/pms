# PMS - Project Management ~~System~~ Sucks

![Sample Image](./pms.gif)

**A TUI utility to open your projects in your favorite editor.**

## Installation

### Using Curl

```
 curl -sSfL https://goblin.run/github.com/JammUtkarsh/pms | sh
```

### With Go Installed 

```bash
go install github.com/JammUtkarsh/pms@latest
```

## TUI

- [x] `pms`: Open TUI
- [x] `/`: Search project
- [x] `enter` : Open the project in default editor
- [x] `#` : Delete a project
- [x] `ctrl + c` : Quit

## CLI

- [x] `pms add [arg]` : adds `[arg]` directory to your projects.
<!--
- [ ] `pms ls`: lists all the projects you have added.
- [ ] `pms rm [arg]` : removes the  `[arg]` directory.
- [ ] `pms op [arg]` : opens  `[arg]` directory in default editor.
- [ ] `pms c [arg]` : changes  to `[arg]` directory.
- [ ] `pms config -e "vscode"` : sets the default editor.
-->

> `[arg]` would be a string directory. Like `.` *or* `~/code/go/cShare` *or* `/root/go/src/grec` *or* `$GOPATH/ugit`
