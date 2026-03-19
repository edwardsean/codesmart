import { User } from "@/types/entity";
import { AxiosInstance } from 'axios';

export const authService = (api: AxiosInstance) => {
    return {
        login: async (data: {email: string, password: string}) => {
            const response = await api.post("/api/auth/login", data, {withCredentials: true});
            return response.data as { user: User, access_token: string };
        }, 
        register: async (data: {username: string, email: string, password: string, confirm_password: string}) => {
            const response = await api.post("/api/auth/register", data, {withCredentials: true});
            return response.data as { user: User, access_token: string};
        },
        me: async () => {
            const response = await api.get("/api/auth/me");
            return response.data as { user: User };
        }
    }
    
    
}