"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import {
  Files,
  LayoutDashboard,
  Network,
  Database,
  ChevronsLeft,
  ChevronsRight,
  ChevronLeft,
} from "lucide-react";
import ThemeToggle from "@/components/ui/ThemeToggle";
import ResizablePanel from "@/components/ide/ResizablePanel";
import FileTree from "@/components/ide/FileTree";
import { IDEProvider, useIDE } from "@/context/IDEProvider";
import { useAuthStore } from "@/stores/authStore";
import { NAV_LINKS } from "@/lib/constants/navigation";
import { FileNode } from "@/types/entity";
import { fileService } from "@/services/fileService";
import useAxiosPrivate from "@/hooks/useAxiosPrivate";
import axios from "axios";

type ActivityPanel = "files" | "nav" | null;

interface IDELayoutProps {
  children: React.ReactNode;
  projectId: string;
  projectTitle?: string;
}

function IDEShell({ children, projectId, projectTitle }: IDELayoutProps) {
  const pathname = usePathname();
  const { account } = useAuthStore();
  const {
    files,
    filesLoading,
    isGithub,
    activeFile,
    setActiveFile,
    setFileContent,
    setFiles,
  } = useIDE();

  console.log("files in shell: ", files);
  const api = useAxiosPrivate();
  const apiFile = fileService(api);

  const [activePanel, setActivePanel] = useState<ActivityPanel>("files");
  const [panelCollapsed, setPanelCollapsed] = useState(false);

  const isDesign = pathname.includes("/design");
  const isDatabase = pathname.includes("/database");
  const isCode = !isDesign && !isDatabase;

  const userInitial = account?.user.username?.[0]?.toUpperCase() ?? "U";

  function togglePanel(panel: ActivityPanel) {
    if (activePanel === panel) {
      setPanelCollapsed((c) => !c);
    } else {
      setActivePanel(panel);
      setPanelCollapsed(false);
    }
  }

  async function handleFileSelect(node: FileNode) {
    if (node.type === "dir") return;
    setActiveFile(node.path);
    try {
      const content = await apiFile.getFileContent(
        Number(projectId),
        node.path,
      );

      console.log("file content: ", content);
      setFileContent(content.content);
    } catch {
      setFileContent("// Failed to load file content");
    }
  }

  async function handleCreateFile(parentPath: string, isDir: boolean) {
    const name = prompt(isDir ? "Folder name:" : "File name:");
    if (!name) return;
    const path = parentPath ? `${parentPath}/${name}` : name;
    try {
      await apiFile.createFile(Number(projectId), {
        path,
        content: "",
        is_dir: isDir,
      });
      const tree = await apiFile.getFileTree(Number(projectId));
      setFiles(tree);
    } catch (err) {
      if (axios.isAxiosError(err)) console.error(err.response?.data?.message);
    }
  }

  async function handleRename(node: FileNode, newName: string) {
    if (!node.id) return;
    const parts = node.path.split("/");
    parts[parts.length - 1] = newName;
    const newPath = parts.join("/");
    try {
      await apiFile.updateFile(Number(projectId), node.id, { path: newPath });
      const tree = await apiFile.getFileTree(Number(projectId));
      setFiles(tree);
      if (activeFile === node.path) setActiveFile(newPath);
    } catch (err) {
      if (axios.isAxiosError(err)) console.error(err.response?.data?.message);
    }
  }

  async function handleDelete(node: FileNode) {
    if (!node.id) return;
    if (!confirm(`Delete ${node.name}?`)) return;
    try {
      await apiFile.deleteFile(Number(projectId), node.id);
      const tree = await apiFile.getFileTree(Number(projectId));
      setFiles(tree);
      if (activeFile === node.path) {
        setActiveFile(undefined);
        setFileContent("");
      }
    } catch (err) {
      if (axios.isAxiosError(err)) console.error(err.response?.data?.message);
    }
  }

  return (
    <div className="flex h-screen bg-white dark:bg-[#1e1e1e] overflow-hidden">
      {/* ── Activity bar ── */}
      <div
        className="w-12 flex-shrink-0 border-r border-black/[0.06] dark:border-white/[0.06]
        flex flex-col items-center py-2 gap-0.5 bg-gray-50 dark:bg-[#252526]"
      >
        <Link
          href="/dashboard"
          className="w-8 h-8 rounded-md border border-[#dc503c]/70 flex items-center justify-center mb-2 flex-shrink-0"
        >
          <svg width="12" height="12" viewBox="0 0 14 14" fill="none">
            <path
              d="M2 4.5L5.5 7L2 9.5"
              stroke="#dc503c"
              strokeWidth="1.5"
              strokeLinecap="round"
              strokeLinejoin="round"
            />
            <path
              d="M7.5 9.5H12"
              stroke="#dc503c"
              strokeWidth="1.5"
              strokeLinecap="round"
            />
          </svg>
        </Link>

        <ActivityButton
          icon={Files}
          active={activePanel === "files" && !panelCollapsed}
          onClick={() => togglePanel("files")}
          tooltip="Explorer"
        />
        <ActivityButton
          icon={LayoutDashboard}
          active={activePanel === "nav" && !panelCollapsed}
          onClick={() => togglePanel("nav")}
          tooltip="Navigation"
        />

        <div className="flex-1" />

        <Link href={`/dashboard/projects/${projectId}/design`}>
          <ActivityButton
            icon={Network}
            active={isDesign}
            onClick={() => {}}
            tooltip="System design"
          />
        </Link>
        <Link href={`/dashboard/projects/${projectId}/database`}>
          <ActivityButton
            icon={Database}
            active={isDatabase}
            onClick={() => {}}
            tooltip="Database design"
          />
        </Link>

        <div className="mt-2 mb-1">
          <ThemeToggle />
        </div>

        <div
          className="w-7 h-7 rounded-full bg-[#dc503c]/20 flex items-center justify-center
          text-xs font-mono text-[#dc503c] mb-1"
        >
          {userInitial}
        </div>
      </div>

      {/* ── Side panel ── */}
      {activePanel && !panelCollapsed && (
        <ResizablePanel
          defaultWidth={220}
          minWidth={160}
          maxWidth={480}
          collapsed={panelCollapsed}
          onCollapse={setPanelCollapsed}
        >
          {activePanel === "files" && (
            <FileTree
              loading={filesLoading}
              files={files}
              activeFile={activeFile}
              projectTitle={projectTitle ?? "Project"}
              isGithub={isGithub}
              onFileSelect={handleFileSelect}
              onCreateFile={handleCreateFile}
              onRename={handleRename}
              onDelete={handleDelete}
            />
          )}

          {activePanel === "nav" && (
            <div
              className="h-full flex flex-col bg-gray-50 dark:bg-zinc-900/60
              border-r border-black/[0.06] dark:border-white/[0.06]"
            >
              <div className="px-3 py-2 border-b border-black/[0.04] dark:border-white/[0.04]">
                <p className="text-xs font-mono text-gray-400 dark:text-zinc-500 uppercase tracking-widest">
                  Navigation
                </p>
              </div>
              <div className="flex-1 py-3 px-2 space-y-0.5">
                {NAV_LINKS.map(({ name, href, icon: Icon }) => (
                  <Link
                    key={href}
                    href={href}
                    className="flex items-center gap-2.5 px-3 py-2 rounded-lg text-sm
                      text-gray-500 dark:text-zinc-400
                      hover:bg-black/[0.04] dark:hover:bg-white/[0.04]
                      hover:text-gray-900 dark:hover:text-zinc-100 transition-colors"
                  >
                    <Icon size={14} />
                    {name}
                  </Link>
                ))}
                <div className="pt-3 mt-2 border-t border-black/[0.06] dark:border-white/[0.06]">
                  <Link
                    href="/dashboard/projects"
                    className="flex items-center gap-2 text-xs text-gray-400 dark:text-zinc-500
                      hover:text-[#dc503c] transition-colors px-3 py-1.5"
                  >
                    <ChevronLeft size={12} />
                    Back to projects
                  </Link>
                </div>
              </div>
            </div>
          )}
        </ResizablePanel>
      )}

      {/* Collapsed strip */}
      {activePanel && panelCollapsed && (
        <div className="w-px border-r border-black/[0.06] dark:border-white/[0.06] relative group">
          <button
            onClick={() => setPanelCollapsed(false)}
            className="absolute top-1/2 -translate-y-1/2 -right-3 w-6 h-6 rounded-full
              bg-gray-200 dark:bg-zinc-700 flex items-center justify-center
              opacity-0 group-hover:opacity-100 transition-opacity z-10"
          >
            <ChevronsRight
              size={12}
              className="text-gray-600 dark:text-zinc-300"
            />
          </button>
        </div>
      )}

      {/* ── Main content ── */}
      <div className="flex-1 flex flex-col min-w-0 overflow-hidden">
        {/* Top bar */}
        <div
          className="h-9 flex-shrink-0 border-b border-black/[0.06] dark:border-white/[0.06]
          flex items-center px-4 gap-2 bg-gray-50/80 dark:bg-[#252526]/80 backdrop-blur-sm"
        >
          <Link
            href="/dashboard/projects"
            className="flex items-center gap-1 text-xs text-gray-400 dark:text-zinc-500
              hover:text-gray-700 dark:hover:text-zinc-300 transition-colors flex-shrink-0"
          >
            <ChevronLeft size={12} />
            Projects
          </Link>
          <span className="text-gray-200 dark:text-zinc-700 text-xs">/</span>
          <span className="text-xs text-gray-600 dark:text-zinc-400 truncate">
            {projectTitle ?? `project ${projectId}`}
          </span>

          {activePanel && (
            <button
              onClick={() => setPanelCollapsed((c) => !c)}
              className="text-gray-400 dark:text-zinc-600 hover:text-gray-700 dark:hover:text-zinc-300 transition-colors flex-shrink-0"
              title={panelCollapsed ? "Show panel" : "Hide panel"}
            >
              {panelCollapsed ? (
                <ChevronsRight size={13} />
              ) : (
                <ChevronsLeft size={13} />
              )}
            </button>
          )}

          <div className="flex-1" />

          {/* Mode tabs */}
          <div className="flex items-center gap-0.5 bg-gray-100 dark:bg-zinc-800 p-0.5 rounded-lg">
            {[
              {
                href: `/dashboard/projects/${projectId}`,
                label: "Code",
                icon: null,
                active: isCode,
              },
              {
                href: `/dashboard/projects/${projectId}/design`,
                label: "Design",
                icon: Network,
                active: isDesign,
              },
              {
                href: `/dashboard/projects/${projectId}/database`,
                label: "Database",
                icon: Database,
                active: isDatabase,
              },
            ].map(({ href, label, icon: Icon, active }) => (
              <Link
                key={href}
                href={href}
                className={`text-xs px-3 py-1 rounded-md transition-colors flex items-center gap-1 ${
                  active
                    ? "bg-white dark:bg-zinc-700 text-gray-900 dark:text-zinc-100 shadow-sm"
                    : "text-gray-500 dark:text-zinc-400 hover:text-gray-700 dark:hover:text-zinc-300"
                }`}
              >
                {Icon && <Icon size={10} />}
                {label}
              </Link>
            ))}
          </div>
        </div>

        {/* Page content */}
        <div className="flex-1 overflow-hidden">{children}</div>
      </div>
    </div>
  );
}

//provides context
export default function IDELayout({
  children,
  projectId,
  projectTitle,
}: IDELayoutProps) {
  const api = useAxiosPrivate();
  const apiFile = fileService(api);
  const { _hasHydrated, account } = useAuthStore();

  const [files, setFiles] = useState<FileNode[]>([]);
  const [filesLoading, setFilesLoading] = useState(true);
  const [isGithub, setIsGithub] = useState(false);

  useEffect(() => {
    if (!_hasHydrated || !account?.accessToken) return;

    const fetchFiles = async () => {
      setFilesLoading(true);
      try {
        const tree = await apiFile.getFileTree(Number(projectId));
        setFiles(tree);
      } catch (err) {
        if (axios.isAxiosError(err)) {
          console.error(
            "failed to fetch file tree:",
            err.response?.data?.message,
          );
        }
      } finally {
        setFilesLoading(false);
      }
    };

    fetchFiles();
  }, [_hasHydrated, account?.accessToken, projectId]);

  return (
    <IDEProvider
      files={files}
      filesLoading={filesLoading}
      isGithub={isGithub}
      onSetFiles={setFiles}
    >
      <IDEShell projectId={projectId} projectTitle={projectTitle}>
        {children}
      </IDEShell>
    </IDEProvider>
  );
}

function ActivityButton({
  icon: Icon,
  active,
  onClick,
  tooltip,
}: {
  icon: React.ElementType;
  active: boolean;
  onClick: () => void;
  tooltip: string;
}) {
  return (
    <button
      onClick={onClick}
      title={tooltip}
      className={`w-9 h-9 rounded-lg flex items-center justify-center transition-colors ${
        active
          ? "bg-gray-200 dark:bg-zinc-700 text-gray-900 dark:text-zinc-100"
          : "text-gray-400 dark:text-zinc-500 hover:bg-gray-100 dark:hover:bg-[#2d2d2d] hover:text-gray-700 dark:hover:text-zinc-300"
      }`}
    >
      <Icon size={16} />
    </button>
  );
}
