import Link from "next/link";
import { ButtonHTMLAttributes } from "react";
import clsx from "clsx";

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  size?: "sm" | "md" | "lg";
  variant?: "primary" | "secondary";
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
          "bg-blue-600 text-white hover:bg-blue-700": variant === "primary",
          "bg-gray-200 text-gray-900 hover:bg-gray-300":
            variant === "secondary",
        },

        {
          "opacity-50 cursor-not-allowed": loading,
        },
        className
      )}
      disabled={loading}
      {...rest}
    >
      {children}
    </button>
  );
}

export function HomeButton() {
  return (
    <Link
      href="/"
      className="text-gray-600 hover:text-red-600 text-sm font-medium transition-colors"
    >
      ← Back to Home
    </Link>
  );
}
