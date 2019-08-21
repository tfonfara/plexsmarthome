package app

import (
	"encoding/json"
	"log"
	"os"
	"github.com/tfonfara/plexsmarthome/helper"
)

type Configuration struct {
	Players  []string     `json:"players"`
	Hue      Hue          `json:"hue"`
	LaMetric LaMetric     `json:"laMetric"`
	Forwards []Forward `json:"forward"`
}

type Hue struct {
	IpAddress    string   `json:"ip"`
	ApiKey       string   `json:"apiKey"`
	MediaTypes   []string `json:"mediaTypes"`
	GroupId      string   `json:"groupId"`
	SceneIdPlay  string   `json:"sceneIdPlay"`
	SceneIdPause string   `json:"sceneIdPause"`
}

type LaMetric struct {
	IpAddress  string   `json:"ip"`
	ApiKey     string   `json:"apiKey"`
	MediaTypes []string `json:"mediaTypes"`
}

type Forward struct {
	Account     int    `json:"account"`
	Destination string `json:"destination"`
}

func NewConfiguration() *Configuration {
	configuration := Configuration{}
	configPath := helper.ConfigurationFilePath()

	file, err := os.Open(configPath);
	if err != nil {
		log.Fatalf("Configuration file not found: %v\n", err)
	}

	if err == nil {
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&configuration)
		if err != nil {
			log.Fatalf("Configuration file contains invalid json: %v\n", err)
		}
	}

	return &configuration
}
