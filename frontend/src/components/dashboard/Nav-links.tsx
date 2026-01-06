"use client";

import { Home, GitBranch, Settings } from "lucide-react";
import Link from "next/link";
import { usePathname } from "next/navigation";

const navigation = [
  { name: "Dashboard", href: "/dashboard", icon: Home },
  {
    name: "Repositories",
    href: "/dashboard/repositories",
    icon: GitBranch,
  },
  { name: "Settings", href: "/dashboard/settings", icon: Settings },
];

export default function NavLinks() {
  const pathname = usePathname();

  return (
    <>
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
    </>
  );
}
