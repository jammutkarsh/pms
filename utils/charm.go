package utils

import (
	"slices"
	"strings"

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
	m := &Model{input: in, actions: ReadConfig().Actions}
	m.reload()
	return m
}

func (m *Model) Init() tea.Cmd {
	return nil
}

// reload reads the projects from the config file and re-applies the search.
func (m *Model) reload() {
	m.projects = ReadConfig().getProjects()
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

	rows := max(m.height-5, 1)
	start := max(m.cursor-rows+1, 0)
	for i := start; i < min(start+rows, len(m.matches)); i++ {
		p := m.matches[i]
		var line string
		if i == m.cursor {
			line = selectedStyle.Render("▌ "+p.ProjectName) + "  " + dimStyle.Render(p.ProjectPath)
		} else {
			line = "  " + p.ProjectName + "  " + dimStyle.Render(p.ProjectPath)
		}
		// a wrapped row would push the list past the screen height (#14)
		if m.width > 0 {
			line = ansi.Truncate(line, m.width, "…")
		}
		b.WriteString(line + "\n")
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
	b.WriteString("\n" + dimStyle.Render(strings.Join(hints, " · ")))
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
	row("ctrl+d", "remove project from list")
	row("esc / ctrl+c", "quit")
	b.WriteString("\n" + dimStyle.Render("Edit actions and keys in "+ConfigPath()+" · ? or esc to close"))
}
