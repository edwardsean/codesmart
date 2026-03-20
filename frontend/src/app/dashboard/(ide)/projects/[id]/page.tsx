// app/dashboard/(ide)/projects/[id]/page.tsx
"use client";

import dynamic from "next/dynamic";
import { Sparkles } from "lucide-react";
import { useIDE } from "@/context/IDEProvider";

const CodeEditor = dynamic(() => import("@/components/ide/CodeEditor"), {
  ssr: false,
  loading: () => (
    <div className="h-full bg-[#1e1e1e] flex items-center justify-center">
      <p className="text-xs font-mono text-zinc-500">Loading editor...</p>
    </div>
  ),
});

export default function ProjectCodePage() {
  const { activeFile, fileContent, setFileContent } = useIDE();

  return (
    <div className="h-full flex">
      {/* ── Editor area ── */}
      <div className="flex-1 flex flex-col min-w-0">
        {/* Tab bar — only show if a file is open */}
        {activeFile && (
          <div
            className="flex items-center border-b border-black/[0.06] dark:border-white/[0.06]
            bg-gray-50 dark:bg-[#252526] flex-shrink-0 h-8 overflow-x-auto"
          >
            <div
              className="flex items-center h-full px-4 gap-2
              border-r border-black/[0.06] dark:border-white/[0.06]
              bg-white dark:bg-[#1e1e1e]"
            >
              <span className="text-xs font-mono text-gray-600 dark:text-zinc-400">
                {activeFile.split("/").pop()}
              </span>
            </div>
          </div>
        )}

        {/* Editor or welcome screen */}
        <div className="flex-1 overflow-hidden">
          {activeFile ? (
            <CodeEditor
              value={fileContent}
              path={activeFile}
              onChange={setFileContent}
            />
          ) : (
            // no file selected — show welcome
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
          )}
        </div>

        {/* Status bar — only show if a file is open */}
        {activeFile && (
          <div
            className="h-6 flex-shrink-0 flex items-center justify-between px-4
            bg-[#dc503c] text-white text-xs font-mono"
          >
            <span>{activeFile.split(".").pop()}</span>
            <span>UTF-8</span>
          </div>
        )}
      </div>

      {/* ── Smarty panel ── */}
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
          <div>
            <p className="text-sm font-medium text-gray-900 dark:text-zinc-100 mb-1">
              No levels yet
            </p>
            <p className="text-xs text-gray-400 dark:text-zinc-600 leading-relaxed">
              Smarty will analyze your project and generate a structured
              learning path.
            </p>
          </div>
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
