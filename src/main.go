package src

import (
	"log"
	"path"
	"regexp"

	"github.com/robfig/cron"
	"github.com/sibwaf/gire/src/config"
	"github.com/sibwaf/gire/src/providers"
	"github.com/sibwaf/gire/src/util"
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

		groupName := util.Coalesce(source.GroupName, provider.GetGroupName(), "_")
		basePath := path.Join(appConfig.RepositoryPath, groupName)

		for _, url := range urls {
			if !checkTextMatchesFilter(url, source.Include, source.Exclude) {
				continue
			}

			result, err := SynchronizeRepository(basePath, url)
			switch result {
			case SYNC_RESULT_OK:
				log.Println("Updated repository:", url)
			case SYNC_RESULT_UPTODATE:
				log.Println("Already up-to-date:", url)
			default:
				log.Printf("Failed to synchronize repository: %s\n%v\n", url, err)
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
