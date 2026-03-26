"use client";

import { useRef, useEffect } from "react";
import Editor, { OnMount } from "@monaco-editor/react";
import * as monaco from "monaco-editor";

interface CodeEditorProps {
  value: string;
  language?: string;
  onChange?: (value: string) => void;
  onSave?: () => void;
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
  language,
  onChange,
  onSave,
  readOnly = false,
  path,
}: CodeEditorProps) {
  const editorRef = useRef<monaco.editor.IStandaloneCodeEditor | null>(null);
  const onSaveRef = useRef(onSave);
  const isMountedRef = useRef(true);

  useEffect(() => {
    onSaveRef.current = onSave;
  }, [onSave]);

  useEffect(() => {
    isMountedRef.current = true;
    return () => {
      isMountedRef.current = false;
    };
  }, []);

  const ext = path?.split(".").pop() ?? "";
  const monacoLang =
    LANG_MAP[ext] || LANG_MAP[language ?? ""] || language || "plaintext";

  const isDark =
    typeof document !== "undefined"
      ? document.documentElement.classList.contains("dark")
      : true;

  const handleMount: OnMount = (editor, monacoInstance) => {
    if (isMountedRef.current) {
      editorRef.current = editor;
      editor.focus();

      editor.addCommand(
        monacoInstance.KeyMod.CtrlCmd | monacoInstance.KeyCode.KeyS,
        () => {
          if (isMountedRef.current) {
            onSaveRef.current?.();
          }
        },
      );
    }
  };

  return (
    <Editor
      height="100%"
      language={monacoLang}
      value={value}
      theme={isDark ? "vs-dark" : "light"}
      onChange={(val) => {
        if (isMountedRef.current) {
          onChange?.(val ?? "");
        }
      }}
      onMount={handleMount}
      path={path}
      options={{
        fontSize: 13,
        fontFamily: "'JetBrains Mono', 'Fira Code', monospace",
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
