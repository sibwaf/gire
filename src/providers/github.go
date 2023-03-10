package providers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/sibwaf/gire/src/util"
	"golang.org/x/sync/errgroup"
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

func (p GithubProvider) GetGroupName() string {
	return extractName(p.Url)
}

func (p GithubProvider) ListRepositories() ([]string, error) {
	name := extractName(p.Url)
	if name == "" {
		return nil, errors.New(fmt.Sprintf("Failed to extract user/organization name from url %s", p.Url))
	}

	wg, ctx := errgroup.WithContext(context.Background())
	repositories := make(chan githubRepository)

	wg.Go(func() error {
		return callGithubApi(fmt.Sprintf("%s/users/%s/repos", GITHUB_API_HOST, name), p.AuthToken, repositories, ctx)
	})

	if p.AuthToken != "" {
		wg.Go(func() error {
			return callGithubApi(fmt.Sprintf("%s/user/repos", GITHUB_API_HOST), p.AuthToken, repositories, ctx)
		})
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

	return result.ToSlice(), wg.Wait()
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

// todo: ratelimit handling?
func callGithubApi(url string, token string, repositories chan githubRepository, ctx context.Context) error {
	page := 1

	for true {
		queryParams := fmt.Sprintf("page=%d&per_page=%d", page, GITHUB_PAGE_SIZE)

		req, err := http.NewRequest("GET", fmt.Sprintf("%s?%s", url, queryParams), nil)
		if err != nil {
			return err
		}

		if token != "" {
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil || res.StatusCode != 200 {
			if err == nil {
				text, _ := io.ReadAll(res.Body)
				err = errors.New(fmt.Sprintf("HTTP %d: %s", res.StatusCode, text))
			}

			res.Body.Close()
			return err
		}

		parsedResponse := make([]githubRepository, 0)
		json.NewDecoder(res.Body).Decode(&parsedResponse)
		res.Body.Close()

		if len(parsedResponse) == 0 {
			break
		}

		for _, repository := range parsedResponse {
			select {
			case repositories <- repository:
				continue
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		page += 1
	}

	return nil
}
