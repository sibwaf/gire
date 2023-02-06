package providers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/sibwaf/gire/src/util"
)

const (
	GITHUB_API_HOST  = "https://api.github.com"
	GITHUB_PAGE_SIZE = 100
)

type GithubProvider struct {
	Url       string
	AuthToken string
}

type githubRepository struct {
	HtmlUrl string `json:"html_url"`
	SshUrl  string `json:"ssh_url"`
}

func (p GithubProvider) ListRepositories() ([]string, error) {
	name := extractName(p.Url)
	if name == "" {
		return nil, errors.New(fmt.Sprintf("Failed to extract user/organization name from url %s", p.Url))
	}

	var wg sync.WaitGroup
	repositories := make(chan githubRepository)

	wg.Add(1)
	go callGithubApi(fmt.Sprintf("%s/users/%s/repos", GITHUB_API_HOST, name), p.AuthToken, repositories, &wg)

	if p.AuthToken != "" {
		wg.Add(1)
		go callGithubApi(fmt.Sprintf("%s/user/repos", GITHUB_API_HOST), p.AuthToken, repositories, &wg)
	}

	go func() {
		wg.Wait()
		close(repositories)
	}()

	sourceUrlLower := strings.ToLower(p.Url)

	result := util.NewSet[string]()
	for repository := range repositories {
		if !strings.HasPrefix(strings.ToLower(repository.HtmlUrl), sourceUrlLower) {
			continue
		}

		result.Add(repository.SshUrl)
	}

	return result.ToSlice(), nil
}

func extractName(url string) string {
	parts := strings.Split(url, "/")
	for i := len(parts) - 1; i >= 0; i-- {
		if parts[i] != "" {
			return parts[i]
		}
	}

	return ""
}

// todo: error handling
// todo: ratelimit handling?
func callGithubApi(url string, token string, repositories chan githubRepository, wg *sync.WaitGroup) {
	defer wg.Done()

	page := 1

	for true {
		queryParams := fmt.Sprintf("page=%d&per_page=%d", page, GITHUB_PAGE_SIZE)

		req, err := http.NewRequest("GET", fmt.Sprintf("%s?%s", url, queryParams), nil)
		if err != nil {
			return
		}

		if token != "" {
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil || res.StatusCode != 200 {
			res.Body.Close()
			return
		}

		parsedResponse := make([]githubRepository, 0)
		json.NewDecoder(res.Body).Decode(&parsedResponse)
		res.Body.Close()

		if len(parsedResponse) == 0 {
			break
		}

		for _, repository := range parsedResponse {
			repositories <- repository
		}

		page += 1
	}
}
