import { Repository } from "@/types/entity";
import { AxiosInstance } from "axios"


export const githubService = (api: AxiosInstance) => ({
    getRepositories: async () => {
        const response = await api.get("/api/github/repositories");
        return response.data as { repositories: Repository[]}
    }
})