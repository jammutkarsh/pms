# PMS - Project Management ~~System~~ Sucks

![SampleImage](/Sample.png)

**A TUI utility to open your projects in your favorite editor.**

## TUI

### Development

- [x] `pms`: Open TUI in terminal and list all projects

### Global Bindings

- [x] `/`: Search Project
- [x] `esc` or `q` : Quit TUI

### Selected Project Key Bindings

- [ ] `O`: open the project in default editor
- [ ] `C`: change directory to the project

## CLI

### Development

- [ ] `pms ls`: lists all the projects you have added.
- [ ] `pms add [arg]` : adds the `[arg]` directory to your projects.
- [ ] `pms rm [arg]` : removes the  `[arg]` directory.
- [ ] `pms op [arg]` : opens  `[arg]` directory in default editor.
- [ ] `pms c [arg]` : changes  to `[arg]` directory.
- [ ] `pms config -e "vscode"` : sets the default editor.

> `[arg]` would be a string directory. Like `.` *or* `~/code/go/cShare` *or* `/root/go/src/grec` *or* `$GOPATH/ugit`
