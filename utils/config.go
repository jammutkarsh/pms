package utils

import (
	"encoding/json"
	"io/ioutil"
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
			byteData, err := json.MarshalIndent(config{"code", []Project{}}, "", "    ")
			if err != nil {
				cobra.CheckErr(err)
			}
			ioutil.WriteFile(ConfigPath(), byteData, 0777)
		} else {
			cobra.CheckErr(err)
		}
	} else {
		cobra.CheckErr(err)
	}
}

type config struct {
	DefaultEditor string    `json:"defaultEditor"`
	Projects      []Project `json:"projects"`
}

// getProjects returns all projects
func (c config) getProjects() []Project {
	return c.Projects
}

// getDefaultEditor returns default editor
func (c config) getDefaultEditor() string {
	return c.DefaultEditor
}

// readConfig returns go struct of config file at ~/.pms.json
func readConfig() (c config) {
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
	return c
}

// AddProject adds JSON representation of new project in ~/.pms.json
func AddProject(workingPath string) error {
	config := readConfig()

	directories := strings.Split(workingPath, "/")

	var newProject Project = Project{
		ProjectName: strings.Split(workingPath, string(os.PathSeparator))[len(directories)-1],
		ProjectPath: workingPath,
	}

	config.Projects = append(config.Projects, newProject)
	byteData, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(ConfigPath(), byteData, 0777)
}

func ConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		cobra.CheckErr(err)
		return ""
	}
	return home + "/" + configFile
}
