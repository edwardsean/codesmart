import { Home, GitBranch, Settings, LucideIcon } from "lucide-react";

export interface NavItem {
  name: string;
  href: string;
  icon: LucideIcon;
}

export const NAV_LINKS: NavItem[] = [
  { name: "Dashboard", href: "/dashboard", icon: Home },
  { name: "Projects", href: "/dashboard/projects", icon: GitBranch },
  { name: "Settings", href: "/dashboard/settings", icon: Settings },
];