import axiosInstance from '@/lib/api';
import { User } from "@/types/entity";

export const authService = {
    // me: () => instance.get<{}>("/auth/me").then(res => res.data),
    login: async (data: {email: string, password: string}) => {
        const response = await axiosInstance.post("/api/auth/login", data, {withCredentials: true});
        return response.data as { user: User, access_token: string };
    }, 
    register: async (data: {username: string, email: string, password: string, confirm_password: string}) => {
        const response = await axiosInstance.post("/api/auth/register", data, {withCredentials: true});
        return response.data as { user: User, access_token: string};
    }
}