"use client";

import clsx from "clsx";
import { forwardRef, InputHTMLAttributes } from "react";

interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  error?: string;
}

const Input = forwardRef<HTMLInputElement, InputProps>(
  ({ label, error, className, ...rest }, ref) => {
    return (
      <div className="flex flex-col gap-1.5">
        {label && (
          <label className="text-xs font-medium text-gray-500 dark:text-zinc-400 tracking-wide">
            {label}
          </label>
        )}
        <input
          ref={ref}
          className={clsx(
            "w-full px-3 py-2 rounded-lg border text-sm outline-none transition-colors",
            "bg-white dark:bg-zinc-900",
            "text-gray-900 dark:text-zinc-100",
            "placeholder:text-gray-400 dark:placeholder:text-zinc-600",
            error
              ? "border-red-500 focus:border-red-500"
              : "border-gray-200 dark:border-zinc-700 focus:border-gray-400 dark:focus:border-zinc-500",
            className,
          )}
          {...rest}
        />
        {error && <p className="text-xs text-red-500">{error}</p>}
      </div>
    );
  },
);

Input.displayName = "Input";
export default Input;
