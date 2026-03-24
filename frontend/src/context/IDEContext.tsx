"use client";

import { createContext, useContext, useState, ReactNode } from "react";
import { FileNode } from "@/types/project.file.types";

interface IDEContextValue {
  //project id
  projectId: string;
  // file tree
  files: FileNode[];
  filesLoading: boolean;
  isGithub: boolean;

  // active file
  activeFile: string | null;

  //openTabs
  openTabs: string[];
  openFile: (path: string) => void;
  closeTab: (path: string) => void;
  editedContent: Record<string, string>;
  setEditedContent: (path: string, content: string) => void;
  clearEditedContent: (path: string) => void;

  // actions
  setActiveFile: (path: string) => void;
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
  const [activeFile, setActiveFile] = useState<string | null>(null); //takes the path
  const [openTabs, setOpenTabs] = useState<string[]>([]);
  const [editedContent, setEditedContentMap] = useState<Record<string, string>>(
    {},
  );

  function openFile(path: string) {
    const isAlreadyOpen = openTabs.some((p) => p === path);

    if (!isAlreadyOpen) {
      setOpenTabs((prev) => [...prev, path]);
    }
    setActiveFile(path);
  }

  function closeTab(path: string) {
    setOpenTabs((prev) => {
      const tabToClose = prev.find((p) => p === path);
      if (!tabToClose) return prev;

      const next = prev.filter((p) => p !== path);

      if (path === activeFile) {
        const idx = prev.findIndex((p) => p === path);
        const nextActive = next[idx] ?? next[idx - 1];
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
