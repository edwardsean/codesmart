"use client";

import { Github } from "lucide-react";
import { Skeleton } from "@/components/ui/Skeleton";
import RepositoryCard from "@/components/github/RepositoryCard";
import { Repository } from "@/types/repository.types";
import { useAuthStore } from "@/stores/auth.store";
import GithubButton from "@/components/auth/GithubButton";

interface RepositoryListProps {
  repos: Repository[];
  loading: boolean;
  error: string;
  selectedId?: string;
  onSelect: (repo: Repository) => void;
}

export default function RepositoryList({
  repos,
  loading,
  error,
  selectedId,
  onSelect,
}: RepositoryListProps) {
  const { account, _hasHydrated } = useAuthStore();
  const isGithubConnected =
    !!account?.user.github_id && account.user.github_id !== null;
  if (_hasHydrated) console.log("account: ", account);
  // not connected to github
  if (!isGithubConnected && _hasHydrated) {
    return (
      <div className="rounded-xl border border-dashed border-black/[0.08] dark:border-white/[0.08] p-8 flex flex-col items-center text-center gap-4">
        <div className="w-10 h-10 rounded-lg bg-gray-100 dark:bg-zinc-800 flex items-center justify-center">
          <Github size={18} className="text-gray-400 dark:text-zinc-500" />
        </div>
        <div>
          <p className="text-sm font-medium text-gray-700 dark:text-zinc-300 mb-1">
            GitHub not connected
          </p>
          <p className="text-xs text-gray-400 dark:text-zinc-600 max-w-xs">
            Connect your GitHub account to import repositories directly.
          </p>
        </div>
        <GithubButton connect redirect="/dashboard/projects/new" />
      </div>
    );
  }

  if (loading || !_hasHydrated) {
    return (
      <div className="space-y-2">
        {[...Array(5)].map((_, i) => (
          <Skeleton key={i} className="h-12" />
        ))}
      </div>
    );
  }

  if (error) {
    return (
      <p className="text-xs text-red-500 bg-red-500/10 border border-red-500/20 rounded-lg px-3 py-2">
        {error}
      </p>
    );
  }

  if (repos.length === 0) {
    return (
      <p className="text-sm text-gray-400 dark:text-zinc-500 text-center py-8">
        No repositories found.
      </p>
    );
  }

  return (
    <div className="border border-black/[0.06] dark:border-white/[0.06] rounded-xl overflow-hidden max-h-64 overflow-y-auto">
      {repos.map((repo, i) => (
        <div
          key={repo.id}
          className={
            i !== 0
              ? "border-t border-black/[0.04] dark:border-white/[0.04]"
              : ""
          }
        >
          <RepositoryCard
            repo={repo}
            selected={selectedId === repo.id}
            onSelect={onSelect}
          />
        </div>
      ))}
    </div>
  );
}
