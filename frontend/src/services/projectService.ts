import { Project, Language, ProjectMode, SourceType, Difficulty } from "@/types/entity"
import { AxiosInstance } from "axios"

interface CreateProjectRequest {
    title: string
    description?: string
    language: Language
    mode: ProjectMode
    source_type: SourceType
    source_url?: string
    difficulty?: Difficulty
}


export const projectService = (api: AxiosInstance) => ({
    createProject: async (payload: CreateProjectRequest) => {
        const response = await api.post("/api/projects", payload);
        return response.data as { project: Project}
    }
})