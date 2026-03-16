"use client";

import { useState } from "react";
import Button from "@/components/ui/Button";
import Input from "@/components/ui/Input";
import { Github } from "lucide-react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { loginSchema, LoginFormData } from "@/lib/schemas/authSchema";
import { useSearchParams, useRouter } from "next/navigation";
import { authService } from "@/services/authService";
import { useAuthStore } from "@/stores/authStore";
import axios from "axios";

export default function LoginForm() {
  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<LoginFormData>({
    resolver: zodResolver(loginSchema),
  });

  const router = useRouter();
  const searchParams = useSearchParams();
  const redirectTo = searchParams.get("redirect") || "/";
  const { setAuth } = useAuthStore();

  const [error, setError] = useState<string>("");

  const handleGitHubLogin = () => {
    window.location.href = `${process.env.NEXT_PUBLIC_BASE_URL}/api/auth/github/login`;
  };

  const onSubmit = async (data: LoginFormData) => {
    try {
      const { user, access_token } = await authService.login(data);

      console.log("data: ", user, access_token);
      if (user && access_token) {
        setAuth(user, access_token);
        router.replace(redirectTo);
      } else {
        throw new Error("Login failed");
      }
    } catch (error) {
      if (axios.isAxiosError(error)) {
        setError(error.response?.data.message);
        return;
      }
      setError("Invalid Credentials, please try again");
    }
  };

  return (
    <div className="space-y-6">
      <div className="space-y-3">
        <Button
          type="button"
          variant="secondary"
          className="w-full flex items-center justify-center gap-3"
          onClick={handleGitHubLogin}
        >
          <Github className="w-5 h-5" />
          Continue with GitHub
        </Button>
      </div>

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

      <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
        <div>
          <Input
            label="Email"
            type="email"
            placeholder="Enter your email"
            {...register("email")}
          />
          {errors.email && (
            <p className="text-red-500 text-sm mt-1">{errors.email.message}</p>
          )}
        </div>

        <div>
          <Input
            label="Password"
            type="password"
            placeholder="Enter your password"
            {...register("password")}
          />
          {errors.password && (
            <p className="text-red-500 text-sm mt-1">
              {errors.password.message}
            </p>
          )}
        </div>

        {error && <p className="text-red-500 text-sm mt-1">{error}</p>}

        <Button
          type="submit"
          variant="primary"
          className="w-full"
          disabled={isSubmitting}
        >
          {isSubmitting ? "Signing in..." : "Sign In with Email"}
        </Button>
      </form>
    </div>
  );
}
