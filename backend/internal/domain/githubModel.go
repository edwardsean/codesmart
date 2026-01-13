package domain

type Repository struct {
	ID              int       `json:"id"`
	Name            string    `json:"name"`
	FullName        string    `json:"full_name"`
	Description     string    `json:"description"`
	HTMLURL         string    `json:"html_url"`
	Language        string    `json:"language"`
	StargazersCount int       `json:"stargazers_count"`
	ForksCount      int       `json:"forks_count"`
	UpdatedAt       string    `json:"updated_at"`
	Private         bool      `json:"private"`
	DefaultBranch   string    `json:"default_branch"`
	Owner           RepoOwner `json:"owner"`
}

type RepoOwner struct {
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
}
