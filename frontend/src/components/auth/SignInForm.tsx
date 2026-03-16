"use client";

import { useSearchParams, useRouter } from "next/navigation";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { registerSchema, RegisterFormData } from "@/lib/schemas/authSchema";
import { Github } from "lucide-react";
import Button from "@/components/ui/Button";
import Input from "@/components/ui/Input";
import { authService } from "@/services/authService";
import { useAuthStore } from "@/stores/authStore";
import axios from "axios";

export default function SignInForm() {
  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<RegisterFormData>({
    resolver: zodResolver(registerSchema),
  });

  const router = useRouter();
  const searchParams = useSearchParams();
  const redirectTo = searchParams.get("redirect") || "/";
  const { setAuth } = useAuthStore();

  const [error, setError] = useState<string>("");

  const handleGitHubLogin = () => {
    window.location.href = `${process.env.NEXT_PUBLIC_BASE_URL}/api/auth/github/login`;
  };

  const onSubmit = async (data: RegisterFormData) => {
    try {
      const { user, access_token } = await authService.register(data);

      if (user && access_token) {
        setAuth(user, access_token);
        router.replace(redirectTo); //may be push
      } else {
        throw new Error("Login failed");
      }
    } catch (error) {
      if (axios.isAxiosError(error)) {
        console.log("error: ", error.response);
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
      <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
        <div>
          <Input
            label="Username"
            type="text"
            placeholder="Create a new username"
            {...register("username")}
          />
          {errors.username && (
            <p className="text-red-500 text-sm mt-1">
              {errors.username.message}
            </p>
          )}
        </div>

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
            placeholder="Create a password"
            {...register("password")}
          />
          {errors.password && (
            <p className="text-red-500 text-sm mt-1">
              {errors.password.message}
            </p>
          )}
        </div>

        <div>
          <Input
            label="Confirm Password"
            type="password"
            placeholder="Confirm password"
            {...register("confirm_password")}
          />
          {errors.confirm_password && (
            <p className="text-red-500 text-sm mt-1">
              {errors.confirm_password.message}
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
          {isSubmitting ? "Creating Account..." : "Create Account"}
        </Button>
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
