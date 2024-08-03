package models

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Device struct {
	Hostname string `yaml:"hostname" json:"hostname"`
	Host     string `yaml:"host" json:"host"`
	Platform string `yaml:"platform" json:"platform"`
	Protocol string `yaml:"protocol,omitempty" json:"protocol,omitempty"`
}

type Inventory struct {
	Devices []Device `yaml:"hosts" json:"hosts"`
}

func LoadInventory(filename string) (*Inventory, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	inventory := &Inventory{}
	switch {
	case filename[len(filename)-4:] == "yaml":
		err = yaml.Unmarshal(data, inventory)
	case filename[len(filename)-4:] == "json":
		err = json.Unmarshal(data, inventory)
	default:
		err = errors.New("unsupported file type")
	}

	if err != nil {
		return nil, err
	}

	return inventory, nil
}
