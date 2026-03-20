interface BadgeProps {
  label: string;
  variant?: "default" | "success" | "warning" | "danger" | "accent";
}

const variants = {
  default: "bg-gray-100 dark:bg-zinc-800 text-gray-600 dark:text-zinc-400",
  success:
    "bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-400",
  warning:
    "bg-amber-100 dark:bg-amber-900/30 text-amber-700 dark:text-amber-400",
  danger: "bg-red-100 dark:bg-red-900/30 text-red-700 dark:text-red-400",
  accent: "bg-[#dc503c]/10 text-[#dc503c]",
};

export default function Badge({ label, variant = "default" }: BadgeProps) {
  return (
    <span
      className={`inline-flex items-center px-2 py-0.5 rounded-md text-xs font-mono ${variants[variant]}`}
    >
      {label}
    </span>
  );
}
