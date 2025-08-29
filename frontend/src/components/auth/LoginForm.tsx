"use client";

import Button from "@/components/ui/Button";
import Input from "@/components/ui/Input";
import { Github } from "lucide-react";

export default function LoginForm() {
  const handleGitHubLogin = () => {
    window.location.href = `${process.env.NEXT_PUBLIC_GOLANG_API_URL}/auth/github/login`;
  };
  return (
    <div className="space-y-6">
      {/* OAuth Providers */}
      <div className="space-y-3">
        <Button
          type="button"
          variant="secondary"
          className="w-full flex items-center justify-center gap-3"
          onClick={handleGitHubLogin}
          //   loading={loading}
        >
          <Github className="w-5 h-5" />
          Continue with GitHub
        </Button>
      </div>

      {/* Divider */}
      <div className="relative">
        <div className="absolute inset-0 flex items-center">
          <div className="w-full border-t border-gray-300" />
        </div>
        <div className="relative flex justify-center text-sm">
          <span className="px-2 bg-white text-gray-500">
            Or continue with email
          </span>
        </div>
      </div>

      {/* Email/Password Form */}
      <form
        onSubmit={(e) => {
          e.preventDefault();
        }}
        className="space-y-4"
      >
        <Input
          label="Email"
          type="email"
          name="email"
          placeholder="Enter your email"
        />

        <Input
          label="Password"
          type="password"
          name="password"
          placeholder="Enter your password"
        />

        <Button type="submit" variant="primary" className="w-full">
          Sign In with Email
        </Button>
      </form>
    </div>
  );
}
