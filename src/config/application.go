package config

import "github.com/kelseyhightower/envconfig"

type ApplicationConfig struct {
	SourcesPath    string `split_words:"true"`
	RepositoryPath string `split_words:"true"`
	Cron           string
}

func ReadApplicationConfig() (ApplicationConfig, error) {
	result := ApplicationConfig{}
	result.SourcesPath = "sources.yaml"
	result.RepositoryPath = "repositories"
	result.Cron = "@daily"

	err := envconfig.Process("GIRE", &result)
	return result, err
}
