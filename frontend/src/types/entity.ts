
// these are domain entities that mirror backend

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
export const LANGUAGES = ["go", "python", "typescript", "javascript", "java", "c++", "c", "rust", "sql", "markdown", "yaml", "json", "css", "html", "bash"] as const
export const DIFFICULTIES = ["beginner", "intermediate", "advanced"] as const
export const STEPS = ["source", "details", "mode"] as const;
export const PROJECT_MODE = ["help", "learn"] as const;
export const SOURCE_TYPE = ["github", "scratch"] as const
export const PROJECT_STATUS = ["active", "completed"] as const


export type Language = typeof LANGUAGES[number]
export type ProjectMode = typeof PROJECT_MODE[number]
export type SourceType = typeof SOURCE_TYPE[number]
export type Difficulty = typeof DIFFICULTIES[number]
export type ProjectStatus = typeof PROJECT_STATUS[number]
export type ProjectStep = typeof STEPS[number]
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

//project file
export interface FileNode {
  id?: number;
  path: string;
  name: string;
  type: "file" | "dir";
  children?: FileNode[];
}

export interface FileContent {
  id?: number;
  path: string;
  name: string;
  content: string;
}
