import { Difficulty, Language, ProjectStep } from "@/types/entity"

export const LANGUAGES: Language[] = ["go", "python", "typescript", "javascript", "java", "c++"] as const
export const DIFFICULTIES: Difficulty[] = ["beginner", "intermediate", "advanced"] as const
export const STEPS: ProjectStep[] = ["source", "details", "mode"];