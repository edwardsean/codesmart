import { Check } from "lucide-react";
import { Repository } from "@/types/repository.types";
import { getLanguageColor } from "@/utils/utils";

interface RepositoryCardProps {
  repo: Repository;
  selected?: boolean;
  onSelect?: (repo: Repository) => void;
}

export default function RepositoryCard({
  repo,
  selected = false,
  onSelect,
}: RepositoryCardProps) {
  return (
    <button
      onClick={() => onSelect?.(repo)}
      className={`w-full flex items-center justify-between px-4 py-3 text-left transition-colors
        ${selected ? "bg-[#dc503c]/5" : "hover:bg-gray-50 dark:hover:bg-zinc-800/50"}`}
    >
      <div className="flex items-center gap-3 min-w-0">
        {/* language dot */}
        <div
          className={`w-2.5 h-2.5 rounded-full flex-shrink-0 ${getLanguageColor(repo.language)}`}
        />

        <div className="min-w-0">
          <p className="text-sm text-gray-900 dark:text-zinc-100 truncate">
            {repo.name}
          </p>
          {repo.description && (
            <p className="text-xs text-gray-400 dark:text-zinc-500 truncate mt-0.5">
              {repo.description}
            </p>
          )}
        </div>

        {repo.language && (
          <span className="text-xs font-mono text-gray-400 dark:text-zinc-600 flex-shrink-0 ml-2">
            {repo.language}
          </span>
        )}
      </div>

      {selected && (
        <Check size={13} className="text-[#dc503c] flex-shrink-0 ml-3" />
      )}
    </button>
  );
}
