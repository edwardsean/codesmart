import Link from "next/link";

interface SectionHeaderProps {
  label: string;
  viewAllHref?: string;
}

export default function SectionHeader({
  label,
  viewAllHref,
}: SectionHeaderProps) {
  return (
    <div className="flex items-center justify-between mb-4">
      <p className="text-xs font-mono text-gray-400 dark:text-zinc-500 tracking-wide uppercase">
        {label}
      </p>
      {viewAllHref && (
        <Link
          href={viewAllHref}
          className="text-xs text-gray-400 dark:text-zinc-500 hover:text-[#dc503c] transition-colors"
        >
          View all →
        </Link>
      )}
    </div>
  );
}
