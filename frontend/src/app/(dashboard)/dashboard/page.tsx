"use client";
import { GitBranch, Plus, Code2, Trophy, Clock } from "lucide-react";
import StatCard from "@/components/ui/StatCard";
import QuickActionCard from "@/components/ui/QuickActionCard";
import EmptyState from "@/components/ui/EmptyState";
import SectionHeader from "@/components/ui/SectionHeader";
import { useAuthStore } from "@/stores/authStore";
import { Skeleton } from "@/components/ui/Skeleton";

export default function Dashboard() {
  const { account, _hasHydrated } = useAuthStore();
  const username = account?.user.username;

  if (!_hasHydrated) {
    return (
      <div>
        <div className="mb-10">
          <Skeleton className="h-3 w-20 mb-2" />
          <Skeleton className="h-9 w-64 mb-2" />
          <Skeleton className="h-4 w-48" />
        </div>
        <div className="grid grid-cols-2 gap-3 mb-10">
          <Skeleton className="h-20" />
          <Skeleton className="h-20" />
        </div>
        <div className="grid grid-cols-3 gap-3 mb-10">
          <Skeleton className="h-24" />
          <Skeleton className="h-24" />
          <Skeleton className="h-24" />
        </div>
      </div>
    );
  }

  return (
    <div>
      {/* Welcome */}
      <div className="mb-10">
        <p className="text-xs font-mono text-[#dc503c] tracking-widest uppercase mb-2">
          dashboard
        </p>
        <h1 className="text-3xl font-serif font-normal text-gray-900 dark:text-zinc-100 tracking-tight">
          Welcome back, {username}.
        </h1>
        <p className="text-sm text-gray-500 dark:text-zinc-500 mt-1.5">
          Pick up where you left off, or start something new.
        </p>
      </div>

      {/* Quick actions */}
      <div className="grid grid-cols-2 gap-3 mb-10">
        <QuickActionCard
          href="/dashboard/projects/new"
          label="New project"
          description="From scratch or GitHub"
          icon={Plus}
          accent
        />
        <QuickActionCard
          href="/dashboard/projects"
          label="All projects"
          description="View your workspace"
          icon={GitBranch}
        />
      </div>

      {/* Stats */}
      <div className="grid grid-cols-3 gap-3 mb-10">
        <StatCard label="Projects" value="0" icon={Code2} />
        <StatCard label="XP earned" value="0" icon={Trophy} />
        <StatCard label="Hours spent" value="0" icon={Clock} />
      </div>

      {/* Recent projects */}
      <SectionHeader
        label="Recent projects"
        viewAllHref="/dashboard/projects"
      />
      <EmptyState
        icon={Code2}
        title="No projects yet"
        description="Import a GitHub repo or start from scratch — Smarty will turn it into a learning path."
        action={{
          label: "Create your first project",
          href: "/dashboard/projects/new",
        }}
      />
    </div>
  );
}
