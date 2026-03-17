"use client";

import Link from "next/link";
import { LogOut } from "lucide-react";
import useAxiosPrivate from "@/hooks/useAxiosPrivate";
import axios from "axios";
import { useAuthStore } from "@/stores/authStore";
import ThemeToggle from "@/components/ui/ThemeToggle";
import { useState } from "react";
import styles from "./Header.module.css";

interface AppHeaderProps {
  showAuth?: boolean; //true = show logout (dashboard)
  showNav?: boolean; //true = show sign in / get started (home)
}

export default function Header({
  showAuth = false,
  showNav = false,
}: AppHeaderProps) {
  const { logout, account } = useAuthStore();
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
    <nav className={styles.nav}>
      <Link href="/" className={styles.logo}>
        <div className={styles.logoIcon}>
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

      <div className={styles.navLinks}>
        <ThemeToggle />

        {/* home page — not logged in */}
        {showNav && (
          <>
            <Link href="/auth/login" className={styles.link}>
              Sign in
            </Link>
            <Link href="/auth/signup" className={styles.btn}>
              Get started
            </Link>
          </>
        )}

        {/* dashboard — logged in */}
        {showAuth && (
          <>
            <span className={styles.username}>{account?.user.username}</span>
            <button
              className={styles.logoutBtn}
              onClick={handleLogout}
              disabled={isLoggingOut}
            >
              <LogOut size={14} />
              {isLoggingOut ? "Signing out..." : "Sign out"}
            </button>
          </>
        )}
      </div>
    </nav>
  );
}
