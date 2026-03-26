package domain

import "time"

type Project struct {
	ID          int         `json:"id"             gorm:"primaryKey"`
	UserID      int         `json:"user_id"` //Which user owns this project, cascade deletes if user is deleted
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Language    Language    `json:"language"` //Primary programming language, determines which executor to use and syntax highlighting
	Mode        ProjectMode `json:"mode"`     // "help" or "learn"

	//For docker container
	ContainerID string `json:"container_id"`
	VolumeName  string `json:"volume_name"`

	// AI fills these after analysis
	ProjectType string     `json:"project_type,omitempty"` //Category of project, e.g. "web", "cli", "api"
	Difficulty  Difficulty `json:"difficulty,omitempty"`   //shown on project card
	TotalLevels int        `json:"total_levels"`           //How many levels the AI generated for this project, shown in progress bar

	SourceType         SourceType    `json:"source_type,omitempty"`
	SourceURL          string        `json:"source_url,omitempty"`                        //GitHub repo URL if imported, stored so you can re-sync later
	GithubMetadata     []byte        `json:"github_metadata,omitempty" gorm:"type:jsonb"` //JSONB blob of GitHub repo info (stars, description, default branch), avoids extra API calls
	CompletedLevels    int           `json:"completed_levels"`
	Status             ProjectStatus `json:"status"`
	ProgressPercentage int           `json:"progress_percentage"` //Derived from completed/total levels, shown in the progress bar on the project card
	CreatedAt          time.Time     `json:"created_at"`
	UpdatedAt          time.Time     `json:"updated_at"`
}

type ProjectFile struct {
	ID        int       `json:"id"        gorm:"primaryKey"`
	ProjectID int       `json:"project_id"`
	FilePath  string    `json:"file_path"`                  //Full path like backend/cmd/main.go, used to render the file tree in the left panel of your wireframes
	Content   string    `json:"content"   gorm:"type:text"` //Raw file content, displayed in the Monaco editor
	Language  Language  `json:"language,omitempty"`         //File language for syntax highlighting ("go", "typescript", etc.)
	CreatedAt time.Time `json:"created_at"`
}
