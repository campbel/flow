package types

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Flowfile struct {
	Flows map[string]Flow `yaml:"flows"`
}

type Flow struct {
	Name  string            `yaml:"name"`
	Envs  map[string]string `yaml:"envs"`
	Vars  map[string]any    `yaml:"vars"`
	If    *Step             `yaml:"if"`
	Steps []Step            `yaml:"steps"`
}

type Step struct {
	Name  string            `yaml:"name"`
	Envs  map[string]string `yaml:"envs"`
	Vars  map[string]any    `yaml:"vars"`
	Shell string            `yaml:"sh"`

	// Controls
	If    *Step  `yaml:"if"`
	Go    *Step  `yaml:"go"`
	Defer *Step  `yaml:"defer"`
	Group []Step `yaml:"group"`
	Range []any  `yaml:"range"`
}

func LoadFlowfile(path string) (Flowfile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var flowConfig Flowfile
	if err := yaml.Unmarshal(data, &flowConfig); err != nil {
		panic(err)
	}

	return flowConfig, nil
}
