package src

import (
	"errors"
	"os"
	"os/exec"
	"path"
	"strings"
)

func SynchronizeRepository(baseDir string, url string) error {
	err := os.MkdirAll(baseDir, os.ModePerm)
	if err != nil {
		return err
	}

	repositoryName := extractRepositoryName(url)
	if repositoryName == "" {
		return errors.New("Failed to extract repository name: " + url)
	}

	repositoryName = repositoryName + ".git"
	repositoryDir := path.Join(baseDir, repositoryName)

	if !checkRepository(repositoryDir) {
		return cloneRepository(url, baseDir, repositoryName)
	}

	return syncRepository(repositoryDir)
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

func cloneRepository(url string, baseDir string, name string) error {
	cmd := exec.Command("git", "clone", "--mirror", url, name)
	cmd.Dir = baseDir

	data, err := cmd.CombinedOutput()
	if err != nil && data != nil {
		return errors.New(strings.TrimSpace(string(data)))
	}

	return err
}

func syncRepository(path string) error {
	cmd := exec.Command("git", "remote", "update", "--prune")
	cmd.Dir = path

	data, err := cmd.CombinedOutput()
	if err != nil && data != nil {
		return errors.New(strings.TrimSpace(string(data)))
	}

	return err
}
