package config

import (
	"os"
	"reflect"
	"regexp"

	"gopkg.in/yaml.v3"
)

type Source struct {
	GroupName string `yaml:"groupName"`
	Url       string
	Type      string
	AuthToken string `yaml:"authToken"`

	RawInclude []string           `yaml:"include"`
	Include    [](*regexp.Regexp) `yaml:"_include_internal"`

	RawExclude []string           `yaml:"exclude"`
	Exclude    [](*regexp.Regexp) `yaml:"_exclude_internal"`
}

const (
	SOURCE_TYPE_REPOSITORY = "repository"
	SOURCE_TYPE_GITHUB     = "github"
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
		sourceValue := reflect.ValueOf(source).Elem()

		for i := 0; i < sourceValue.NumField(); i++ {
			field := sourceValue.Field(i)
			if field.Kind() == reflect.String && field.String() != "" {
				field.SetString(os.ExpandEnv(field.String()))
			}
		}

		if source.Type == "" {
			source.Type = SOURCE_TYPE_REPOSITORY
		}

		source.Include = make([](*regexp.Regexp), 0, len(source.RawInclude))
		for _, regexRaw := range source.RawInclude {
			regex, err := regexp.Compile(regexRaw)
			if err != nil {
				return nil, err
			}
			source.Include = append(source.Include, regex)
		}

		source.Exclude = make([](*regexp.Regexp), 0, len(source.RawExclude))
		for _, regexRaw := range source.RawExclude {
			regex, err := regexp.Compile(regexRaw)
			if err != nil {
				return nil, err
			}
			source.Exclude = append(source.Exclude, regex)
		}
	}

	return result, nil
}
