// import { Repository, User } from "@/types/entity";
// // const API_URL = process.env.NEXT_PUBLIC_GOLANG_API_URL

// let API_URL = process.env.NEXT_PUBLIC_GOLANG_API_URL

// export class APIError extends Error {
//     status: number;
//     constructor(status: number, message: string) {
//         super(message)
//         this.status = status;
//     }
// }

// async function request<T>(url:string, access_token: string, options?: RequestInit): Promise<T> {
//     console.log(`URL: ${API_URL}/${url}`)
//     const fetchOptions: RequestInit = {
//         ...options,
//     }


//     if(access_token != "") {
//         console.log("access token in api.ts: ", access_token)
//         fetchOptions.headers = {
//             ...(options?.headers || {}),
//             Cookie: `access_token=${access_token}`
//         };

//         API_URL = process.env.DOCKER_PUBLIC_GOLANG_API_URL
//     }  else {
//         console.log("no access token in api.ts")
//         fetchOptions.credentials = "include"
//     }


//     try {
//         console.log("URL in Request: ", `${API_URL}/${url}`)
//         const res = await fetch(`${API_URL}/${url}`, fetchOptions);

//         if(!res.ok) {
//             const text = await res.text();
//             throw new APIError(res.status, text || `API error: ${res.status}`)
//         }
//         return res.json() as Promise<T>;
//     } catch(err: any) {
//         console.log(`error: ${err.message}`)
//         if(err instanceof APIError) {
//             console.log("instance of api error")
//             throw err
//         }
//         console.log("non api error")

//         throw new APIError(0, err.message || "Network Error")
        
//     }
    
// }

// async function requestVoid(url:string, access_token: string, options?: RequestInit): Promise<void> {
//     console.log(`URL: ${API_URL}/${url}`)
//      const fetchOptions: RequestInit = {
//         ...options,
//     }

//     if(access_token != "") {
//         fetchOptions.headers = {
//             ...(options?.headers || {}),
//             Cookie: `access_token=${access_token}`
//         };

//         API_URL = process.env.DOCKER_PUBLIC_GOLANG_API_URL
//     }  else {
//         fetchOptions.credentials = "include"
//     }


//     try {
//         const res = await fetch(`${API_URL}/${url}`, fetchOptions);

//         if(!res.ok) {
//             const text = await res.text();
//             throw new APIError(res.status, text || `API error: ${res.status}`)
//         }
//     } catch(err: any) {
//         console.log(`error: ${err.message}`)
//         if(err instanceof APIError) {
//             console.log("instance of api error")
//             throw err
//         }
//         console.log("non api error")

//         throw new APIError(0, err.message || "Network Error")
        
//     }

// }

// export const api = {
//     me: (access_token?: string) => request<{user_data: User}>("auth/me", access_token || "", undefined),
//     refresh: (access_token?: string) => requestVoid("auth/refresh",  access_token || "" ,{method: "POST"}),
//     login: (email: string, password: string, access_token?: string) =>  request<({token: string})>("/auth/login", access_token || "", {method: "POST", headers: {"Content-Type" : "application/json",}, body : JSON.stringify({email, password})}),
//     register: (access_token?: string) => requestVoid("auth/register", access_token || "", {method: "POST"}),
//     logout:  (access_token?: string) => requestVoid("auth/logout", access_token || "", {method: "POST"}),
//     get_repositories: (access_token?: string) => request<{repositories: Repository[]}>("getRepositories", access_token || "", undefined),
//     testing: (access_token?: string) => request<string>("debug", access_token || "", undefined)
    
// }

import {  Repository } from '@/types/entity';
import axios from 'axios';

export const instance = axios.create({
    baseURL: process.env.NEXT_PUBLIC_GOLANG_API_URL,
    withCredentials: true,
})

export const api = {
    me: () => instance.get<{access_token: string}>("/auth/me").then(res => res.data),
    login: (email: string, password: string) => instance.post("/auth/login", {
        email: email, 
        password: password
    }),
    refresh: () => instance.post<{access_token: string}>("/auth/refresh").then(res => res.data),
    register: () => instance.post("/auth/register"),
    logout: () => instance.post("/auth/logout"),
    get_repositories: () => instance.get<{repositories: Repository[]}>("/getRepositories").then(res => res.data)
}



