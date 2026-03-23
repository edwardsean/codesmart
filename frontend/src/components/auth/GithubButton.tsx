"use client";

import { useState } from "react";
import Button from "@/components/ui/Button";
import { Github } from "lucide-react";
import useAxiosPrivate from "@/hooks/useAxiosPrivate";
import { githubService } from "@/services/github.service";
import axios from "axios";

interface GithubButtonProps {
  connect?: boolean;
  redirect?: string;
}

export default function GithubButton({
  connect = false,
  redirect = "/dashboard",
}: GithubButtonProps) {
  const api = useAxiosPrivate();
  const apiGithub = githubService(api);
  const [loading, setLoading] = useState(false);

  async function handleClick() {
    const encoded = encodeURIComponent(redirect);

    if (connect) {
      setLoading(true);
      try {
        const response = await apiGithub.connectGithub(encoded);
        const connect_code = response.connect_code;

        window.location.href = `${process.env.NEXT_PUBLIC_BASE_URL}/api/auth/github/connect?code=${connect_code}`;
      } catch (error) {
        if (axios.isAxiosError(error)) {
          console.error(
            "error when connecting to github: ",
            error.response?.data.message,
          );
        }
        setLoading(false);
      }
    } else {
      window.location.href = `${process.env.NEXT_PUBLIC_BASE_URL}/api/auth/github/login?redirect=${encoded}`;
    }
  }

  return (
    <Button
      type="button"
      variant="ghost"
      onClick={handleClick}
      loading={loading}
      className="w-full flex items-center justify-center gap-2.5 py-2.5 text-sm"
    >
      <Github size={15} />
      {connect ? "Connect GitHub" : "Continue with GitHub"}
    </Button>
  );
}
