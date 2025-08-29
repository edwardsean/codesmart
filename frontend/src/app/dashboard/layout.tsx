"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";

import { Home, Upload, GitBranch, Users, Settings, LogOut } from "lucide-react";

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const navigation = [
    { name: "Dashboard", href: "/dashboard", icon: Home },
    {
      name: "Add Repositories",
      href: "/dashboard/repositories",
      icon: GitBranch,
    },
    { name: "Settings", href: "/dashboard/settings", icon: Settings },
  ];

  const pathname = usePathname();

  return (
    <div className="min-h-screen bg-gray-50">
      <nav className="bg-white shadow-sm border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16 items-center">
            <div className="flex items-center">
              <Link
                href="/dashboard"
                className="text-xl font-bold text-red-600"
              >
                CodeSmart
              </Link>
            </div>

            <div className="flex items-center space-x-4">
              <span className="text-sm text-gray-700">Username</span>
              <button
                // onClick={() => signOut({ callbackUrl: "/" })}
                className="btn-secondary px-3 py-1 text-sm flex items-center gap-2"
              >
                <LogOut className="w-4 h-4" />
                Sign Out
              </button>
            </div>
          </div>
        </div>
      </nav>

      <div className="flex">
        {/* Sidebar Navigation */}
        <div className="w-64 bg-white shadow-sm border-r border-gray-200 min-h-screen">
          <nav className="mt-8">
            <div className="px-4 space-y-2">
              {navigation.map((item) => {
                const Icon = item.icon;
                const isActive = pathname === item.href;

                return (
                  <Link
                    key={item.name}
                    href={item.href}
                    className={`group flex items-center px-3 py-2 text-sm font-medium rounded-md transition-colors ${
                      isActive
                        ? "bg-red-50 text-red-600 border-r-2 border-red-600"
                        : "text-gray-700 hover:bg-gray-50 hover:text-gray-900"
                    }`}
                  >
                    <Icon className="w-5 h-5 mr-3" />
                    {item.name}
                  </Link>
                );
              })}
            </div>
          </nav>
        </div>

        {/* Main Content */}
        <div className="flex-1">{children}</div>
      </div>
    </div>
  );
}
