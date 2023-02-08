package src

import (
	"errors"
	"os"
	"os/exec"
	"path"
	"strings"
)

type SyncResult string

const (
	SYNC_RESULT_OK       SyncResult = "ok"
	SYNC_RESULT_UPTODATE SyncResult = "uptodate"
	SYNC_RESULT_FAIL     SyncResult = "fail"
)

func SynchronizeRepository(baseDir string, url string) (SyncResult, error) {
	err := os.MkdirAll(baseDir, os.ModePerm)
	if err != nil {
		return SYNC_RESULT_FAIL, err
	}

	repositoryName := extractRepositoryName(url)
	if repositoryName == "" {
		return SYNC_RESULT_FAIL, errors.New("Failed to extract repository name: " + url)
	}

	repositoryName = repositoryName + ".git"
	repositoryDir := path.Join(baseDir, repositoryName)

	if !checkRepository(repositoryDir) {
		return cloneRepository(url, baseDir, repositoryName)
	}

	return updateRepository(repositoryDir)
}

func extractRepositoryName(url string) string {
	url = strings.TrimSpace(url)

	for strings.HasSuffix(url, "/") {
		url = url[:len(url)-1]
	}
	if strings.HasSuffix(url, ".git") {
		url = url[:len(url)-4]
	}

	lastSlashIndex := strings.LastIndex(url, "/")
	if lastSlashIndex < 0 {
		return ""
	}

	return url[lastSlashIndex+1:]
}

func checkRepository(path string) bool {
	cmd := exec.Command("git", "show", "--summary")
	cmd.Dir = path

	return cmd.Run() == nil
}

func cloneRepository(url string, baseDir string, name string) (SyncResult, error) {
	cmd := exec.Command("git", "clone", "--mirror", url, name)
	cmd.Dir = baseDir

	data, err := cmd.CombinedOutput()
	if err != nil {
		if data != nil {
			return SYNC_RESULT_FAIL, errors.New(strings.TrimSpace(string(data)))
		}

		return SYNC_RESULT_FAIL, err
	}

	return SYNC_RESULT_OK, nil
}

func updateRepository(path string) (SyncResult, error) {
	cmd := exec.Command("git", "remote", "update", "--prune")
	cmd.Dir = path

	data, err := cmd.CombinedOutput()
	if err != nil {
		if data != nil {
			return SYNC_RESULT_FAIL, errors.New(strings.TrimSpace(string(data)))
		}

		return SYNC_RESULT_FAIL, err
	}

	if len(data) > 0 {
		return SYNC_RESULT_OK, nil
	} else {
		return SYNC_RESULT_UPTODATE, nil
	}
}
