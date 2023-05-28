package charm

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

const configFile = ".pms.json"

type config struct {
	DefaultEditor string    `json:"defaultEditor"`
	Projects      []Project `json:"projects"`
}

// readConfig returns go struct of config file at ~/.pms.json
func readConfig() (c config) {
	home, err := os.UserHomeDir()
	if err != nil {
		cobra.CheckErr(err)
		return
	}

	file, err := os.Open(home + "/" + configFile)
	if err != nil {
		cobra.CheckErr(err)
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&c); err != nil {
		cobra.CheckErr(err)
		return
	}
	return c
}

// getProjects returns all projects
func getProjects() []Project {
	c := readConfig()
	hideHomeDir(&c)
	return c.Projects
}

// getDefaultEditor returns default editor
func getDefaultEditor() string {
	c := readConfig()
	return c.DefaultEditor
}

func hideHomeDir(c *config) {
	home, err := os.UserHomeDir()
	if err != nil {
		cobra.CheckErr(err)
	}
	for i, project := range c.Projects {
		c.Projects[i].ProjectPath = strings.Replace(project.ProjectPath, home, "~", 1)
	}
}