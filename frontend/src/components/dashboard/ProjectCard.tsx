import Link from "next/link";
import { GitBranch, Code2, Clock } from "lucide-react";
import { ProjectListItem } from "@/types/entity";

interface ProjectCardProps {
  project: ProjectListItem;
}

const MODE_LABELS = {
  help: "Help mode",
  learn: "Learn mode",
} as const;

const STATUS_STYLES = {
  active: "bg-blue-500/10 text-blue-600 dark:text-blue-400",
  completed: "bg-green-500/10 text-green-600 dark:text-green-400",
  archived: "bg-gray-500/10 text-gray-500",
} as const;

export default function ProjectCard({ project }: ProjectCardProps) {
  return (
    <Link
      href={`/dashboard/projects/${project.id}`}
      className="group block p-5 rounded-xl border border-black/[0.06] dark:border-white/[0.06]
        bg-gray-50/50 dark:bg-zinc-900/50
        hover:border-black/10 dark:hover:border-white/10
        hover:bg-gray-100/50 dark:hover:bg-zinc-800/50
        transition-all"
    >
      {/* Header */}
      <div className="flex items-start justify-between gap-3 mb-4">
        <div className="min-w-0">
          <p className="text-sm font-medium text-gray-900 dark:text-zinc-100 truncate">
            {project.title}
          </p>
          <div className="flex items-center gap-2 mt-1">
            <span className="text-xs font-mono text-gray-400 dark:text-zinc-500">
              {project.language}
            </span>
            {project.source_type === "github" && (
              <GitBranch
                size={11}
                className="text-gray-300 dark:text-zinc-600"
              />
            )}
          </div>
        </div>
        <span
          className={`text-xs px-2 py-0.5 rounded-md flex-shrink-0 ${STATUS_STYLES[project.status]}`}
        >
          {project.status}
        </span>
      </div>

      {/* Progress bar */}
      {project.total_levels > 0 && (
        <div className="mb-4">
          <div className="flex justify-between text-xs text-gray-400 dark:text-zinc-600 mb-1.5">
            <span>
              {project.completed_levels} / {project.total_levels} levels
            </span>
            <span>{project.progress_percentage}%</span>
          </div>
          <div className="h-1 bg-gray-200 dark:bg-zinc-700 rounded-full overflow-hidden">
            <div
              className="h-full bg-[#dc503c] rounded-full transition-all"
              style={{ width: `${project.progress_percentage}%` }}
            />
          </div>
        </div>
      )}

      {/* Footer */}
      <div className="flex items-center justify-between">
        <span className="text-xs text-gray-400 dark:text-zinc-600 flex items-center gap-1.5">
          <Code2 size={11} />
          {MODE_LABELS[project.mode]}
        </span>
        <span className="text-xs text-gray-400 dark:text-zinc-600 flex items-center gap-1.5">
          <Clock size={11} />
          {new Date(project.updated_at).toLocaleDateString("en-US", {
            month: "short",
            day: "numeric",
          })}
        </span>
      </div>
    </Link>
  );
}
