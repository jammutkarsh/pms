package charm

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

var defaultEditor string 

/*Custom Type*/
type Project struct {
	// rank int    
	name string 
	path string 
}

// FilterValue allows us to filter the options by name.
func (p Project) FilterValue() string { return p.name }

// Name returns the name of the project.
func (p Project) Name() string { return p.name }

// Path returns the path of the project.
func (p Project) Path() string { return p.path }

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

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	return m.list.View()
}
