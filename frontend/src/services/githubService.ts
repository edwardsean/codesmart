import { Repository } from "@/types/entity";
import { AxiosInstance } from "axios"


export const githubService = (api: AxiosInstance) => ({
    connectGithub: async (encoded: string) => {
        const response = await api.post(`/api/auth/github/connect/init?redirect=${encoded}`);
        return response.data as { connect_code: string };
    },
    getRepositories: async () => {
        const response = await api.get("/api/github/repositories");
        return response.data as { repositories: Repository[]}
    }
})