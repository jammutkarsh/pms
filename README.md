# PMS - Project Management ~~System~~ Sucks

![Sample Image](./pms.gif)

**A TUI utility to open your projects in your favorite editor.**

## Installation

<!-- Install Go -->

1. Download Go from [here](https://go.dev/dl/).
2. Run the following command

```bash
go install github.com/JammUtkarsh/pms@latest
```

>Go installation is necessary, because I haven't set up a binary release which you can simply download and use.

## TUI

- [x] `pms`: Open TUI in terminal and list all projects
- [x] `/`: Search Project
- [x] `esc` or `q` : Quit TUI
- [x] `enter` or `spacebar`: open the project in default editor
- [x] `backspace` or `delete`: deletes the selected project
- [ ] `C`: change directory to the project (Not possible due)

## CLI

- [x] `pms add [arg]` : adds the `[arg]` directory to your projects.
- [ ] `pms ls`: lists all the projects you have added.
- [ ] `pms rm [arg]` : removes the  `[arg]` directory.
- [ ] `pms op [arg]` : opens  `[arg]` directory in default editor.
- [ ] `pms c [arg]` : changes  to `[arg]` directory.
- [ ] `pms config -e "vscode"` : sets the default editor.

> `[arg]` would be a string directory. Like `.` *or* `~/code/go/cShare` *or* `/root/go/src/grec` *or* `$GOPATH/ugit`
