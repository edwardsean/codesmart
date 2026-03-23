import axiosInstance from "@/lib/api";
import axios from "axios";
import { processQueue, refreshState } from "@/utils/auth.utils";
import { useAuthStore } from "@/stores/auth.store";

const useRefreshToken = () => {
  const { refreshAccessToken, logout, account } = useAuthStore();
  const currentUser = account?.user;

  const refresh = async () => {
    if (refreshState.isRefreshing) {
      return new Promise<string>((resolve, reject) => {
        refreshState.pendingQueue.push({ resolve, reject });
      });
    }

    refreshState.isRefreshing = true;

    try {
      const response = await axiosInstance.post("/api/auth/refresh", {
        withCredentials: true,
      });

      const { access_token } = response.data;

      if (currentUser && access_token && access_token !== "")
        refreshAccessToken(currentUser, access_token);
    } catch (error) {
      if (axios.isAxiosError(error)) {
        console.warn("Refresh failed: ", error);
      }
      processQueue(error, null);

      try {
        await axiosInstance.post(
          "/api/auth/logout",
          {},
          {
            withCredentials: true,
          },
        );
      } catch (err) {
        if (axios.isAxiosError(err)) {
          console.error("Error when logging out: ", err.response?.data);
        }
      }

      logout();
      return null;
    } finally {
      refreshState.isRefreshing = false;
    }
  };

  return refresh;
};

export default useRefreshToken;
