"use client";

import { useState } from "react";
import Input from "@/components/ui/Input";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { registerSchema, RegisterFormData } from "@/schemas/authSchema";
import { useSearchParams, useRouter } from "next/navigation";
import { authService } from "@/services/auth.service";
import { useAuthStore } from "@/stores/auth.store";
import axios from "axios";
import axiosInstance from "@/lib/api";

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
  const redirectTo = searchParams.get("redirect") || "/dashboard";
  const { setAuth } = useAuthStore();
  const [error, setError] = useState("");
  const apiAuth = authService(axiosInstance);

  const onSubmit = async (data: RegisterFormData) => {
    try {
      const { user, access_token } = await apiAuth.register(data);
      if (user && access_token) {
        setAuth(user, access_token);
        router.replace(redirectTo);
      }
    } catch (err) {
      if (axios.isAxiosError(err)) {
        setError(err.response?.data.message ?? "Registration failed");
        return;
      }
      setError("Something went wrong, please try again");
    }
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-3">
      <Input
        label="Username"
        type="text"
        placeholder="Create a username"
        error={errors.username?.message}
        {...register("username")}
      />

      <Input
        label="Email"
        type="email"
        placeholder="Enter your email"
        error={errors.email?.message}
        {...register("email")}
      />

      <Input
        label="Password"
        type="password"
        placeholder="Create a password"
        error={errors.password?.message}
        {...register("password")}
      />

      <Input
        label="Confirm password"
        type="password"
        placeholder="Confirm your password"
        error={errors.confirm_password?.message}
        {...register("confirm_password")}
      />

      {error && (
        <p className="text-xs text-red-500 bg-red-500/10 border border-red-500/20 rounded-lg px-3 py-2">
          {error}
        </p>
      )}

      <button
        type="submit"
        disabled={isSubmitting}
        className="w-full py-2.5 px-4 rounded-lg text-sm font-medium transition-opacity
          bg-gray-900 dark:bg-zinc-100 text-white dark:text-zinc-900
          hover:opacity-85 disabled:opacity-50 disabled:cursor-not-allowed mt-1"
      >
        {isSubmitting ? "Creating account..." : "Create account"}
      </button>

      <p className="text-xs text-gray-400 dark:text-zinc-600 text-center mt-1">
        By creating an account you agree to our{" "}
        <a href="#" className="text-[#dc503c] hover:opacity-80">
          Terms of Service
        </a>{" "}
        and{" "}
        <a href="#" className="text-[#dc503c] hover:opacity-80">
          Privacy Policy
        </a>
      </p>
    </form>
  );
}
