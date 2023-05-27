package charm

import (
	"encoding/json"
	"os"

	"github.com/spf13/cobra"
)

var configFile = "SampleConfig.json"

type config struct {
	DefaultEditor string `json:"defaultEditor"`
	Projects []Project `json:"projects"`
}


func readMetaData() (c config) {
	// home, err := os.UserHomeDir()
	home, err := os.Getwd() // for testing
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

func filterProjects() []Project{
	c:= readMetaData()
	return c.Projects
}

func getDefaultEditor() string {
	c := readMetaData()
	return c.DefaultEditor
}