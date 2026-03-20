"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { Plus, Code2 } from "lucide-react";
import SectionHeader from "@/components/ui/SectionHeader";
import EmptyState from "@/components/ui/EmptyState";
import ProjectCard from "@/components/dashboard/ProjectCard";
import { Skeleton } from "@/components/ui/Skeleton";
import useAxiosPrivate from "@/hooks/useAxiosPrivate";
import { projectService } from "@/services/projectService";
import { ProjectListItem } from "@/types/entity";
import { useAuthStore } from "@/stores/authStore";

export default function ProjectsPage() {
  const api = useAxiosPrivate();
  const apiProject = projectService(api);
  const { _hasHydrated, account } = useAuthStore();

  const [projects, setProjects] = useState<ProjectListItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    if (!_hasHydrated || !account?.accessToken) return;

    const fetchProjects = async () => {
      try {
        const { projects } = await apiProject.getProjects();
        setProjects(projects);
      } catch {
        setError("Failed to load projects.");
      } finally {
        setLoading(false);
      }
    };

    fetchProjects();
  }, [_hasHydrated, account?.accessToken]);

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
          {error}
        </p>
      )}

      {/* Loading */}
      {loading && (
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
          {[...Array(4)].map((_, i) => (
            <Skeleton key={i} className="h-40" />
          ))}
        </div>
      )}

      {/* Empty */}
      {!loading && !error && projects.length === 0 && (
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
      {!loading && projects.length > 0 && (
        <>
          {/* Active */}
          {projects.filter((p) => p.status === "active").length > 0 && (
            <div className="mb-8">
              <SectionHeader label="Active" />
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                {projects
                  .filter((p) => p.status === "active")
                  .map((p) => (
                    <ProjectCard key={p.id} project={p} />
                  ))}
              </div>
            </div>
          )}

          {/* Completed */}
          {projects.filter((p) => p.status === "completed").length > 0 && (
            <div>
              <SectionHeader label="Completed" />
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                {projects
                  .filter((p) => p.status === "completed")
                  .map((p) => (
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
