"use client";

import {
  createContext,
  useContext,
  useState,
  useEffect,
  useLayoutEffect,
} from "react";
import { AuthContextType } from "@/types/context";
import { api, instance } from "@/lib/api";

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  //   const [user, setUser] = useState<User | null>(null);
  const [access_token, setAccessToken] = useState<string | null>(null);

  useEffect(() => {
    const fetchMe = async () => {
      try {
        const { access_token } = await api.me();
        setAccessToken(access_token);
      } catch {
        setAccessToken(null);
      }
    };

    fetchMe();
  }, []);

  useLayoutEffect(() => {
    //everytime the token changes, update axios instance to put token in header
    const authInterceptor = instance.interceptors.request.use((config) => {
      config.headers = config.headers || {};
      config.headers.Authorization =
        !config._retry && access_token
          ? `Bearer ${access_token}`
          : config.headers.Authorization;
      return config;
    });

    return () => {
      instance.interceptors.request.eject(authInterceptor);
    };
  }, [access_token]);

  useLayoutEffect(() => {
    //we're only checking if there is an error
    const refreshInterceptor = instance.interceptors.response.use(
      (res) => res,
      async (error) => {
        const originalRequest = error.config;

        console.error("Interceptor error:", error);
        if (
          error.response.status === 401 &&
          error.response.data.message == "permission denied"
        ) {
          console.log("refreshing token");
          try {
            const { access_token } = await api.refresh();
            console.log("refresh token success: ", access_token);

            setAccessToken(access_token);

            originalRequest._retry = true;
            originalRequest.headers.Authorization = `Bearer ${access_token}`;
            return instance(originalRequest);
          } catch {
            setAccessToken(null);
            console.log("refresh token failed");
          }
        }

        return Promise.reject(error); //passess either 401 if refresh token fail or other error if error is not 401
      }
    );

    return () => {
      instance.interceptors.request.eject(refreshInterceptor); //cleanup
    };
  }, []);

  const value: AuthContextType = { access_token };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) {
    throw new Error("useAuth must be used within an AuthProvider");
  }

  return ctx;
}

//AuthProvider = Radio station that broadcasts (provides) auth state
//useAuth() = Radio receiver that listens to the broadcast (consumes auth state)
//Any component can "tune in" to get auth state without props
