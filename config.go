package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

var config Config

type Config struct {
	Server struct {
		Host              string
		PublicDirectory   string `yaml:"public_directory"`
		TemplateDirectory string `yaml:"template_directory"`
	}

	Database struct {
		Endpoint string
		Redis    string
	}
}

func ReadConfig(f string, c *Config) error {
	log.Printf("Reading config file %q", f)

	d, err := ioutil.ReadFile(f)
	if err != nil {
		return fmt.Errorf("config: could not read config. (%v)", err)
	}

	if err := yaml.Unmarshal(d, c); err != nil {
		return fmt.Errorf("config: could not read config. (%v)", err)
	}

	return nil
}
