interface ProgressBarProps {
  value: number; // 0–100
  showLabel?: boolean;
}

export default function ProgressBar({
  value,
  showLabel = false,
}: ProgressBarProps) {
  const clamped = Math.min(100, Math.max(0, value));

  return (
    <div className="flex items-center gap-3">
      <div className="flex-1 h-1.5 bg-gray-100 dark:bg-zinc-800 rounded-full overflow-hidden">
        <div
          className="h-full bg-[#dc503c] rounded-full transition-all duration-300"
          style={{ width: `${clamped}%` }}
        />
      </div>
      {showLabel && (
        <span className="text-xs font-mono text-gray-400 dark:text-zinc-500 w-8 text-right">
          {clamped}%
        </span>
      )}
    </div>
  );
}
