package core

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type GitConfig struct {
	Sections map[string]map[string]string
}

func ParseConfig(repoPath string) (*GitConfig, error) {
	configPath := filepath.Join(repoPath, RepoDirName, "config")
	f, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("open config: %w", err)
	}
	defer f.Close()

	cfg := &GitConfig{Sections: make(map[string]map[string]string)}

	scanner := bufio.NewScanner(f)
	currentSection := ""

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			section := strings.TrimSpace(line[1 : len(line)-1])
			currentSection = section
			if _, ok := cfg.Sections[currentSection]; !ok {
				cfg.Sections[currentSection] = make(map[string]string)
			}

			continue
		}

		if eq := strings.Index(line, "="); eq != -1 {
			key := strings.TrimSpace(line[:eq])
			val := strings.TrimSpace(line[eq+1:])
			if currentSection == "" {
				return nil, fmt.Errorf("key-value pair outside section: %s", line)
			}
			cfg.Sections[currentSection][key] = val
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan config: %w", err)
	}

	return cfg, nil
}

func ListConfig(repoPath string) error {
	cfg, err := ParseConfig(repoPath)
	if err != nil {
		return fmt.Errorf("parse config: %w", err)
	}

	for section, keys := range cfg.Sections {
		for key, value := range keys {
			fmt.Printf("%s.%s=%s\n", section, key, value)
		}
	}

	return nil
}

func GetConfig(repoPath, section, key string) (string, error) {
	cfg, err := ParseConfig(repoPath)
	if err != nil {
		return "", fmt.Errorf("parse config: %w", err)
	}

	sectionMap, ok := cfg.Sections[section]
	if !ok {
		return "", fmt.Errorf("section '%s' not found", section)
	}

	value, ok := sectionMap[key]
	if !ok {
		return "", fmt.Errorf("key '%s' not found in section '%s'", key, section)
	}

	return value, nil
}

func SetConfig(repoPath, section, key, value string) error {
	cfg, err := ParseConfig(repoPath)
	if err != nil {
		return fmt.Errorf("parse config: %w", err)
	}

	if _, ok := cfg.Sections[section]; !ok {
		cfg.Sections[section] = make(map[string]string)
	}

	cfg.Sections[section][key] = value

	configPath := filepath.Join(repoPath, RepoDirName, "config")
	return writeConfig(configPath, cfg)
}

func writeConfig(configPath string, cfg *GitConfig) error {
	f, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("create config file: %w", err)
	}
	defer f.Close()

	writer := bufio.NewWriter(f)
	defer writer.Flush()

	for section, keys := range cfg.Sections {
		if _, err := fmt.Fprintf(writer, "[%s]\n", section); err != nil {
			return fmt.Errorf("write section: %w", err)
		}

		for key, value := range keys {
			if _, err := fmt.Fprintf(writer, "\t%s = %s\n", key, value); err != nil {
				return fmt.Errorf("write key-value: %w", err)
			}
		}

		if _, err := writer.WriteString("\n"); err != nil {
			return fmt.Errorf("write newline: %w", err)
		}
	}

	return nil
}

func CreateDefaultConfig(repoPath string, isBare bool) error {
	cfg := &GitConfig{
		Sections: make(map[string]map[string]string),
	}

	cfg.Sections["core"] = map[string]string{
		"repositoryformatversion": "0",
		"filemode":                "true",
		"bare":                    fmt.Sprintf("%t", isBare),
		"logallrefupdates":        "true",
	}

	configPath := filepath.Join(repoPath, RepoDirName, "config")
	return writeConfig(configPath, cfg)
}
