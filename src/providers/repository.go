package providers

type RepositoryProvider struct {
	Url string
}

func (p RepositoryProvider) GetGroupName() string {
	return ""
}

func (p RepositoryProvider) ListRepositories() ([]string, error) {
	return []string{p.Url}, nil
}
