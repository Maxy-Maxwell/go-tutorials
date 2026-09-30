package config

import (
	"encoding/json"
	"log"
	"os"
)

const fileName string = ".gatorconfig.json"

type Config struct {
	Db_url            string
	Current_user_name string
}

// Config related functions
func getConfigFilePath() string {
	confFilePath, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	return confFilePath + "/" + fileName
}

func Read() *Config {
	// Get the file path
	confFilePath := getConfigFilePath()

	// Actually read the file
	fileContents, err := os.ReadFile(confFilePath)
	conf := Config{}

	if err != nil {
		log.Fatal(err)
		return nil
	}

	// Parse it into the Config struct
	if err = json.Unmarshal(fileContents, &conf); err != nil {
		log.Fatal(err)
		return nil
	}

	return &conf
}

func (c *Config) SetUser(username string) error {
	c.Current_user_name = username

	err := writeConfigFile(c)
	if err != nil {
		return err
	}

	return nil
}

func writeConfigFile(c *Config) error {
	// Serialize Config struct
	fileContents, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	// Actually write it to the file
	confFilePath := getConfigFilePath()
	err = os.WriteFile(confFilePath, fileContents, 0644)
	if err != nil {
		return err
	}

	return nil
}
