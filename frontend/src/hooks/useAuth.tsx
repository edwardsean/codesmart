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
            const response = await instance.post<{ access_token: string }>(
              "/auth/refresh"
            );
            console.log("refresh token success: ", response.data.access_token);

            setAccessToken(response.data.access_token);

            originalRequest._retry = true;
            originalRequest.headers.Authorization = `Bearer ${response.data.access_token}`;
            return instance(originalRequest);
          } catch {
            setAccessToken(null);
            console.log("refresh token failed");
          }
        }

        return Promise.reject(error);
      }
    );

    return () => {
      instance.interceptors.request.eject(refreshInterceptor);
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
