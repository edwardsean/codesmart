"use client";

import { Github } from "lucide-react";
import Button from "@/components/ui/Button";
import Input from "@/components/ui/Input";

export default function SignInForm() {
  const handleGitHubLogin = () => {
    window.location.href = `${process.env.NEXT_PUBLIC_GOLANG_API_URL}/auth/github/login`;
  };
  return (
    <div className="space-y-6">
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

      {/* Email Signup Form */}
      <form
        onSubmit={(e) => {
          e.preventDefault();
        }}
        className="space-y-4"
      >
        <Input
          label="Username"
          type="username"
          name="username"
          placeholder="Create a new username"
        />

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
          placeholder="Create a password"
        />

        <Input
          label="Confirm Password"
          type="password"
          name="Confirm Password"
          placeholder="Confirm password"
        />

        <Button type="submit" variant="primary" className="w-full">
          Create Account
        </Button>

        {/* Submit Button
        <button type="submit" disabled={loading} className="btn-primary">
          {loading ? (
            <>
              <div className="spinner"></div>
              Creating Account...
            </>
          ) : (
            "Create Account"
          )}
        </button> */}
      </form>

      {/* Terms */}
      <p className="text-xs text-gray-500 text-center">
        By creating an account, you agree to our{" "}
        <a href="#" className="text-red-600 hover:text-red-700">
          Terms of Service
        </a>{" "}
        and{" "}
        <a href="#" className="text-red-600 hover:text-red-700">
          Privacy Policy
        </a>
      </p>
    </div>
  );
}
