"use client";

import Link from "next/link";
import { LogOut } from "lucide-react";
import useAxiosPrivate from "@/hooks/useAxiosPrivate";
import axios from "axios";
import { useAuthStore } from "@/stores/auth.store";
import ThemeToggle from "@/components/ui/ThemeToggle";
import { useState } from "react";
import { Skeleton } from "@/components/ui/Skeleton";

interface HeaderProps {
  showAuth?: boolean;
  showNav?: boolean;
}

export default function Header({
  showAuth = false,
  showNav = false,
}: HeaderProps) {
  const { logout, account, _hasHydrated } = useAuthStore();
  const api = useAxiosPrivate();
  const [isLoggingOut, setLoggingOut] = useState(false);

  async function handleLogout() {
    setLoggingOut(true);
    try {
      await api.post("/api/auth/logout");
    } catch (error) {
      if (axios.isAxiosError(error)) {
        console.error("Error when logout:", error.response?.data.message);
      }
    } finally {
      logout();
      window.location.href = "/";
    }
  }

  return (
    <nav className="fixed top-0 left-0 right-0 z-50 h-[65px] px-12 flex items-center justify-between border-b border-black/[0.06] dark:border-white/[0.06] bg-white/85 dark:bg-[#141414]/85 backdrop-blur-md">
      {/* Logo */}
      <Link
        href="/"
        className="flex items-center gap-2.5 font-mono text-sm font-medium text-gray-900 dark:text-zinc-100 no-underline tracking-tight"
      >
        <div className="w-7 h-7 rounded-md border border-[#dc503c]/70 flex items-center justify-center flex-shrink-0">
          <svg width="14" height="14" viewBox="0 0 14 14" fill="none">
            <path
              d="M2 4.5L5.5 7L2 9.5"
              stroke="#dc503c"
              strokeWidth="1.5"
              strokeLinecap="round"
              strokeLinejoin="round"
            />
            <path
              d="M7.5 9.5H12"
              stroke="#dc503c"
              strokeWidth="1.5"
              strokeLinecap="round"
            />
          </svg>
        </div>
        codesmart
      </Link>

      {/* Right side */}
      <div className="flex items-center gap-2">
        <ThemeToggle />

        {showNav && (
          <>
            <Link
              href="/auth/login"
              className="text-sm text-gray-500 dark:text-zinc-400 hover:text-gray-900 dark:hover:text-zinc-100 px-3 py-1.5 rounded-lg transition-colors"
            >
              Sign in
            </Link>
            <Link
              href="/auth/signup"
              className="text-sm font-medium px-4 py-1.5 rounded-lg transition-opacity bg-gray-900 dark:bg-zinc-100 text-white dark:text-zinc-900 hover:opacity-85"
            >
              Get started
            </Link>
          </>
        )}

        {showAuth && (
          <>
            {_hasHydrated ? (
              <span className="text-sm font-mono text-gray-400 dark:text-zinc-500">
                {account?.user.username}
              </span>
            ) : (
              <Skeleton className="h-4 w-20" />
            )}
            <button
              onClick={handleLogout}
              disabled={isLoggingOut}
              className="flex items-center gap-1.5 text-sm text-gray-500 dark:text-zinc-400 border border-black/[0.08] dark:border-white/[0.07] px-3 py-1.5 rounded-lg hover:text-gray-900 dark:hover:text-zinc-100 hover:border-black/20 dark:hover:border-white/20 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <LogOut size={13} />
              {isLoggingOut ? "Signing out..." : "Sign out"}
            </button>
          </>
        )}
      </div>
    </nav>
  );
}
