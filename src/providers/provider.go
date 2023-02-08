package providers

import (
	"errors"
	"fmt"

	"github.com/sibwaf/gire/src/config"
)

type Provider interface {
	GetGroupName() string
	ListRepositories() ([]string, error)
}

func CreateProviderFor(source config.Source) (Provider, error) {
	switch source.Type {
	case config.SOURCE_TYPE_REPOSITORY:
		return RepositoryProvider{Url: source.Url}, nil
	case config.SOURCE_TYPE_GITHUB:
		return GithubProvider{Url: source.Url, AuthToken: source.AuthToken}, nil
	default:
		return nil, errors.New(fmt.Sprintf("No provider available for source type %s", source.Type))
	}
}
