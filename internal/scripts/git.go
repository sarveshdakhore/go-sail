package scripts

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
)

func GitClone(ctx context.Context, projectName, templateType, templateURL string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	if templateType == "" || templateURL == "" {
		return fmt.Errorf("project template not found")
	}

	currentDir, _ := os.Getwd()

	folder := filepath.Join(currentDir, projectName)

	_, errPlainClone := git.PlainClone(
		folder,
		false,
		&git.CloneOptions{
			URL: getAbsoluteURL(templateURL),
		},
	)
	if errPlainClone != nil {
		return fmt.Errorf("repository `%v` was not cloned", templateURL)
	}

	err := os.RemoveAll(filepath.Join(folder, ".git"))
	if err != nil {
		return fmt.Errorf("failed to remove .git directory: %v", err)
	}

	// github configuration files removal
	githubFiles := []string{
		"CODEOWNERS",
		"CONTRIBUTING.md",
		"FUNDING.yml",
		"ISSUE_TEMPLATE",
		"PULL_REQUEST_TEMPLATE",
		"SECURITY.md",
	}

	for _, file := range githubFiles {
		err := os.RemoveAll(filepath.Join(folder, file))
		if err != nil {
			return fmt.Errorf("failed to remove GitHub configuration file %s: %v", file, err)
		}
	}

	// CI/CD configuration files removal
	ciFiles := []string{
		".github",
		".travis.yml",
		"circle.yml",
	}

	for _, file := range ciFiles {
		err := os.RemoveAll(filepath.Join(folder, file))
		if err != nil {
			return fmt.Errorf("failed to remove CI/CD configuration file %s: %v", file, err)
		}
	}

	// documentation files being removed
	docFiles := []string{
		"README.md",
		"CHANGELOG.md",
		"LICENSE",
	}

	for _, file := range docFiles {
		err := os.RemoveAll(filepath.Join(folder, file))
		if err != nil {
			return fmt.Errorf("failed to remove documentation file %s: %v", file, err)
		}
	}

	return nil
}

func getAbsoluteURL(templateURL string) string {
	templateURL = strings.TrimSpace(templateURL)
	u, _ := url.Parse(templateURL)

	if u.Scheme == "" {
		u.Scheme = "https"
	}

	return u.String()
}
