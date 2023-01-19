package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Source struct {
	GroupName string
	Url       string
	Type      string
	Include   []string
	Exclude   []string
}

const (
	SOURCE_TYPE_REPOSITORY = "repository"
)

func ReadSourceList(path string) ([]*Source, error) {
	result := [](*Source){}

	configContent, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(configContent, &result)
	if err != nil {
		return nil, err
	}

	for _, source := range result {
		if source.GroupName == "" {
			source.GroupName = "_"
		}
		if source.Type == "" {
			source.Type = SOURCE_TYPE_REPOSITORY
		}
	}

	return result, nil
}
