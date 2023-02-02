package providers

import (
	"errors"
	"fmt"

	"github.com/sibwaf/gire/src/config"
)

type Provider interface {
	// GetGroupName() string // todo?
	ListRepositories() ([]string, error)
}

func CreateProviderFor(source config.Source) (Provider, error) {
	switch source.Type {
	case config.SOURCE_TYPE_REPOSITORY:
		return RepositoryProvider{Url: source.Url}, nil
	default:
		return nil, errors.New(fmt.Sprintf("No provider available for source type %s", source.Type))
	}
}
