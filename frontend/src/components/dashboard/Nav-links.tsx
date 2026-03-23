"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { NAV_LINKS } from "@/constants/navigation";

export default function NavLinks() {
  const pathname = usePathname();

  return (
    <>
      {NAV_LINKS.map((item) => {
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
