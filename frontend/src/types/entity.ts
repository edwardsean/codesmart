export type User = {
    id: number;
    email: string;
    username: string;
    github_id?: number;
    createdAt: string;
}

export interface Repository {
  id: string;
  name: string;
  full_name: string;
  description: string;
  html_url: string
  language: string
  stargazers_count: number
  forks_count: number
  updated_at: string
  private: boolean
  default_branch: string
  owner: {
    login: string
    avatar_url: string
  }
}


//project
export type Language = "go" | "python" | "typescript" | "javascript" | "java" | "c++" 
export type ProjectMode = "help" | "learn"
export type SourceType = "github" | "scratch"
export type Difficulty = "beginner" | "intermediate" | "advanced"
export type ProjectStatus = "active" | "completed" | "archived"
export type ProjectStep = "source" | "details" | "mode"
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
