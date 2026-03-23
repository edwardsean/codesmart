export type ProjectStep = "source" | "details" | "mode"
export type GithubInputMethod = "repos" | "url"

export type Language = "go" | "python" | "typescript" | "javascript" | "java" | "c++" | "c" | "rust" | "sql" | "markdown" | "yaml" | "json" | "css" | "html" | "bash"
export type ProjectMode = "help" | "learn"
export type SourceType = "github" | "scratch"
export type Difficulty = "beginner" | "intermediate" | "advanced"
export type ProjectStatus = "active" | "completed" | "archived"
export interface Project {
    id: number
    title: string
    description: string
    language: Language
    mode: ProjectMode
    difficulty?: Difficulty
    source_type: SourceType
    source_url?: string
    total_levels: number
    completed_levels: number
    progress_percentage: number
    status: ProjectStatus
    created_at: string
    updated_at: string
}

export interface ProjectListItem {
    id: number
    title: string
    language: Language
    mode: ProjectMode
    difficulty?: Difficulty
    source_type: SourceType
    total_levels: number
    completed_levels: number
    progress_percentage: number
    status: ProjectStatus
    updated_at: string
}