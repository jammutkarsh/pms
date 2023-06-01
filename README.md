# PMS - Project Management ~~System~~ Sucks

![SampleImage](./pms.gif)

**A TUI utility to open your projects in your favorite editor.**

## TUI

- [x] `pms`: Open TUI in terminal and list all projects

- [x] `/`: Search Project
- [x] `esc` or `q` : Quit TUI
- [x] `enter` or `spacebar`: open the project in default editor
- [ ] `C`: change directory to the project (Not possible due)

## CLI

- [ ] `pms ls`: lists all the projects you have added.
- [ ] `pms add [arg]` : adds the `[arg]` directory to your projects.
- [ ] `pms rm [arg]` : removes the  `[arg]` directory.
- [ ] `pms op [arg]` : opens  `[arg]` directory in default editor.
- [ ] `pms c [arg]` : changes  to `[arg]` directory.
- [ ] `pms config -e "vscode"` : sets the default editor.

> `[arg]` would be a string directory. Like `.` *or* `~/code/go/cShare` *or* `/root/go/src/grec` *or* `$GOPATH/ugit`
