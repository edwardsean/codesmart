"use client";

import { Home, GitBranch, Settings } from "lucide-react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import styles from "./Nav-links.module.css";

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
            className={`${styles.link} ${isActive ? styles.active : ""}`}
          >
            <Icon size={15} />
            {item.name}
          </Link>
        );
      })}
    </>
  );
}
