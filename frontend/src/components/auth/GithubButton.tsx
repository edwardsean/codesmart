"use client";

import Button from "@/components/ui/Button";
import { Github } from "lucide-react";

interface GithubButtonProps {
  connect?: boolean;
  redirect?: string;
}

export default function GithubButton({
  connect = false,
  redirect = "/dashboard",
}: GithubButtonProps) {
  function handleClick() {
    const encoded = encodeURIComponent(redirect);
    const endpoint = connect ? "github/connect" : "github/login";
    window.location.href = `${process.env.NEXT_PUBLIC_BASE_URL}/api/auth/${endpoint}?redirect=${encoded}`;
  }

  return (
    <Button
      type="button"
      variant="ghost"
      onClick={handleClick}
      className="w-full flex items-center justify-center gap-2.5 py-2.5 text-sm"
    >
      <Github size={15} />
      {connect ? "Connect GitHub" : "Continue with GitHub"}
    </Button>
  );
}
