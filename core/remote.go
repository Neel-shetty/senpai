package core

import (
	"fmt"
	"strings"
)

type Remote struct {
	Name string
	URL  string
}

func ListRemotes(repoPath string) ([]Remote, error) {
	cfg, err := ParseConfig(repoPath)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	var remotes []Remote
	for section, keys := range cfg.Sections {
		if strings.HasPrefix(section, "remote ") {
			remoteName := strings.Trim(section[7:], "\"")
			url := keys["url"]
			remotes = append(remotes, Remote{
				Name: remoteName,
				URL:  url,
			})
		}
	}

	return remotes, nil
}

func AddRemote(repoPath, name, url string) error {
	cfg, err := ParseConfig(repoPath)
	if err != nil {
		return fmt.Errorf("parse config: %w", err)
	}

	section := fmt.Sprintf("remote \"%s\"", name)

	if _, exists := cfg.Sections[section]; exists {
		return fmt.Errorf("remote '%s' already exists", name)
	}

	if err := SetConfig(repoPath, section, "url", url); err != nil {
		return fmt.Errorf("set remote url: %w", err)
	}

	fetchRefspec := fmt.Sprintf("+refs/heads/*:refs/remotes/%s/*", name)
	if err := SetConfig(repoPath, section, "fetch", fetchRefspec); err != nil {
		return fmt.Errorf("set remote fetch: %w", err)
	}

	return nil
}

func RemoveRemote(repoPath, name string) error {
	cfg, err := ParseConfig(repoPath)
	if err != nil {
		return fmt.Errorf("parse config: %w", err)
	}

	section := fmt.Sprintf("remote \"%s\"", name)

	if _, exists := cfg.Sections[section]; !exists {
		return fmt.Errorf("remote '%s' does not exist", name)
	}

	delete(cfg.Sections, section)

	configPath := fmt.Sprintf("%s/%s/config", repoPath, RepoDirName)
	return writeConfig(configPath, cfg)
}

func GetRemoteURL(repoPath, name string) (string, error) {
	section := fmt.Sprintf("remote \"%s\"", name)

	url, err := GetConfig(repoPath, section, "url")
	if err != nil {
		return "", fmt.Errorf("remote '%s' does not exist", name)
	}

	return url, nil
}

func SetRemoteURL(repoPath, name, url string) error {
	section := fmt.Sprintf("remote \"%s\"", name)

	cfg, err := ParseConfig(repoPath)
	if err != nil {
		return fmt.Errorf("parse config: %w", err)
	}

	if _, exists := cfg.Sections[section]; !exists {
		return fmt.Errorf("remote '%s' does not exist", name)
	}

	if err := SetConfig(repoPath, section, "url", url); err != nil {
		return fmt.Errorf("set remote url: %w", err)
	}

	return nil
}
