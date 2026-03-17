"use client";
import { setTheme } from "@/app/actions/theme";
import { useState, useEffect } from "react";

export default function ThemeToggle() {
  const [theme, setLocalTheme] = useState<"light" | "dark">("dark");

  useEffect(() => {
    //read from cookie on mount
    const current =
      (document.cookie
        .split("; ")
        .find((row) => row.startsWith("theme="))
        ?.split("=")[1] as "light" | "dark") ?? "dark";
    setLocalTheme(current);
  }, []);

  async function toggle() {
    const next = theme === "dark" ? "light" : "dark";
    setLocalTheme(next); //instant UI update
    await setTheme(next); //persist to cookie
  }

  return (
    <button
      onClick={toggle}
      style={{
        width: 32,
        height: 32,
        borderRadius: 6,
        border: "1px solid var(--border)",
        background: "transparent",
        cursor: "pointer",
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        color: "var(--fg-muted)",
        transition: "color 0.2s",
      }}
      aria-label="Toggle theme"
    >
      {theme === "dark" ? (
        <svg
          width="14"
          height="14"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          strokeWidth="1.5"
          strokeLinecap="round"
        >
          <circle cx="12" cy="12" r="4" />
          <path d="M12 2v2M12 20v2M4.93 4.93l1.41 1.41M17.66 17.66l1.41 1.41M2 12h2M20 12h2M4.93 19.07l1.41-1.41M17.66 6.34l1.41-1.41" />
        </svg>
      ) : (
        <svg
          width="14"
          height="14"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          strokeWidth="1.5"
          strokeLinecap="round"
        >
          <path d="M21 12.79A9 9 0 1111.21 3 7 7 0 0021 12.79z" />
        </svg>
      )}
    </button>
  );
}
