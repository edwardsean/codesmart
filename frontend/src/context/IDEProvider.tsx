"use client";

import { createContext, useContext, useState, ReactNode } from "react";
import { FileNode, FileContent } from "@/types/entity";

interface IDEContextValue {
  // file tree
  files: FileNode[];
  filesLoading: boolean;
  isGithub: boolean;

  // active file
  activeFile: string | undefined;
  fileContent: string;

  // actions
  setActiveFile: (path: string | undefined) => void;
  setFileContent: (content: string) => void;
  setFiles: (files: FileNode[]) => void;
}

const IDEContext = createContext<IDEContextValue | null>(null);

export function IDEProvider({
  children,
  files,
  filesLoading,
  isGithub,
  onSetFiles,
}: {
  children: ReactNode;
  files: FileNode[];
  filesLoading: boolean;
  isGithub: boolean;
  onSetFiles: (files: FileNode[]) => void;
}) {
  const [activeFile, setActiveFile] = useState<string | undefined>();
  const [fileContent, setFileContent] = useState("");

  console.log("files in provider: ", files);
  return (
    <IDEContext.Provider
      value={{
        files,
        filesLoading,
        isGithub,
        activeFile,
        fileContent,
        setActiveFile,
        setFileContent,
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
