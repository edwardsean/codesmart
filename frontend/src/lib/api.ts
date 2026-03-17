import axios from 'axios';

export default axios.create({
    baseURL: process.env.BASE_URL,
});

export const instance = axios.create({
    baseURL: process.env.BASE_URL,
    withCredentials: true,
})

//make a hook for this to use the useAxiosPrivate hook to automatically add the access token to the header and handle refresh token logic
// export const api = {
//     me: () => instance.get<{access_token: string}>("/auth/me").then(res => res.data),
//     login: (email: string, password: string) => instance.post("/auth/login", {
//         email: email, 
//         password: password
//     }),
//     register: () => instance.post("/auth/register"),
//     logout: async () => {
//         try {
//             await instance.post("/auth/logout");
//         } finally {
//             localStorage.removeItem("access_token");
//         }
        

//     },
//     get_repositories: () => instance.get<{repositories: Repository[]}>("/github/getRepositories").then(res => res.data),
//     debug: () => instance.get<string>("/github/debug").then(res => res.data)
// }



