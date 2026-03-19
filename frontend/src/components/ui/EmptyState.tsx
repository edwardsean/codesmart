import Link from "next/link";
import { LucideIcon } from "lucide-react";

interface EmptyStateProps {
  icon: LucideIcon;
  title: string;
  description: string;
  action?: {
    label: string;
    href: string;
  };
}

export default function EmptyState({
  icon: Icon,
  title,
  description,
  action,
}: EmptyStateProps) {
  return (
    <div className="rounded-xl border border-dashed border-black/[0.08] dark:border-white/[0.08] p-12 flex flex-col items-center justify-center text-center">
      <div className="w-10 h-10 rounded-lg bg-gray-100 dark:bg-zinc-800 flex items-center justify-center mb-4">
        <Icon size={18} className="text-gray-400 dark:text-zinc-500" />
      </div>
      <p className="text-sm font-medium text-gray-700 dark:text-zinc-300 mb-1">
        {title}
      </p>
      <p className="text-xs text-gray-400 dark:text-zinc-600 mb-5 max-w-xs">
        {description}
      </p>
      {action && (
        <Link
          href={action.href}
          className="text-sm font-medium px-4 py-2 rounded-lg bg-gray-900 dark:bg-zinc-100 text-white dark:text-zinc-900 hover:opacity-85 transition-opacity"
        >
          {action.label}
        </Link>
      )}
    </div>
  );
}
