import { LucideIcon } from "lucide-react";

interface StatCardProps {
  label: string;
  value: string | number;
  icon: LucideIcon;
}

export default function StatCard({ label, value, icon: Icon }: StatCardProps) {
  return (
    <div className="p-5 rounded-xl border border-black/[0.06] dark:border-white/[0.06] bg-gray-50/50 dark:bg-zinc-900/50">
      <Icon size={14} className="text-gray-400 dark:text-zinc-500 mb-3" />
      <p className="text-2xl font-mono font-medium text-gray-900 dark:text-zinc-100">
        {value}
      </p>
      <p className="text-xs text-gray-500 dark:text-zinc-500 mt-0.5">{label}</p>
    </div>
  );
}
