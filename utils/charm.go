package utils

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/sahilm/fuzzy"
)

/*Custom Type*/
type Project struct {
	ProjectName string `json:"Name"`
	ProjectPath string `json:"Path"`
	Opens       int    `json:"opens,omitempty"`
	LastOpened  int64  `json:"lastOpened,omitempty"` // unix seconds
	Pinned      bool   `json:"pinned,omitempty"`
}

// Action is something to do with the selected project, triggered by any of its Keys.
// Run is a shell command executed inside the project directory; empty means just cd.
type Action struct {
	Name string   `json:"name"`
	Keys []string `json:"keys"`
	Run  string   `json:"run"`
}

var (
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true)
	dimStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	branchStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("110"))
	changedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	cleanStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("108"))
	missingStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("203"))
)

/*Main Model*/
type Model struct {
	input    textinput.Model
	projects []Project
	matches  []Project
	cursor   int
	height   int
	width    int
	actions  []Action
	help     bool
	git      map[string]gitStatus

	// Chosen and Action are set when the user triggers an action.
	Chosen Project
	Action *Action
}

// NewModel returns a new Model with some sane defaults.
func NewModel() *Model {
	in := textinput.New()
	in.Prompt = "❯ "
	in.Placeholder = "Search projects…"
	in.Focus()
	m := &Model{input: in, actions: ReadConfig().Actions, git: map[string]gitStatus{}}
	m.reload()
	return m
}

func (m *Model) Init() tea.Cmd {
	cmds := make([]tea.Cmd, len(m.projects))
	for i, p := range m.projects {
		cmds[i] = fetchGitStatus(p.ProjectPath)
	}
	return tea.Batch(cmds...)
}

// reload reads the projects from the config file and re-applies the search.
func (m *Model) reload() {
	m.projects = rank(ReadConfig().getProjects(), time.Now())
	m.filter()
}

// filter fuzzy-matches the search query against project paths.
func (m *Model) filter() {
	q := m.input.Value()
	if q == "" {
		m.matches = m.projects
	} else {
		paths := make([]string, len(m.projects))
		for i, p := range m.projects {
			paths[i] = p.ProjectPath
		}
		m.matches = nil
		for _, match := range fuzzy.Find(q, paths) {
			m.matches = append(m.matches, m.projects[match.Index])
		}
	}
	m.cursor = min(m.cursor, max(len(m.matches)-1, 0))
}

// Update updates the model on an event like a key press.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height, m.width = msg.Height, msg.Width
		// with no width, textinput only renders the placeholder's first rune
		m.input.SetWidth(max(msg.Width-lipgloss.Width(m.input.Prompt)-1, 1))
		return m, nil
	case gitStatusMsg:
		m.git[msg.path] = msg.status
		return m, nil
	case tea.KeyPressMsg:
		k := msg.String()
		if m.help {
			switch k {
			case "ctrl+c":
				return m, tea.Quit
			case "?", "esc", "q":
				m.help = false
			}
			return m, nil
		}
		switch k {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "?":
			m.help = true
			return m, nil
		case "up", "ctrl+p", "ctrl+k":
			m.cursor = max(m.cursor-1, 0)
			return m, nil
		case "down", "ctrl+n", "ctrl+j":
			m.cursor = min(m.cursor+1, max(len(m.matches)-1, 0))
			return m, nil
		case "ctrl+t":
			if len(m.matches) > 0 {
				if err := TogglePin(m.matches[m.cursor].ProjectPath); err != nil {
					return m, tea.Quit
				}
				m.reload()
			}
			return m, nil
		case "ctrl+d":
			if len(m.matches) > 0 {
				if err := DeleteProjet(m.matches[m.cursor].ProjectPath); err != nil {
					return m, tea.Quit
				}
				m.reload()
			}
			return m, nil
		}
		for i, a := range m.actions {
			if slices.Contains(a.Keys, k) {
				if len(m.matches) > 0 {
					m.Chosen, m.Action = m.matches[m.cursor], &m.actions[i]
					return m, tea.Quit
				}
				return m, nil
			}
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	m.filter()
	return m, cmd
}

// View renders the search box, the matching projects and the key hints, or the help screen.
func (m *Model) View() tea.View {
	var b strings.Builder
	if m.help {
		m.writeHelp(&b)
	} else {
		m.writeList(&b)
	}
	v := tea.NewView(b.String())
	v.AltScreen = true
	return v
}

func (m *Model) writeList(b *strings.Builder) {
	b.WriteString(m.input.View() + "\n\n")

	// each project takes a name line, a details line and a blank line
	rows := max((m.height-5)/3, 1)
	start := max(m.cursor-rows+1, 0)
	for i := start; i < min(start+rows, len(m.matches)); i++ {
		p := m.matches[i]
		name := p.ProjectName
		if p.Pinned {
			name = "★ " + name
		}
		// git details first so a long path is what gets truncated
		var details string
		if g, ok := m.git[p.ProjectPath]; ok {
			switch {
			case g.missing:
				details = missingStyle.Render("missing") + "  "
			case g.changed > 0:
				details = branchStyle.Render("⎇ "+g.branch) + "  " + changedStyle.Render(fmt.Sprintf("%d changed", g.changed)) + "  "
			default:
				details = branchStyle.Render("⎇ "+g.branch) + "  " + cleanStyle.Render("clean") + "  "
			}
		}
		details += dimStyle.Render(tildePath(p.ProjectPath))
		bar := "  "
		if i == m.cursor {
			bar, name = selectedStyle.Render("▌ "), selectedStyle.Render(name)
		}
		for _, line := range []string{bar + name, bar + details} {
			// a wrapped row would push the list past the screen height (#14)
			if m.width > 0 {
				line = ansi.Truncate(line, m.width, "…")
			}
			b.WriteString(line + "\n")
		}
		b.WriteString("\n")
	}
	if len(m.projects) == 0 {
		b.WriteString(dimStyle.Render("  No projects yet. Run `pms add` inside a project.") + "\n")
	}

	hints := make([]string, 0, len(m.actions)+1)
	for _, a := range m.actions {
		if len(a.Keys) > 0 {
			hints = append(hints, a.Keys[0]+" "+a.Name)
		}
	}
	hints = append(hints, "? help")
	// pin the hints to the last line of the screen
	pad := max(m.height-1-strings.Count(b.String(), "\n"), 1)
	b.WriteString(strings.Repeat("\n", pad) + dimStyle.Render(strings.Join(hints, " · ")))
}

// tildePath shortens paths under $HOME to ~/….
func tildePath(path string) string {
	home, err := os.UserHomeDir()
	if err == nil && strings.HasPrefix(path, home+string(os.PathSeparator)) {
		return "~" + path[len(home):]
	}
	return path
}

func (m *Model) writeHelp(b *strings.Builder) {
	row := func(keys, desc string) {
		b.WriteString("  " + selectedStyle.Width(26).Render(keys) + desc + "\n")
	}
	b.WriteString("Actions\n\n")
	for _, a := range m.actions {
		desc := a.Name + dimStyle.Render("  "+a.Run)
		if a.Run == "" {
			desc = a.Name + dimStyle.Render("  (change directory)")
		}
		row(strings.Join(a.Keys, " / "), desc)
	}
	b.WriteString("\nNavigation\n\n")
	row("↑ ↓ / ctrl+p ctrl+n", "move")
	row("ctrl+t", "pin / unpin project")
	row("ctrl+d", "remove project from list")
	row("esc / ctrl+c", "quit")
	b.WriteString("\n" + dimStyle.Render("Edit actions and keys in "+ConfigPath()+" · ? or esc to close"))
}
