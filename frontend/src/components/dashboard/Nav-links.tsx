"use client";

import { Home, GitBranch, Settings } from "lucide-react";
import Link from "next/link";
import { usePathname } from "next/navigation";

const navigation = [
  { name: "Dashboard", href: "/dashboard", icon: Home },
  { name: "Projects", href: "/dashboard/projects", icon: GitBranch },
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
            className={`
              flex items-center gap-2.5 px-3 py-2 rounded-lg text-sm transition-colors
              ${
                isActive
                  ? "bg-[#dc503c]/10 text-[#dc503c] font-medium"
                  : "text-gray-500 dark:text-zinc-400 hover:bg-black/[0.04] dark:hover:bg-white/[0.04] hover:text-gray-900 dark:hover:text-zinc-100"
              }
            `}
          >
            <Icon size={15} />
            {item.name}
          </Link>
        );
      })}
    </>
  );
}
