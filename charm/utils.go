package charm

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sort"
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
func (c *config) getProjects() []Project {
	return c.Projects
}

// getDefaultEditor returns default editor
func (c config) getDefaultEditor() string {
	return c.DefaultEditor
}


func WriteConfig(path []byte) error {
	cf := readConfig()
	
	length := len(cf.Projects)
	name := strings.Split(string(path), "/")

	var newProject Project = Project{
		ProjectPath:           string(path),
		ProjectName:           strings.Split(string(path), "/")[len(name)-1],
		LastRecentlyUsedRank:  length + 1,
	}
	cf.Projects = append(cf.Projects, newProject)
	sort.SliceStable(cf.Projects, func(i, j int) bool {
		return cf.Projects[i].LastRecentlyUsedRank < cf.Projects[j].LastRecentlyUsedRank
	})
	byteData, err := json.MarshalIndent(cf, "", "    ")
	if err != nil {
		return err
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	
	file := home + "/" + configFile

	return ioutil.WriteFile(file, byteData, 0777)
}
