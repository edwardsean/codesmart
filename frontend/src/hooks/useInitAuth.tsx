import { useEffect } from "react";
import { useAuthStore } from "@/stores/auth.store";
import useRefreshToken from "@/hooks/useRefreshToken";

export const useInitAuth = () => {
  const { _hasHydrated, account } = useAuthStore();
  const refresh = useRefreshToken();

  useEffect(() => {
    if (!_hasHydrated) return;

    console.log("account: ", account);

    if (account?.user && !account?.accessToken) refresh();
  }, [_hasHydrated]);
};

// const AuthContext = createContext<AuthContextType | undefined>(undefined);

// export function AuthProvider({ children }: { children: React.ReactNode }) {
//   //   const [user, setUser] = useState<User | null>(null);
//   const [access_token, setAccessToken] = useState<string | null>(null);
//   // const [loading, setLoading] = useState<boolean>(true);

//   useEffect(() => {
//     const fetchMe = async () => {
//       try {
//         // setLoading(true);
//         const { access_token } = await api.me();
//         console.log("me access token: ", access_token);
//         setAccessToken(access_token);
//       } catch (err) {
//         console.log("me fetch failed: ", err);
//         setAccessToken(null);
//       }
//     };

//     fetchMe();
//   }, [access_token]);

//   useLayoutEffect(() => {
//     //everytime the token changes, update axios instance to put token in header
//     const authInterceptor = instance.interceptors.request.use(
//       (config) => {
//         config.headers = config.headers || {};
//         config.headers.Authorization =
//           !config._retry && access_token
//             ? `Bearer ${access_token}`
//             : config.headers.Authorization;
//         return config;
//       },
//       (error) => Promise.reject(error),
//     );

//     return () => {
//       instance.interceptors.request.eject(authInterceptor);
//     };
//   }, [access_token]);

//   useLayoutEffect(() => {
//     //we're only checking if there is an error
//     const refreshInterceptor = instance.interceptors.response.use(
//       (res) => res,
//       async (error) => {
//         const originalRequest = error.config;

//         if (error.response.status === 401) {
//           console.log("refreshing token");
//           try {
//             const { access_token } = await api.refresh();
//             console.log("refresh token success: ", access_token);

//             setAccessToken(access_token);

//             originalRequest._retry = true;
//             originalRequest.headers.Authorization = `Bearer ${access_token}`;
//             return instance(originalRequest);
//           } catch {
//             setAccessToken(null);
//             console.log("refresh token failed");
//           }
//         }

//         return Promise.reject(error); //passess either 401 if refresh token fail or other error if error is not 401
//       },
//     );

//     return () => {
//       instance.interceptors.response.eject(refreshInterceptor); //cleanup
//     };
//   }, []);

//   const value: AuthContextType = { access_token };

//   return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
// }

// export function useAuth() {
//   const ctx = useContext(AuthContext);
//   if (!ctx) {
//     throw new Error("useAuth must be used within an AuthProvider");
//   }

//   return ctx;
// }

//AuthProvider = Radio station that broadcasts (provides) auth state
//useAuth() = Radio receiver that listens to the broadcast (consumes auth state)
//Any component can "tune in" to get auth state without props
