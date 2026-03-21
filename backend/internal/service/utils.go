package service

import (
	"fmt"
	"path/filepath"
	"strings"
)

func ParseGithubURL(url string) (owner, repo string, err error) {
	// https://github.com/owner/repo
	parts := strings.Split(strings.TrimPrefix(url, "https://github.com/"), "/")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid github url: %s", url)
	}
	return parts[0], parts[1], nil
}

func DetectLanguage(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	langMap := map[string]string{
		".go":   "go",
		".ts":   "typescript",
		".tsx":  "typescript",
		".js":   "javascript",
		".jsx":  "javascript",
		".py":   "python",
		".java": "java",
		".cpp":  "c++",
		".c":    "c",
		".rs":   "rust",
		".sql":  "sql",
		".md":   "markdown",
		".yaml": "yaml",
		".yml":  "yaml",
		".json": "json",
		".css":  "css",
		".html": "html",
		".sh":   "bash",
	}
	if lang, ok := langMap[ext]; ok {
		return lang
	}
	return ""
}
