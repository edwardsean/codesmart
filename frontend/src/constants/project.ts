import { Language, Difficulty, ProjectStatus } from "@/types/entity"
import { ProjectStep } from "@/types/project.types"

export const LANGUAGES: Language[] = [
  "go", "python", "typescript", "javascript",
  "java", "c++", "c", "rust",
  "sql", "markdown", "yaml", "json",
  "css", "html", "bash",
]

export const DIFFICULTIES: Difficulty[] = [
  "beginner", "intermediate", "advanced",
]

export const PROJECT_STATUSES: ProjectStatus[] = [
  "active", "completed", "archived",
]

export const STEPS: ProjectStep[] = [
  "source", "details", "mode",
]