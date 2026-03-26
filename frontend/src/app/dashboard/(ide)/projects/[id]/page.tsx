"use client";

import dynamic from "next/dynamic";
import { Sparkles } from "lucide-react";
import { useIDE } from "@/context/IDEContext";
import {
  useFileContent,
  useUpdateFileContent,
} from "@/hooks/query/useFileContent";
import { X } from "lucide-react";
import { useState, useEffect } from "react";
import TerminalPanel from "@/components/ide/TerminalPanel";

const CodeEditor = dynamic(() => import("@/components/ide/CodeEditor"), {
  ssr: false,
  loading: () => (
    <div className="h-full bg-[#1e1e1e] flex items-center justify-center">
      <p className="text-xs font-mono text-zinc-500">Loading editor...</p>
    </div>
  ),
});

export default function ProjectCodePage() {
  const {
    activeFile,
    projectId,
    openTabs,
    openFile,
    closeTab,
    editedContent,
    setEditedContent,
    clearEditedContent,
  } = useIDE();

  const { data: fileData, isLoading: contentLoading } = useFileContent(
    projectId,
    activeFile ?? "",
  );
  console.log(`file for ${activeFile}: `, fileData);
  const { mutate: updateFile } = useUpdateFileContent();

  //to show in editor either unsaved local edits (typed but not saved) or server content from react query
  const displayContent = activeFile
    ? (editedContent[activeFile] ?? fileData?.content ?? "")
    : "";

  //unsaved if the active file now is edited and not the same as from the react query
  const isUnsaved =
    activeFile &&
    editedContent[activeFile] !== undefined &&
    editedContent[activeFile] !== fileData?.content;

  const [saveStatus, setSaveStatus] = useState<
    "idle" | "saving" | "saved" | "error"
  >("idle");
  const [terminalOpen, setTerminalOpen] = useState(false);

  useEffect(() => {
    const handler = (e: KeyboardEvent) => {
      if ((e.ctrlKey || e.metaKey) && e.key === "`") {
        setTerminalOpen((o) => !o);
      }
    };
    window.addEventListener("keydown", handler);
    return () => window.removeEventListener("keydown", handler);
  }, []);

  async function handleSave() {
    if (!activeFile || !fileData) return; //if no active file or not in workspace
    const content = editedContent[activeFile];

    if (content === undefined) return; //nothing to save

    setSaveStatus("saving");
    updateFile(
      {
        projectId,
        payload: { content, path: activeFile },
      },
      {
        onSuccess: () => {
          clearEditedContent(activeFile);
          setSaveStatus("saved");
          setTimeout(() => setSaveStatus("idle"), 2000);
        },
        onError: () => {
          setSaveStatus("error");
          setTimeout(() => setSaveStatus("idle"), 2000);
        },
      },
    );
  }

  return (
    <div className="h-full flex">
      <div className="flex-1 flex flex-col min-w-0">
        {/* Tab bar */}
        {openTabs.length > 0 && (
          <div
            className="flex items-center border-b border-black/[0.06] dark:border-white/[0.06]
            bg-gray-50 dark:bg-[#252526] flex-shrink-0 h-8 overflow-x-auto"
          >
            {openTabs.map((path) => {
              const isActive = path === activeFile;
              const isDirty = path in editedContent;

              return (
                <div
                  key={path}
                  onClick={() => openFile(path)}
                  className={`group flex items-center gap-2 h-full px-3 border-r
                    border-black/[0.06] dark:border-white/[0.06]
                    cursor-pointer flex-shrink-0 transition-colors ${
                      isActive
                        ? "bg-white dark:bg-[#1e1e1e] text-gray-900 dark:text-zinc-100"
                        : "text-gray-400 dark:text-zinc-500 hover:bg-gray-100 dark:hover:bg-zinc-800"
                    }`}
                >
                  <span className="text-xs font-mono">
                    {path.split("/").pop()}
                  </span>

                  <div className="w-4 h-4 flex items-center justify-center flex-shrink-0">
                    {isDirty ? (
                      <div
                        className="w-1.5 h-1.5 rounded-full bg-[#dc503c]"
                        onClick={(e) => {
                          e.stopPropagation();
                          closeTab(path);
                        }}
                      />
                    ) : (
                      <button
                        onClick={(e) => {
                          e.stopPropagation();
                          closeTab(path);
                        }}
                        className="opacity-0 group-hover:opacity-100 transition-opacity
                          text-gray-400 hover:text-gray-700 dark:hover:text-zinc-200"
                      >
                        <X size={11} />
                      </button>
                    )}
                  </div>
                </div>
              );
            })}
          </div>
        )}

        {/* Editor or welcome */}
        <div className="flex-1 overflow-hidden">
          {!activeFile ? (
            <div className="h-full bg-white dark:bg-[#1e1e1e] flex flex-col items-center justify-center gap-3">
              <svg
                width="32"
                height="32"
                viewBox="0 0 14 14"
                fill="none"
                opacity={0.2}
              >
                <path
                  d="M2 4.5L5.5 7L2 9.5"
                  stroke="currentColor"
                  strokeWidth="1.5"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                />
                <path
                  d="M7.5 9.5H12"
                  stroke="currentColor"
                  strokeWidth="1.5"
                  strokeLinecap="round"
                />
              </svg>
              <p className="text-xs font-mono text-gray-300 dark:text-zinc-600">
                Select a file to start editing
              </p>
            </div>
          ) : contentLoading ? (
            <div className="h-full bg-white dark:bg-[#1e1e1e] flex items-center justify-center">
              <p className="text-xs font-mono text-zinc-500">Loading...</p>
            </div>
          ) : (
            <CodeEditor
              value={displayContent}
              path={activeFile}
              onChange={(val) => setEditedContent(activeFile, val)}
              onSave={handleSave}
            />
          )}
        </div>

        {terminalOpen && (
          <TerminalPanel
            projectId={projectId}
            onClose={() => setTerminalOpen(false)}
          />
        )}
        {/* Status bar */}
        {activeFile && (
          <div
            className="h-6 flex-shrink-0 flex items-center justify-between px-4
    bg-[#dc503c] text-white text-xs font-mono"
          >
            <span>{activeFile.split(".").pop()}</span>

            <span className="flex items-center gap-1.5 transition-all">
              {saveStatus === "saving" && (
                <>
                  <svg
                    className="animate-spin w-2.5 h-2.5"
                    viewBox="0 0 24 24"
                    fill="none"
                  >
                    <circle
                      className="opacity-25"
                      cx="12"
                      cy="12"
                      r="10"
                      stroke="currentColor"
                      strokeWidth="4"
                    />
                    <path
                      className="opacity-75"
                      fill="currentColor"
                      d="M4 12a8 8 0 018-8v8z"
                    />
                  </svg>
                  saving...
                </>
              )}
              {saveStatus === "saved" && "saved"}
              {saveStatus === "error" && "save failed"}
              {saveStatus === "idle" && (isUnsaved ? "unsaved" : "UTF-8")}
            </span>
          </div>
        )}
      </div>

      {/* Smarty panel */}
      <div
        className="w-72 flex-shrink-0 border-l border-black/[0.06] dark:border-white/[0.06]
        flex flex-col bg-gray-50 dark:bg-[#252526]"
      >
        <div className="px-4 py-3 border-b border-black/[0.06] dark:border-white/[0.06] flex items-center gap-2">
          <Sparkles size={14} className="text-[#dc503c]" />
          <span className="text-xs font-mono text-gray-600 dark:text-zinc-400 uppercase tracking-wide">
            Smarty
          </span>
        </div>
        <div className="flex-1 flex flex-col items-center justify-center p-6 text-center gap-4">
          <div className="w-12 h-12 rounded-xl bg-[#dc503c]/10 flex items-center justify-center">
            <Sparkles size={20} className="text-[#dc503c]" />
          </div>
          <p className="text-sm font-medium text-gray-900 dark:text-zinc-100 mb-1">
            No levels yet
          </p>
          <p className="text-xs text-gray-400 dark:text-zinc-600 leading-relaxed">
            Smarty will analyze your project and generate a structured learning
            path.
          </p>
          <button
            className="flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium
            bg-[#dc503c] text-white hover:opacity-85 transition-opacity w-full justify-center"
          >
            <Sparkles size={14} />
            Generate levels
          </button>
        </div>
      </div>
    </div>
  );
}
