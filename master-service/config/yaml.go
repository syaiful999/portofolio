package config

import (
	yaml "github.com/asim/go-micro/plugins/config/encoder/yaml/v4"
	"go-micro.dev/v4/config"
	"go-micro.dev/v4/config/reader"
	"go-micro.dev/v4/config/reader/json"
	"go-micro.dev/v4/config/source/file"
)

func ReadYamlConfig(path string, scan interface{}, pathTo ...string) error {
	enc := yaml.NewEncoder()

	// new config
	c, err := config.NewConfig(
		config.WithReader(
			json.NewReader(
				reader.WithEncoder(enc),
			),
		),
	)
	if err != nil {
		return err
	}

	if err := c.Load(
		file.NewSource(
			file.WithPath(path),
		),
	); err != nil {
		return err
	}

	if err := c.Get(pathTo...).Scan(scan); err != nil {
		return err
	}

	return nil
}
