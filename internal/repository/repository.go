package repository

// package repository provides tools for scraping a .git directory for useful information.

import (
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// GetRepoSlug parses the "config" file in a git repo directory
// for the repository slug of the repo
func GetRepoSlug(u string) (string, error) {
	var slug string
	if strings.HasPrefix(u, "http") {
		// Parse as HTTP
		p, err := url.ParseRequestURI(u)
		if err != nil {
			return "", err
		}
		slug = strings.TrimLeft(p.Path, "/")
	} else {
		// Parse as SSH
		items := strings.SplitN(u, ":", 2)
		if len(items) != 2 {
			return "", errors.New("get slug: can't parse " + u)
		}
		slug = items[1]
	}

	if slug == "" {
		return "", errors.New("repository has no remote")
	}
	if strings.HasSuffix(slug, ".git") {
		return slug[:len(slug)-4], nil
	} else {
		return slug, nil
	}
}

// GetRepoTemplates looks for `.md` files in the subdirectory
// `.gitlab/merge_request_templates`, returning a slice of all filenames found.
// Returns an error if either directory doesn't exist.
func GetRepoTemplates() ([]string, error) {

	var files []string
	err := filepath.Walk(
		".gitlab/merge_request_templates",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(path, ".md") {
				files = append(files, path)
			}
			return nil
		},
	)
	return files, err
}
