package main

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// Config main config file structure
type Config struct {
	SlackToken     string `yaml:"slack_token"`
	Port           int    `yaml:"port"`
	SlackWebookURL string `yaml:"slack_webhook_url"`
	Channel        string `yaml:"channel"`
	Database       string `yaml:"database"`
}

var config Config

func readConfig() {
	if len(os.Args) < 2 {
		log.Fatal("Missing required config parameter")
	}

	contents, err := ioutil.ReadFile(os.Args[1])

	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(contents, &config)

	if err != nil {
		log.Fatalf("Error parsing config: %v", err)
	}
}
