// Copyright (c) 2017 Christian Saide <Supernomad>
// Licensed under the MPL-2.0, for details see https://github.com/Supernomad/protond/blob/master/LICENSE

package common

import (
	"encoding/json"
	"io/ioutil"
	"path"

	"gopkg.in/yaml.v2"
)

// PluginConfig is a struct representing a input plugin configuration.
type PluginConfig struct {
	Type   string            `json:"type" yaml:"type"`
	Name   string            `json:"name" yaml:"name"`
	Config map[string]string `json:"config" yaml:"config"`
}

// ParsePluginConfigs parses a directory of files and returns the resulting array of configs.
func ParsePluginConfigs(dir string, log *Logger) ([]*PluginConfig, error) {
	configs := make([]*PluginConfig, 0)
	if PathExists(dir) {
		inputFiles, err := ioutil.ReadDir(dir)
		if err != nil {
			return nil, err
		}

		for i := 0; i < len(inputFiles); i++ {
			name := inputFiles[i].Name()
			ext := path.Ext(name)

			var cfg PluginConfig
			switch ext {
			case ".yml", ".yaml":
				fileData, err := ioutil.ReadFile(path.Join(dir, name))
				if err != nil {
					return nil, err
				}

				err = yaml.Unmarshal(fileData, &cfg)
				if err != nil {
					return nil, err
				}
			case ".json":
				fileData, err := ioutil.ReadFile(path.Join(dir, name))
				if err != nil {
					return nil, err
				}

				err = json.Unmarshal(fileData, &cfg)
				if err != nil {
					return nil, err
				}
			default:
				log.Warn.Printf("Input or Output configuration file '%s' is not one of the compatible configuration file types: 'json', 'yml', or 'yaml'.", name)
				continue
			}

			configs = append(configs, &cfg)
		}
	} else {
		log.Warn.Printf("The specified directory '%s' does not exist.", dir)
	}

	return configs, nil
}
