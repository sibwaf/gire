package src

import (
	"log"
	"path"
	"regexp"

	"github.com/robfig/cron"
	"github.com/sibwaf/gire/src/config"
	"github.com/sibwaf/gire/src/providers"
)

func Main() {
	appConfig, err := config.ReadApplicationConfig()
	if err != nil {
		log.Fatalf("Failed to read application configuration\n%v\n", err)
	}

	sourceConfig, err := config.ReadSourceList(appConfig.SourcesPath)
	if err != nil {
		log.Fatalf("Failed to read source configuration: %s\n%v\n", appConfig.SourcesPath, err)
	}

	c := cron.New()

	err = c.AddFunc(appConfig.Cron, func() { synchronize(appConfig, sourceConfig) })
	if err != nil {
		log.Fatalf("Failed to setup cron with expression: %s\n%v\n", appConfig.Cron, err)
	}

	c.Run()
}

func synchronize(appConfig config.ApplicationConfig, sourceConfig []*config.Source) {
	status := []StatusEntry{}

	for _, source := range sourceConfig {
		// todo: move to app startup?
		provider, err := providers.CreateProviderFor(*source)
		if err != nil {
			status = append(status, MakeStatusEntry(source.Url, err))
			log.Printf("Failed to get a provider for source %s:\n%v\n", source.Url, err)
			continue
		}

		urls, err := provider.ListRepositories()
		if err != nil {
			status = append(status, MakeStatusEntry(source.Url, err))
			log.Printf("Failed to list repository URLs for source %s:\n%v\n", source.Url, err)
			continue
		}

		basePath := path.Join(appConfig.RepositoryPath, source.GroupName)

		for _, url := range urls {
			if !checkTextMatchesFilter(url, source.Include, source.Exclude) {
				continue
			}

			updated := true // todo

			err = SynchronizeRepository(basePath, url)
			if err != nil {
				log.Printf("Failed to synchronize repository: %s\n%v\n", url, err)
			} else if updated {
				log.Println("Updated repository:", url)
			}

			status = append(status, MakeStatusEntry(url, err))
		}
	}

	SaveStatus(path.Join(appConfig.RepositoryPath, ".gire.json"), status)
}

func checkTextMatchesFilter(s string, include [](*regexp.Regexp), exclude [](*regexp.Regexp)) bool {
	matchesInclude := len(include) == 0 || checkTextMatchesAny(s, include)
	return matchesInclude && !checkTextMatchesAny(s, exclude)
}

func checkTextMatchesAny(s string, regexps [](*regexp.Regexp)) bool {
	for _, regex := range regexps {
		if regex.MatchString(s) {
			return true
		}
	}

	return false
}
