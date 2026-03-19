package domain

type ProjectMode string
type SourceType string
type Difficulty string
type ProjectStatus string
type Language string

const (
	// ProjectMode
	ModeHelp  ProjectMode = "help"
	ModeLearn ProjectMode = "learn"

	// SourceType
	SourceGithub  SourceType = "github"
	SourceScratch SourceType = "scratch"

	// Difficulty
	DifficultyBeginner     Difficulty = "beginner"
	DifficultyIntermediate Difficulty = "intermediate"
	DifficultyAdvanced     Difficulty = "advanced"

	// ProjectStatus
	StatusActive    ProjectStatus = "active"
	StatusCompleted ProjectStatus = "completed"
	StatusArchived  ProjectStatus = "archived"

	// Language
	LanguageGo         Language = "go"
	LanguagePython     Language = "python"
	LanguageTypeScript Language = "typescript"
	LanguageJavaScript Language = "javascript"
	LanguageJava       Language = "java"
)
