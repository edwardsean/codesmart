"use client";

import Link from "next/link";
import { Plus, Code2 } from "lucide-react";
import SectionHeader from "@/components/ui/SectionHeader";
import EmptyState from "@/components/ui/EmptyState";
import ProjectCard from "@/components/dashboard/ProjectCard";
import { Skeleton } from "@/components/ui/Skeleton";
import { Project } from "@/types/project.types";
import { useProjects } from "@/hooks/query/useProjects";

export default function ProjectsPage() {
  const { data: projects, isLoading, error } = useProjects();
  const active = projects?.filter((p: Project) => p.status === "active") ?? [];
  const completed =
    projects?.filter((p: Project) => p.status === "completed") ?? [];

  return (
    <div>
      {/* Header */}
      <div className="flex items-center justify-between mb-10">
        <div>
          <p className="text-xs font-mono text-[#dc503c] tracking-widest uppercase mb-2">
            projects
          </p>
          <h1 className="text-3xl font-serif font-normal text-gray-900 dark:text-zinc-100 tracking-tight">
            Your projects.
          </h1>
        </div>
        <Link
          href="/dashboard/projects/new"
          className="flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium
            bg-gray-900 dark:bg-zinc-100 text-white dark:text-zinc-900
            hover:opacity-85 transition-opacity"
        >
          <Plus size={14} />
          New project
        </Link>
      </div>

      {/* Error */}
      {error && (
        <p className="text-xs text-red-500 bg-red-500/10 border border-red-500/20 rounded-lg px-3 py-2 mb-6">
          Failed to load projects.
        </p>
      )}

      {/* Loading */}
      {isLoading && (
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
          {[...Array(4)].map((_, i) => (
            <Skeleton key={i} className="h-40" />
          ))}
        </div>
      )}

      {/* Empty */}
      {!isLoading && !error && projects?.length === 0 && (
        <EmptyState
          icon={Code2}
          title="No projects yet"
          description="Import a GitHub repo or start from scratch — Smarty will turn it into a learning path."
          action={{
            label: "Create your first project",
            href: "/dashboard/projects/new",
          }}
        />
      )}

      {/* Project grid */}
      {!isLoading && projects && projects.length > 0 && (
        <>
          {active.length > 0 && (
            <div className="mb-8">
              <SectionHeader label="Active" />
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                {active.map((p: Project) => (
                  <ProjectCard key={p.id} project={p} />
                ))}
              </div>
            </div>
          )}
          {completed.length > 0 && (
            <div>
              <SectionHeader label="Completed" />
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                {completed.map((p: Project) => (
                  <ProjectCard key={p.id} project={p} />
                ))}
              </div>
            </div>
          )}
        </>
      )}
    </div>
  );
}
