import { useLayoutEffect } from "react";
import useRefreshToken from "./useRefreshToken";
import { instance } from "@/lib/api";
import { useAuthStore } from "@/stores/authStore";

const useAxiosPrivate = () => {
  const { account } = useAuthStore();
  const token = account?.accessToken;
  const refresh = useRefreshToken();

  useLayoutEffect(() => {
    const requestInterceptor = instance.interceptors.request.use(
      (config) => {
        if (!config.headers["Authorization"]) {
          config.headers["Authorization"] = `Bearer ${token}`;
        }

        return config;
      },
      (error) => Promise.reject(error),
    );

    const responseInterceptor = instance.interceptors.response.use(
      (response) => response,
      async (error) => {
        const prevRequest = error?.config;

        if (error?.response?.status === 401 && !prevRequest?._retry) {
          prevRequest._retry = true;
          const newAccessToken = await refresh();
          prevRequest.headers["Authorization"] = `Bearer ${newAccessToken}`;
          return instance(prevRequest);
        }

        return Promise.reject(error);
      },
    );

    return () => {
      instance.interceptors.request.eject(requestInterceptor);
      instance.interceptors.response.eject(responseInterceptor);
    };
  }, [token, refresh]); //everytime auth or refresh function changes, we need to re-setup the interceptors

  return instance;
};

export default useAxiosPrivate;
