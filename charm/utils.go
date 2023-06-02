package charm

import (
	"encoding/json"
	"os"

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
func (c *config) getProjects() []Project {
	return c.Projects
}

// getDefaultEditor returns default editor
func (c config) getDefaultEditor() string {
	return c.DefaultEditor
}
