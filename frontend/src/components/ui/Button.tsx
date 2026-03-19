import Link from "next/link";
import { ButtonHTMLAttributes } from "react";
import clsx from "clsx";
import { Github } from "lucide-react";

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  size?: "sm" | "md" | "lg";
  variant?: "primary" | "secondary" | "ghost";
  loading?: boolean;
  children: React.ReactNode;
}

export default function Button({
  children,
  className,
  size = "md",
  variant = "primary",
  loading = false,
  ...rest
}: ButtonProps) {
  return (
    <button
      {...rest}
      className={clsx(
        "font-medium rounded-lg transition-colors",
        {
          "px-3 py-1.5 text-sm": size === "sm",
          "px-4 py-2 text-base": size === "md",
          "px-6 py-3 text-lg": size === "lg",
        },
        {
          "bg-gray-900 dark:bg-zinc-100 text-white dark:text-zinc-900 hover:opacity-85":
            variant === "primary",
          "bg-gray-100 dark:bg-zinc-800 text-gray-900 dark:text-zinc-100 hover:bg-gray-200 dark:hover:bg-zinc-700":
            variant === "secondary",
          // ← ghost: transparent with border, for GitHub/OAuth buttons
          "bg-transparent border border-black/[0.08] dark:border-white/[0.07] text-gray-700 dark:text-zinc-300 hover:bg-black/[0.04] dark:hover:bg-white/[0.04]":
            variant === "ghost",
        },
        { "opacity-50 cursor-not-allowed": loading },
        className,
      )}
      disabled={loading || rest.disabled}
    >
      {children}
    </button>
  );
}

export function HomeButton() {
  return (
    <Link
      href="/"
      className="text-gray-500 dark:text-zinc-500 hover:text-[#dc503c] text-sm transition-colors"
    >
      ← Back to home
    </Link>
  );
}
