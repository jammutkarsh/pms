package utils

import (
	"encoding/json"
	"errors"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

const configFile = ".pms.json"

// init() checks the presence of a valid config file.
// If it doesn't exist, it creates one with default values.
func init() {
	// Create config file if it doesn't exist
	if _, err := os.Stat(ConfigPath()); os.IsNotExist(err) {
		// Write default config
		if _, err := os.Create(ConfigPath()); err == nil {
			byteData, err := json.MarshalIndent(config{Projects: []Project{}, Actions: defaultActions("code")}, "", "    ")
			if err != nil {
				cobra.CheckErr(err)
			}
			os.WriteFile(ConfigPath(), byteData, 0o777)
		} else {
			cobra.CheckErr(err)
		}
	} else {
		cobra.CheckErr(err)
	}
}

type config struct {
	// DefaultEditor is only read to migrate old configs into Actions.
	DefaultEditor string    `json:"defaultEditor,omitempty"`
	Projects      []Project `json:"projects"`
	Actions       []Action  `json:"actions"`
}

// defaultActions seeds ~/.pms.json; everything else (agents etc.) is added there by the user.
// $EDITOR is left for the shell to expand at run time.
func defaultActions(ide string) []Action {
	return []Action{
		{Name: "Editor", Keys: []string{"enter", "ctrl+e"}, Run: "${EDITOR:-vi} ."},
		{Name: "IDE", Keys: []string{"ctrl+shift+e"}, Run: ide + " ."},
		{Name: "cd", Keys: []string{"alt+enter", "tab"}},
	}
}

// getProjects returns all projects
func (c config) getProjects() []Project {
	return c.Projects
}

// ReadConfig returns go struct of config file at ~/.pms.json
func ReadConfig() (c config) {
	configFile, err := os.Open(ConfigPath())
	if err != nil {
		cobra.CheckErr(err)
		return
	}
	defer configFile.Close()

	if err = json.NewDecoder(configFile).Decode(&c); err != nil {
		cobra.CheckErr(err)
		return
	}

	if len(c.Actions) == 0 {
		editor := c.DefaultEditor
		if editor == "" {
			editor = "code"
		}
		c.Actions, c.DefaultEditor = defaultActions(editor), ""
		// persist right away so the actions are there to edit
		cobra.CheckErr(writeConfig(c))
	}
	return c
}

// UpdateProjectList writes JSON representation of projects to ~/.pms.json
func UpdateProjectList(projectPath string) error {
	var (
		currentProjectIndex int
		currentProject      Project
	)
	projects := ReadConfig().getProjects()

	// finding the index of the selected project
	for i, p := range projects {
		if p.ProjectPath == projectPath {
			currentProjectIndex = i
			currentProject = p
			break
		}
	}

	// moving the selected project to the top of the list
	for i := currentProjectIndex; i > 0; i-- {
		projects[i] = projects[i-1]
	}
	projects[0] = currentProject

	c := ReadConfig()
	c.Projects = projects
	return writeConfig(c)
}

// DeleteProjet deletes JSON representation of project from ~/.pms.json
func DeleteProjet(projectPath string) error {
	var updatedProjects []Project
	currentProjects := ReadConfig().getProjects()

	for _, project := range currentProjects {
		if project.ProjectPath != projectPath {
			updatedProjects = append(updatedProjects, project)
		}
	}

	c := ReadConfig()
	c.Projects = updatedProjects
	return writeConfig(c)
}

// AddProject adds JSON representation of new project in ~/.pms.json
func AddProject(workingPath string) error {
	c := ReadConfig()
	directories := strings.Split(workingPath, string(os.PathSeparator))

	var newProject Project = Project{
		ProjectName: strings.Split(workingPath, string(os.PathSeparator))[len(directories)-1],
		ProjectPath: workingPath,
	}

	for _, p := range c.Projects {
		if p.ProjectPath == newProject.ProjectPath {
			return errors.New("project already exists")
		}
	}

	c.Projects = append(c.Projects, newProject)
	return writeConfig(c)
}

// writeConfig writes JSON representation of the config to ~/.pms.json
func writeConfig(c config) error {
	byteData, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(ConfigPath(), byteData, 0o777)
}

// ConfigPath returns the path of the config file
func ConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		cobra.CheckErr(err)
		return ""
	}

	return home + string(os.PathSeparator) + configFile
}
