package providers

type RepositoryProvider struct {
	Url string
}

func (p RepositoryProvider) ListRepositories() ([]string, error) {
	return []string{p.Url}, nil
}
