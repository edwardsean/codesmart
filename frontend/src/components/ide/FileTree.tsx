"use client";

import { useState, useRef } from "react";
import {
  ChevronRight,
  ChevronDown,
  File,
  Folder,
  FolderOpen,
  FilePlus,
  FolderPlus,
  Pencil,
  Trash2,
} from "lucide-react";
import { FileNode } from "@/types/project.file.types";
import { usePrefetchFileContent } from "@/hooks/query/useFileContent";

interface FileTreeProps {
  projectId: string;
  files: FileNode[];
  activeFile?: string | null;
  projectTitle: string;
  isGithub?: boolean;
  loading: boolean;
  onFileSelect: (node: FileNode) => void;
  onCreateFile?: (parentPath: string, isDir: boolean) => void;
  onRename?: (node: FileNode, newName: string) => void;
  onDelete?: (node: FileNode) => void;
}

const LANG_COLORS: Record<string, string> = {
  go: "text-blue-400",
  ts: "text-blue-500",
  tsx: "text-blue-400",
  js: "text-yellow-400",
  jsx: "text-yellow-400",
  py: "text-green-400",
  sql: "text-orange-400",
  md: "text-gray-400",
  css: "text-pink-400",
  json: "text-yellow-300",
  yaml: "text-red-400",
  yml: "text-red-400",
};

function getExt(name: string) {
  return name.split(".").pop() ?? "";
}

function FileTreeNode({
  projectId,
  node,
  depth,
  activeFile,
  isGithub,
  onFileSelect,
  onCreateFile,
  onRename,
  onDelete,
}: {
  projectId: string;
  node: FileNode;
  depth: number;
  activeFile?: string | null;
  isGithub?: boolean;
  onFileSelect: (n: FileNode) => void;
  onCreateFile?: (parentPath: string, isDir: boolean) => void;
  onRename?: (node: FileNode, newName: string) => void;
  onDelete?: (node: FileNode) => void;
}) {
  const [open, setOpen] = useState(depth < 1);
  const [showMenu, setShowMenu] = useState(false);
  const [isRenaming, setIsRenaming] = useState(false);
  const [renameValue, setRenameValue] = useState(node.name);
  const renameRef = useRef<HTMLInputElement>(null);
  const isActive = activeFile && node.path === activeFile;
  const ext = getExt(node.name);
  const color = LANG_COLORS[ext] ?? "text-gray-400 dark:text-zinc-500";
  const prefetch = usePrefetchFileContent(projectId);

  function handleRenameSubmit() {
    if (renameValue.trim() && renameValue !== node.name) {
      onRename?.(node, renameValue.trim());
    }
    setIsRenaming(false);
  }

  const indent = { paddingLeft: `${8 + depth * 12}px` };

  if (node.type === "dir") {
    return (
      <div>
        <div
          className="group flex items-center gap-1 py-0.5 hover:bg-black/[0.04] dark:hover:bg-white/[0.04]
            text-gray-600 dark:text-zinc-400 rounded cursor-pointer relative"
          style={indent}
          onClick={() => setOpen((o) => !o)}
        >
          {open ? (
            <ChevronDown size={11} className="flex-shrink-0 text-gray-400" />
          ) : (
            <ChevronRight size={11} className="flex-shrink-0 text-gray-400" />
          )}
          {open ? (
            <FolderOpen size={13} className="flex-shrink-0 text-[#dc503c]/70" />
          ) : (
            <Folder size={13} className="flex-shrink-0 text-[#dc503c]/70" />
          )}

          {isRenaming ? (
            <input
              ref={renameRef}
              value={renameValue}
              onChange={(e) => setRenameValue(e.target.value)}
              onBlur={handleRenameSubmit}
              onKeyDown={(e) => {
                if (e.key === "Enter") handleRenameSubmit();
                if (e.key === "Escape") setIsRenaming(false);
              }}
              className="text-xs bg-white dark:bg-zinc-800 border border-[#dc503c]/50 rounded px-1 outline-none flex-1"
              autoFocus
              onClick={(e) => e.stopPropagation()}
            />
          ) : (
            <span className="text-xs truncate flex-1">{node.name}</span>
          )}

          {/* Context menu button — scratch only */}
          {!isGithub && !isRenaming && (
            <div
              className="opacity-0 group-hover:opacity-100 flex items-center gap-0.5 mr-1 flex-shrink-0"
              onClick={(e) => e.stopPropagation()}
            >
              <button
                onClick={() => onCreateFile?.(node.path, false)}
                className="p-0.5 hover:text-gray-900 dark:hover:text-zinc-100 transition-colors"
                title="New file"
              >
                <FilePlus size={11} />
              </button>
              <button
                onClick={() => onCreateFile?.(node.path, true)}
                className="p-0.5 hover:text-gray-900 dark:hover:text-zinc-100 transition-colors"
                title="New folder"
              >
                <FolderPlus size={11} />
              </button>
              <button
                onClick={() => {
                  setIsRenaming(true);
                  setRenameValue(node.name);
                }}
                className="p-0.5 hover:text-gray-900 dark:hover:text-zinc-100 transition-colors"
                title="Rename"
              >
                <Pencil size={11} />
              </button>
              <button
                onClick={() => onDelete?.(node)}
                className="p-0.5 hover:text-red-500 transition-colors"
                title="Delete"
              >
                <Trash2 size={11} />
              </button>
            </div>
          )}
        </div>

        {open &&
          node.children?.map((child) => (
            <FileTreeNode
              projectId={projectId}
              key={child.path}
              node={child}
              depth={depth + 1}
              activeFile={activeFile}
              isGithub={isGithub}
              onFileSelect={onFileSelect}
              onCreateFile={onCreateFile}
              onRename={onRename}
              onDelete={onDelete}
            />
          ))}
      </div>
    );
  }

  // File node
  return (
    <div
      className={`group flex items-center gap-1.5 py-0.5 rounded cursor-pointer transition-colors relative ${
        isActive
          ? "bg-[#dc503c]/10 text-gray-900 dark:text-zinc-100"
          : "text-gray-500 dark:text-zinc-400 hover:bg-black/[0.04] dark:hover:bg-white/[0.04]"
      }`}
      style={indent}
      onMouseEnter={() => {
        if (node.type === "file") prefetch(node.path);
      }}
      onClick={() => onFileSelect(node)}
      role="button"
    >
      <File size={12} className={`flex-shrink-0 ${color}`} />

      {isRenaming ? (
        <input
          value={renameValue}
          onChange={(e) => setRenameValue(e.target.value)}
          onBlur={handleRenameSubmit}
          onKeyDown={(e) => {
            if (e.key === "Enter") handleRenameSubmit();
            if (e.key === "Escape") setIsRenaming(false);
          }}
          className="text-xs bg-white dark:bg-zinc-800 border border-[#dc503c]/50 rounded px-1 outline-none flex-1"
          autoFocus
          onClick={(e) => e.stopPropagation()}
        />
      ) : (
        <span className="text-xs truncate flex-1">{node.name}</span>
      )}

      {!isGithub && !isRenaming && (
        <div
          className="opacity-0 group-hover:opacity-100 flex items-center gap-0.5 mr-1 flex-shrink-0"
          onClick={(e) => e.stopPropagation()}
        >
          <button
            onClick={() => {
              setIsRenaming(true);
              setRenameValue(node.name);
            }}
            className="p-0.5 hover:text-gray-900 dark:hover:text-zinc-100 transition-colors"
            title="Rename"
          >
            <Pencil size={11} />
          </button>
          <button
            onClick={() => onDelete?.(node)}
            className="p-0.5 hover:text-red-500 transition-colors"
            title="Delete"
          >
            <Trash2 size={11} />
          </button>
        </div>
      )}
    </div>
  );
}

export default function FileTree({
  projectId,
  loading,
  files,
  activeFile,
  projectTitle,
  isGithub = false,
  onFileSelect,
  onCreateFile,
  onRename,
  onDelete,
}: FileTreeProps) {
  return (
    <div className="h-full flex flex-col bg-gray-50 dark:bg-zinc-900/60 border-r border-black/[0.06] dark:border-white/[0.06]">
      {/* Header */}
      <div className="px-3 py-2 border-b border-black/[0.04] dark:border-white/[0.04] flex items-center justify-between flex-shrink-0">
        <p className="text-xs font-mono text-gray-400 dark:text-zinc-500 uppercase tracking-widest truncate">
          {projectTitle}
        </p>
        {/* Root-level create buttons for scratch projects */}
        {!isGithub && (
          <div className="flex items-center gap-1 flex-shrink-0">
            <button
              onClick={() => onCreateFile?.("", false)}
              className="text-gray-400 dark:text-zinc-500 hover:text-gray-700 dark:hover:text-zinc-300 transition-colors"
              title="New file"
            >
              <FilePlus size={13} />
            </button>
            <button
              onClick={() => onCreateFile?.("", true)}
              className="text-gray-400 dark:text-zinc-500 hover:text-gray-700 dark:hover:text-zinc-300 transition-colors"
              title="New folder"
            >
              <FolderPlus size={13} />
            </button>
          </div>
        )}
      </div>

      {/* Tree */}
      {/* Tree */}
      <div className="flex-1 overflow-y-auto py-1 px-1">
        {loading ? (
          <FileTreeSkeleton />
        ) : files.length === 0 ? (
          <div className="px-4 py-8 text-center">
            <p className="text-xs text-gray-400 dark:text-zinc-600">
              {isGithub
                ? "No files found."
                : "No files yet. Create one to get started."}
            </p>
          </div>
        ) : (
          files.map((node) => (
            <FileTreeNode
              projectId={projectId}
              key={node.path}
              node={node}
              depth={0}
              activeFile={activeFile}
              isGithub={isGithub}
              onFileSelect={onFileSelect}
              onCreateFile={onCreateFile}
              onRename={onRename}
              onDelete={onDelete}
            />
          ))
        )}
      </div>
    </div>
  );
}

function FileTreeSkeleton() {
  const items = [
    { depth: 0, width: "w-24", isDir: true },
    { depth: 1, width: "w-20" },
    { depth: 1, width: "w-16" },
    { depth: 1, width: "w-28", isDir: true },
    { depth: 2, width: "w-14" },
    { depth: 2, width: "w-20" },
    { depth: 0, width: "w-16", isDir: true },
    { depth: 1, width: "w-24" },
    { depth: 0, width: "w-20" },
    { depth: 0, width: "w-12" },
  ];

  return (
    <div className="py-1 px-1 space-y-0.5">
      {items.map((item, i) => (
        <div
          key={i}
          className="flex items-center gap-1.5 py-0.5"
          style={{ paddingLeft: `${8 + item.depth * 12}px` }}
        >
          {/* chevron placeholder */}
          <div className="w-2.5 h-2.5 rounded bg-gray-200 dark:bg-zinc-700 animate-pulse flex-shrink-0" />
          {/* icon placeholder */}
          <div className="w-3 h-3 rounded bg-gray-200 dark:bg-zinc-700 animate-pulse flex-shrink-0" />
          {/* name placeholder */}
          <div
            className={`h-2.5 rounded bg-gray-200 dark:bg-zinc-700 animate-pulse ${item.width}`}
          />
        </div>
      ))}
    </div>
  );
}
