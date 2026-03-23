"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import axios from "@/lib/api";
import { useAuthStore } from "@/stores/auth.store";
import LoadingScreen from "@/components/ui/Loading";

export default function OAuthCallbackContent({
  code,
  redirect,
}: {
  code?: string;
  redirect?: string;
}) {
  const { setAuth } = useAuthStore();
  const router = useRouter();

  useEffect(() => {
    //receive code from auth service google callback
    console.log("code in callback frontend: ", code);

    if (!code) {
      router.replace("/auth/login");
      return;
    }

    const exchange = async () => {
      try {
        const { data } = await axios.get(
          `/api/auth/oauth/exchange?code=${code}`,
          { withCredentials: true }, //backend can set the cookie
        );
        console.log("data in callback frontend: ", data);
        setAuth(data.user, data.access_token);
        router.replace(redirect ?? "/dashboard"); //replace so user cant go back to /oauth.callback not router.push
      } catch {
        router.replace("/auth/login");
      }
    };

    exchange();
  }, []);

  return <LoadingScreen message="Signing you in..." />;
}
