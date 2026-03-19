import clsx from "clsx";

interface SkeletonProps {
  className?: string;
}

export function Skeleton({ className }: SkeletonProps) {
  return (
    <div
      className={clsx(
        "bg-gray-200 dark:bg-zinc-800 rounded animate-pulse",
        className,
      )}
    />
  );
}
