package charm

import (
	"os/exec"
	"sort"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

/*Custom Type*/
type Project struct {
	LastRecentlyUsedRank int    `json:"Rank"` // will be used for sorting projects (based on LRU algorithm from CPU ; )
	ProjectName          string `json:"Name"`
	ProjectPath          string `json:"Path"`
}

// FilterValue allows us to filter the options by name.
func (p Project) FilterValue() string { return p.ProjectName }

/*
Implements DefaultItem interface to render the list items.
https://pkg.go.dev/github.com/charmbracelet/bubbles@v0.15.0/list#DefaultItem
*/

// Title returns the name of the project.
func (p Project) Title() string { return p.ProjectName }

// Description returns the path of the project.
func (p Project) Description() string { return p.ProjectPath }

/*Main Model*/
type Model struct {
	list     list.Model
	quitting bool
}

// NewModel returns a new Model with some sane defaults.
func NewModel() *Model {
	return &Model{}
}

func (m *Model) OpenInEditor() error {
	selectedIProject := m.list.SelectedItem()
	project := selectedIProject.(Project)
	defaultEditor := getDefaultEditor()
	cmd := exec.Command(defaultEditor, project.ProjectPath)
	return cmd.Start()
}

func (m *Model) InitList(width, height int) {
	// setting the default values
	m.list = list.New([]list.Item{}, list.NewDefaultDelegate(), width, height)
	m.list.Title = "Projects"

	// reading the projects from file
	projects := getProjects()
	projectList := make([]list.Item, len(projects))
	sort.SliceStable(projects, func(i, j int) bool {
		return projects[i].LastRecentlyUsedRank < projects[j].LastRecentlyUsedRank
	})
	for i, project := range projects {
		projectList[i] = list.Item(project)
	}
	// setting list items
	m.list.SetItems(projectList)
}

// Move utils.go functionnallity here.
func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.InitList(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "enter", " ":
			if err := m.OpenInEditor(); err != nil {
				cobra.CheckErr(err)
			}
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View renders the list. 
func (m Model) View() string {
	if m.quitting {
		return ""
	}
	return m.list.View() + "\n"
}
