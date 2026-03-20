"use client";

import { useRef } from "react";
import Editor, { OnMount } from "@monaco-editor/react";
import { useAuthStore } from "@/stores/authStore";

interface CodeEditorProps {
  value: string;
  language?: string;
  onChange?: (value: string) => void;
  readOnly?: boolean;
  path?: string;
}

const LANG_MAP: Record<string, string> = {
  go: "go",
  ts: "typescript",
  tsx: "typescript",
  js: "javascript",
  jsx: "javascript",
  py: "python",
  sql: "sql",
  md: "markdown",
  css: "css",
  json: "json",
  yaml: "yaml",
  yml: "yaml",
};

export default function CodeEditor({
  value,
  language = "plaintext",
  onChange,
  readOnly = false,
  path,
}: CodeEditorProps) {
  const editorRef = useRef<unknown>(null);
  const { account } = useAuthStore();

  // derive monaco language from file extension or language prop
  const ext = path?.split(".").pop() ?? "";
  const monacoLang = LANG_MAP[ext] || LANG_MAP[language] || language;

  const handleMount: OnMount = (editor) => {
    editorRef.current = editor;
    editor.focus();
  };

  // determine theme based on current html class
  const isDark =
    typeof document !== "undefined"
      ? document.documentElement.classList.contains("dark")
      : true;

  return (
    <Editor
      height="100%"
      language={monacoLang}
      value={value}
      theme={isDark ? "vs-dark" : "light"}
      onChange={(val) => onChange?.(val ?? "")}
      onMount={handleMount}
      path={path} // keeps separate undo/redo history per file
      options={{
        fontSize: 13,
        fontFamily: "'JetBrains Mono', 'Fira Code', 'Cascadia Code', monospace",
        fontLigatures: true,
        minimap: { enabled: false },
        scrollBeyondLastLine: false,
        lineNumbers: "on",
        wordWrap: "on",
        tabSize: 2,
        readOnly,
        automaticLayout: true,
        padding: { top: 16 },
        renderLineHighlight: "gutter",
        smoothScrolling: true,
        cursorBlinking: "smooth",
        bracketPairColorization: { enabled: true },
      }}
    />
  );
}
