"use client";

import Link from "next/link";
import { LogOut } from "lucide-react";
import { api } from "@/lib/api";
import Button from "@/components/ui/Button";

export default function Header() {
  return (
    <nav className="bg-white shadow-sm border-b border-gray-200">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between h-16 items-center">
          <div className="flex items-center">
            <Link href="/dashboard" className="text-xl font-bold text-red-600">
              CodeSmart
            </Link>
          </div>

          <div className="flex items-center space-x-4">
            <span className="text-sm text-gray-700">Username</span>
            <Button
              type="button"
              variant="secondary"
              className="w-full flex items-center justify-center gap-3"
              onClick={async () => {
                console.log("Clicked log out");
                await api.logout();
                window.location.href = "/";
              }}
              //   loading={loading}
            >
              <LogOut className="w-4 h-4" />
              Sign Out
            </Button>
          </div>
        </div>
      </div>
    </nav>
  );
}
