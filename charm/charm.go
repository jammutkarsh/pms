package charm

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

var defaultEditor string 

/*Custom Type*/
type Project struct {
	LastRecentlyUsedRank int `json:"Rank"` // will be used for sorting projects (based on LRU algorithm from CPU ; )
	ProjectName string `json:"Name"`
	ProjectPath string `json:"Path"`
}

// FilterValue allows us to filter the options by name.
func (p Project) FilterValue() string { return p.ProjectName }

// Name returns the name of the project.
func (p Project) Name() string { return p.ProjectName }

// Path returns the path of the project.
func (p Project) Path() string { return p.ProjectPath }

/*Main Model*/
type Model struct {
	list list.Model
}

// NewModel returns a new Model with some sane defaults.
func NewModel() *Model {
	return &Model{}
}

func (m *Model) InitList(width, height int) {
	// setting the default values
	m.list = list.New([]list.Item{}, list.NewDefaultDelegate(), width, height)
	m.list.Title = "Projects"

	// reading the projects from file
	projects := filterProjects()
	projectList := make([]list.Item, len(projects))
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
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.list.View()
}
