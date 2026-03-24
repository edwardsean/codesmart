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

func ShouldSkipFile(path string) bool {
	skipDirs := []string{
		"node_modules/", ".git/", "vendor/", ".next/",
		"dist/", "build/", "__pycache__/", ".venv/",
	}
	for _, dir := range skipDirs {
		if strings.HasPrefix(path, dir) || strings.Contains(path, "/"+dir) {
			return true
		}
	}

	skipExts := []string{
		".png", ".jpg", ".jpeg", ".gif", ".ico", ".svg",
		".woff", ".woff2", ".ttf", ".eot",
		".zip", ".tar", ".gz", ".exe", ".bin",
		".lock", // package-lock.json is fine but yarn.lock/pnpm-lock is huge
	}
	for _, ext := range skipExts {
		if strings.HasSuffix(strings.ToLower(path), ext) {
			return true
		}
	}

	return false
}
