package providers

import (
	"github.com/sibwaf/gire/src/config"
)

type Provider interface {
	// GetGroupName() string // todo?
	ListRepositories() ([]string, error)
}

func CreateProviderFor(source config.Source) Provider {
	switch source.Type {
	case config.SOURCE_TYPE_REPOSITORY:
		return RepositoryProvider{Url: source.Url}
	default:
		return nil
	}
}
