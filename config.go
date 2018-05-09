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
		PublicDirectory   string `yaml:"public_directory"`
		TemplateDirectory string `yaml:"template_directory"`
	}

	Database struct {
		Endpoint string
		Redis    struct {
			Endpoint      string
			TokenLifetime int `yaml:"token_lifetime"`
		}
	}

	Mailgun struct {
		Domain           string
		MailgunKey       string `yaml:"mailgun_key"`
		MailgunPublicKey string `yaml:"mailgun_public_key"`
		RootTemplate     string `yaml:"root_mail_template"`
		RootURL          string `yaml:"root_url"`
		Email            string `yaml:"email"`
	}

	Logger struct {
		Path     string
		FileName string
	}

	Gcs struct {
		Bucket    string `yaml:"bucket"`
		ProjectID string `yaml:"project_id"`
		PublicURL string `yaml:"public_url"`
	}

	Ocra struct {
		Endpoint string
		AppsKey  string `yaml:"apps_key"`
	}

	Voucher struct {
		Link string
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
