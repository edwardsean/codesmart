"use client";

import { createContext, useContext, useState, ReactNode } from "react";
import { FileNode } from "@/types/project.file.types";

export interface ActiveFile {
  id: number;
  path: string;
}

interface IDEContextValue {
  //project id
  projectId: string;
  // file tree
  files: FileNode[];
  filesLoading: boolean;
  isGithub: boolean;

  // active file
  activeFile: ActiveFile | null;

  //openTabs
  openTabs: ActiveFile[];
  openFile: (activeFile: ActiveFile) => void;
  closeTab: (path: string) => void;
  editedContent: Record<string, string>;
  setEditedContent: (path: string, content: string) => void;
  clearEditedContent: (path: string) => void;

  // actions
  setActiveFile: (activeFile: ActiveFile | null) => void;
  setFiles: (files: FileNode[]) => void;
}

const IDEContext = createContext<IDEContextValue | null>(null);

export function IDEProvider({
  children,
  projectId,
  files,
  filesLoading,
  isGithub,
  onSetFiles,
}: {
  children: ReactNode;
  projectId: string;
  files: FileNode[];
  filesLoading: boolean;
  isGithub: boolean;
  onSetFiles: (files: FileNode[]) => void;
}) {
  const [activeFile, setActiveFile] = useState<ActiveFile | null>(null); //takes the path
  const [openTabs, setOpenTabs] = useState<ActiveFile[]>([]);
  const [editedContent, setEditedContentMap] = useState<Record<string, string>>(
    {},
  );

  function openFile(activeFile: ActiveFile) {
    const isAlreadyOpen = openTabs.some((tab) => tab.path === activeFile.path);

    if (!isAlreadyOpen) {
      setOpenTabs((prev) => [
        ...prev,
        { path: activeFile.path, id: activeFile.id },
      ]);
    }
    setActiveFile({ path: activeFile.path, id: activeFile.id });
  }

  function closeTab(path: string) {
    setOpenTabs((prev) => {
      const tabToClose = prev.find((p) => p.path === path);
      if (!tabToClose) return prev;

      const next = prev.filter((tab) => tab.path !== path);

      if (activeFile && path === activeFile.path) {
        const idx = prev.findIndex((tab) => tab.path === path);
        const nextActive = next[idx] ?? next[idx - 1] ?? null;
        setActiveFile(nextActive);
      }

      return next;
    });
  }

  function setEditedContent(path: string, content: string) {
    setEditedContentMap((prev) => ({ ...prev, [path]: content }));
  }

  function clearEditedContent(path: string) {
    setEditedContentMap((prev) => {
      const next = { ...prev };
      delete next[path];
      return next;
    });
  }

  return (
    <IDEContext.Provider
      value={{
        projectId,
        files,
        filesLoading,
        isGithub,
        activeFile,
        editedContent,
        openTabs,
        openFile,
        closeTab,
        setEditedContent,
        clearEditedContent,
        setActiveFile,
        setFiles: onSetFiles,
      }}
    >
      {children}
    </IDEContext.Provider>
  );
}

export function useIDE() {
  const ctx = useContext(IDEContext);
  if (!ctx) throw new Error("useIDE must be used within IDEProvider");
  return ctx;
}
