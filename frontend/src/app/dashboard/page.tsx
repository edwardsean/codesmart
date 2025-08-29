"use client";

import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import Button from "@/components/ui/Button";
import Link from "next/link";
import {
  Upload,
  GitBranch,
  Users,
  Settings,
  Plus,
  Activity,
  FileCode,
  AlertCircle,
} from "lucide-react";

interface Project {
  id: string;
  name: string;
  description: string;

  repository: {
    owner: string;
    name: string;
    url: string;
    branch: string;
  };

  stats: {
    totalReviews: number;
    activeIssues: number;
    teamMembers: number;
    lastActivity: string;
  };
  language: string;
  status: "active" | "archived";
}

export default function Dashboard() {
  const [projects, setProjects] = useState<Project[]>([
    {
      id: "1",
      name: "E-commerce Platform",
      description: "Full-stack e-commerce application with React and Node.js",
      repository: {
        owner: "mycompany",
        name: "ecommerce-platform",
        url: "https://github.com/mycompany/ecommerce-platform",
        branch: "main",
      },
      stats: {
        totalReviews: 45,
        activeIssues: 3,
        teamMembers: 5,
        lastActivity: "2 hours ago",
      },
      language: "JavaScript",
      status: "active",
    },
  ]);

  return (
    <>
      {true && (
        <div className="p-8">
          {/* Welcome Section */}
          <div className="mb-8">
            <h1 className="text-3xl font-bold text-gray-900">
              Welcome back username!
            </h1>
            <p className="text-gray-600 mt-2">
              Ready to review some code? Lets get started.
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
            {projects.map((project) => (
              <ProjectCard key={project.id} project={project}></ProjectCard>
            ))}
          </div>
        </div>
      )}
    </>
  );
}

interface ProjectCardProps {
  project: Project;
}

function ProjectCard({ project }: ProjectCardProps) {
  const router = useRouter();

  const getLanguageColor = (language: string) => {
    const colors: { [key: string]: string } = {
      JavaScript: "bg-yellow-400",
      TypeScript: "bg-blue-500",
      Python: "bg-green-500",
      Java: "bg-orange-500",
      Go: "bg-cyan-500",
    };

    return colors[language] || "bg-gray-400";
  };

  const handleCardClick = () => {
    router.push(`/dashboard/projects/${project.name}`);
  };

  return (
    <div
      className="bg-white rounded-lg shadow border border-gray-200 p-6 hover:shadow-lg transition-shadow cursor-pointer"
      onClick={handleCardClick}
    >
      {/* Header */}
      <div className="flex justify-between items-start mb-4">
        <div>
          <h3 className="text-lg font-semibold text-gray-900 mb-1">
            {project.name}
          </h3>
          <p className="text-sm text-gray-500">
            {project.repository.owner}/{project.repository.name}
          </p>
        </div>
        {/* <button
          onClick={handleExternalClick}
          className="text-gray-400 hover:text-gray-600 transition-colors"
          title="Open repository in GitHub"
        >
          <ExternalLink className="w-5 h-5" />
        </button> */}
      </div>

      {/* Description */}
      <p className="text-gray-600 text-sm mb-4 line-clamp-2">
        {project.description}
      </p>

      {/* Language & Status */}
      <div className="flex items-center gap-3 mb-4">
        <div className="flex items-center gap-2">
          <div
            className={`w-3 h-3 rounded-full ${getLanguageColor(
              project.language
            )}`}
          ></div>
          <span className="text-sm text-gray-600">{project.language}</span>
        </div>
        <span
          className={`px-2 py-1 rounded-full text-xs font-medium ${
            project.status === "active"
              ? "bg-green-100 text-green-800"
              : "bg-gray-100 text-gray-800"
          }`}
        >
          {project.status}
        </span>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-2 gap-4 mb-4">
        <div className="text-center">
          <p className="text-lg font-semibold text-gray-900">
            {project.stats.totalReviews}
          </p>
          <p className="text-xs text-gray-500">Reviews</p>
        </div>
        <div className="text-center">
          <p className="text-lg font-semibold text-gray-900">
            {project.stats.teamMembers}
          </p>
          <p className="text-xs text-gray-500">Members</p>
        </div>
      </div>

      {/* Footer */}
      <div className="flex justify-between items-center pt-4 border-t border-gray-100">
        <div className="flex items-center gap-2">
          <AlertCircle className="w-4 h-4 text-yellow-500" />
          <span className="text-sm text-gray-600">
            {project.stats.activeIssues} active issues
          </span>
        </div>
        <span className="text-xs text-gray-500">
          {project.stats.lastActivity}
        </span>
      </div>
    </div>
  );
}
