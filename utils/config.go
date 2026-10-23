package utils

import (
	"cmp"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

var setupOnce sync.Once

// setup makes sure a config file exists, moving an old ~/.pms.json over
// or writing defaults. It runs on the first read rather than in init() so
// importing the package (e.g. in tests) doesn't touch the user's config.
func setup() {
	path := ConfigPath()
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		cobra.CheckErr(err)
		return
	}
	cobra.CheckErr(os.MkdirAll(filepath.Dir(path), 0o700))
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	legacy := filepath.Join(home, ".pms.json")
	if data, err := os.ReadFile(legacy); err == nil {
		// older versions wrote the config world-writable, so write a fresh 0600 file rather than rename
		cobra.CheckErr(os.WriteFile(path, data, 0o600))
		cobra.CheckErr(os.Remove(legacy))
		return
	}
	cobra.CheckErr(writeConfig(config{Projects: []Project{}, Actions: defaultActions("code")}))
}

type config struct {
	// DefaultEditor is only read to migrate old configs into Actions.
	DefaultEditor string    `json:"defaultEditor,omitempty"`
	Projects      []Project `json:"projects"`
	Actions       []Action  `json:"actions"`
}

// defaultActions seeds the config file; everything else (agents etc.) is added there by the user.
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

// ReadConfig returns go struct of the config file
func ReadConfig() (c config) {
	setupOnce.Do(setup)
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

// RecordOpen bumps the open count and last-opened time used for ranking.
func RecordOpen(projectPath string) error {
	c := ReadConfig()
	for i := range c.Projects {
		if c.Projects[i].ProjectPath == projectPath {
			c.Projects[i].Opens++
			c.Projects[i].LastOpened = time.Now().Unix()
		}
	}
	return writeConfig(c)
}

// rank orders projects by frecency: open count weighted by how recently they
// were opened (the same buckets zoxide uses). Ties keep their list order.
func rank(projects []Project, now time.Time) []Project {
	score := func(p Project) float64 {
		switch age := now.Sub(time.Unix(p.LastOpened, 0)); {
		case age < time.Hour:
			return float64(p.Opens) * 4
		case age < 24*time.Hour:
			return float64(p.Opens) * 2
		case age < 7*24*time.Hour:
			return float64(p.Opens) / 2
		default:
			return float64(p.Opens) / 4
		}
	}
	slices.SortStableFunc(projects, func(a, b Project) int {
		return cmp.Compare(score(b), score(a))
	})
	return projects
}

// DeleteProjet deletes JSON representation of project from the config file
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

// AddProject adds JSON representation of new project in the config file
func AddProject(workingPath string) error {
	if info, err := os.Stat(workingPath); err != nil {
		return err
	} else if !info.IsDir() {
		return errors.New("not a directory")
	}

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

// writeConfig writes JSON representation of the config to the config file.
// The config holds shell commands pms runs, so it's private to the user (0600),
// and written via rename so a crash can't leave it half-written.
func writeConfig(c config) error {
	byteData, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		return err
	}
	path := ConfigPath()
	tmp, err := os.CreateTemp(filepath.Dir(path), ".pms-*.json")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.Write(byteData); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), path)
}

// ConfigPath returns the path of the config file: $XDG_CONFIG_HOME/pms/config.json,
// falling back to ~/.config like XDG says (os.UserConfigDir doesn't on macOS).
func ConfigPath() string {
	dir := os.Getenv("XDG_CONFIG_HOME")
	if dir == "" {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		dir = filepath.Join(home, ".config")
	}
	return filepath.Join(dir, "pms", "config.json")
}
