import Link from "next/link";
import { LucideIcon, ArrowRight } from "lucide-react";

interface QuickActionCardProps {
  href: string;
  label: string;
  description: string;
  icon: LucideIcon;
  accent?: boolean; // true = red icon bg, false = gray
}

export default function QuickActionCard({
  href,
  label,
  description,
  icon: Icon,
  accent = false,
}: QuickActionCardProps) {
  return (
    <Link
      href={href}
      className="group flex items-center gap-4 p-5 rounded-xl border border-black/[0.06] dark:border-white/[0.06] bg-gray-50/50 dark:bg-zinc-900/50 hover:border-black/10 dark:hover:border-white/10 hover:bg-gray-100/50 dark:hover:bg-zinc-800/50 transition-all"
    >
      <div
        className={`w-9 h-9 rounded-lg flex items-center justify-center flex-shrink-0 ${
          accent ? "bg-[#dc503c]/10" : "bg-gray-200 dark:bg-zinc-800"
        }`}
      >
        <Icon
          size={16}
          className={
            accent ? "text-[#dc503c]" : "text-gray-500 dark:text-zinc-400"
          }
        />
      </div>
      <div>
        <p className="text-sm font-medium text-gray-900 dark:text-zinc-100">
          {label}
        </p>
        <p className="text-xs text-gray-500 dark:text-zinc-500 mt-0.5">
          {description}
        </p>
      </div>
      <ArrowRight
        size={14}
        className="ml-auto text-gray-300 dark:text-zinc-600 group-hover:text-gray-500 dark:group-hover:text-zinc-400 transition-colors"
      />
    </Link>
  );
}
