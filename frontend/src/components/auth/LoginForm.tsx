"use client";

import { useState } from "react";
import Input from "@/components/ui/Input";
import { Github } from "lucide-react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { loginSchema, LoginFormData } from "@/lib/schemas/authSchema";
import { useSearchParams, useRouter } from "next/navigation";
import { authService } from "@/services/authService";
import { useAuthStore } from "@/stores/authStore";
import axios from "axios";

export default function LoginForm({
  styles,
}: {
  styles: Record<string, string>;
}) {
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
  const [error, setError] = useState("");

  const handleGitHubLogin = () => {
    window.location.href = `${process.env.NEXT_PUBLIC_BASE_URL}/api/auth/github/login`;
  };

  const onSubmit = async (data: LoginFormData) => {
    try {
      const { user, access_token } = await authService.login(data);
      if (user && access_token) {
        setAuth(user, access_token);
        router.replace(redirectTo);
      }
    } catch (err) {
      if (axios.isAxiosError(err)) {
        setError(err.response?.data.message ?? "Login failed");
        return;
      }
      setError("Invalid credentials, please try again");
    }
  };

  return (
    <div style={{ display: "flex", flexDirection: "column", gap: "16px" }}>
      <button
        type="button"
        className={styles.githubBtn}
        onClick={handleGitHubLogin}
      >
        <Github size={16} />
        Continue with GitHub
      </button>

      <div className={styles.divider}>
        <span className={styles.dividerText}>or continue with email</span>
      </div>

      <form
        onSubmit={handleSubmit(onSubmit)}
        style={{ display: "flex", flexDirection: "column", gap: "12px" }}
      >
        <div>
          <Input
            label="Email"
            type="email"
            placeholder="Enter your email"
            {...register("email")}
          />
          {errors.email && (
            <p className={styles.fieldError}>{errors.email.message}</p>
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
            <p className={styles.fieldError}>{errors.password.message}</p>
          )}
        </div>

        {error && <p className={styles.formError}>{error}</p>}

        <button
          type="submit"
          className={styles.submitBtn}
          disabled={isSubmitting}
        >
          {isSubmitting ? "Signing in..." : "Sign in"}
        </button>
      </form>
    </div>
  );
}
