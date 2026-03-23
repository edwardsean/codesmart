import { AxiosInstance } from "axios"
import { Language, ProjectMode, SourceType, Difficulty } from '@/types/entity'


export interface CreateProjectRequest {
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
        // return response.data as { project: Project}
        console.log("response from server: ", response);
        return response.data.project;
    },
    getProjectByID: async (id: number) => {
        const response = await api.get(`/api/projects/${id}`);
        return response.data.project;
    },
    deleteProject: async (id: number) => {
        await api.delete(`/api/projects/${id}`);
    },
    getProjects: async () => {
        const response = await api.get(`/api/projects`);
        return response.data.projects;
    } 
})