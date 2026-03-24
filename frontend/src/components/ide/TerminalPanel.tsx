"use client";

import { useEffect, useRef, useState, useCallback } from "react";
import { Terminal as XTerm } from "@xterm/xterm";
import { FitAddon } from "@xterm/addon-fit";
import { WebLinksAddon } from "@xterm/addon-web-links";
import "@xterm/xterm/css/xterm.css";
import useAxiosPrivate from "@/hooks/useAxiosPrivate";
import { Plus, X, ChevronUp, Minus } from "lucide-react";

interface TerminalInstance {
  id: string;
  name: string;
  term: XTerm;
  fitAddon: FitAddon;
  ws: WebSocket | null;
}

interface TerminalPanelProps {
  projectId: string;
  onClose: () => void;
}

// ── Single terminal tab ──
function TerminalTab({
  instance,
  active,
  onClick,
  onClose,
}: {
  instance: TerminalInstance;
  active: boolean;
  onClick: () => void;
  onClose: () => void;
}) {
  return (
    <div
      onClick={onClick}
      className={`group flex items-center gap-2 px-3 h-full border-r
        border-black/[0.06] dark:border-white/[0.06]
        cursor-pointer flex-shrink-0 transition-colors text-xs font-mono ${
          active
            ? "bg-[#1e1e1e] text-zinc-200"
            : "text-zinc-500 hover:text-zinc-300 hover:bg-[#1e1e1e]/50"
        }`}
    >
      <span>{instance.name}</span>
      <button
        onClick={(e) => {
          e.stopPropagation();
          onClose();
        }}
        className="opacity-0 group-hover:opacity-100 transition-opacity hover:text-red-400"
      >
        <X size={10} />
      </button>
    </div>
  );
}

// ── Individual terminal instance ──
function TerminalInstance({
  projectId,
  instance,
  visible,
  api,
}: {
  projectId: string;
  instance: TerminalInstance;
  visible: boolean;
  api: ReturnType<typeof useAxiosPrivate>;
}) {
  const containerRef = useRef<HTMLDivElement>(null);
  const connectedRef = useRef(false);

  useEffect(() => {
    if (!containerRef.current || connectedRef.current) return;
    connectedRef.current = true;

    // mount terminal into this container
    // instance.term.open(containerRef.current);
    // instance.fitAddon.fit();

    const connect = async () => {
      let wsUrl: string;

      const base =
        process.env.NODE_ENV === "development"
          ? "ws://localhost:8080"
          : process.env.NEXT_PUBLIC_WS_URL;
      try {
        const { data } = await api.post("/api/auth/ws-ticket");
        wsUrl = `${base}/api/v1/projects/${projectId}/terminal?ticket=${data.ticket}`;
      } catch {
        instance.term.writeln("\x1b[31mFailed to authenticate terminal\x1b[0m");
        return;
      }

      const ws = new WebSocket(wsUrl);
      ws.binaryType = "arraybuffer";
      instance.ws = ws;

      ws.onopen = () => {
        instance.term.writeln("\x1b[32mConnected\x1b[0m");
        // send initial size
        ws.send(
          JSON.stringify({
            type: "resize",
            cols: instance.term.cols,
            rows: instance.term.rows,
          }),
        );
      };

      ws.onclose = () => {
        instance.term.writeln("\r\n\x1b[31mDisconnected\x1b[0m");
      };

      ws.onerror = () => {
        instance.term.writeln("\r\n\x1b[31mConnection error\x1b[0m");
      };

      // ws.onmessage = (event) => {
      //   instance.term.write(event.data);
      // };
      ws.onmessage = (event) => {
        if (event.data instanceof ArrayBuffer) {
          // binary message from PTY — write as Uint8Array
          instance.term.write(new Uint8Array(event.data));
        } else {
          // text message
          instance.term.write(event.data);
        }
      };
      instance.term.onData((data) => {
        if (ws.readyState === WebSocket.OPEN) {
          ws.send(data);
        }
      });
    };

    setTimeout(() => {
      if (!containerRef.current) return;
      instance.term.open(containerRef.current); //this renders the terminal's canvas into the div, before this, XTerm exists in memory but isn't visible
      instance.fitAddon.fit(); //calculate cols/rows based on div size, resize termina
      instance.term.focus(); //give keyboard focus to the terminal so keypresses go to it
      connect();
    }, 10);

    return () => {
      instance.ws?.close(); // This function is RETURNED but NOT executed yet, will run when unmount
    };
  }, []);

  // fit when visibility changes
  useEffect(() => {
    if (visible) {
      setTimeout(() => {
        instance.fitAddon.fit();
        instance.term.focus();
      }, 50);
    }
  }, [visible]);

  return (
    <div
      ref={containerRef}
      className="h-full w-full bg-[#1e1e1e]"
      style={{ display: visible ? "block" : "none" }}
      onClick={() => instance.term.focus()}
    />
  );
}

// ── Terminal Panel (the full thing with tabs, resize handle) ──
export default function TerminalPanel({
  projectId,
  onClose,
}: TerminalPanelProps) {
  const api = useAxiosPrivate();
  const [instances, setInstances] = useState<TerminalInstance[]>([]);
  const [activeId, setActiveId] = useState<string | null>(null);
  const [height, setHeight] = useState(220);
  const [isMinimized, setIsMinimized] = useState(false);
  const isDragging = useRef(false);
  const startY = useRef(0);
  const startHeight = useRef(0);
  const counter = useRef(1);

  function createInstance(): TerminalInstance {
    const term = new XTerm({
      //create the terminal UI in the browser
      theme: {
        background: "#1e1e1e",
        foreground: "#d4d4d4",
        cursor: "#dc503c",
        selectionBackground: "#dc503c40",
      },
      fontSize: 13,
      fontFamily: "'JetBrains Mono', 'Fira Code', Menlo, monospace",
      cursorBlink: true,
      cursorStyle: "block",
      scrollback: 1000,
    });

    const fitAddon = new FitAddon(); //FitAddon makes terminal fill its container
    term.loadAddon(fitAddon);
    term.loadAddon(new WebLinksAddon());

    return {
      id: crypto.randomUUID(),
      name: `bash ${counter.current++}`,
      term,
      fitAddon,
      ws: null,
    };
  }

  // create first terminal on mount
  useEffect(() => {
    const first = createInstance();
    setInstances([first]);
    setActiveId(first.id);
  }, []);

  function addTerminal() {
    const instance = createInstance();
    setInstances((prev) => [...prev, instance]);
    setActiveId(instance.id);
  }

  function closeTerminal(id: string) {
    setInstances((prev) => {
      const next = prev.filter((t) => t.id !== id);

      // find the instance to close its WebSocket
      const closing = prev.find((t) => t.id === id);
      closing?.ws?.close();
      closing?.term.dispose();

      // if closing active tab, switch to adjacent
      if (id === activeId) {
        const idx = prev.findIndex((t) => t.id === id);
        const nextActive = next[idx] ?? next[idx - 1];
        setActiveId(nextActive?.id ?? null);
      }

      // if no terminals left, close the panel
      if (next.length === 0) onClose();

      return next;
    });
  }

  // ── Drag to resize ──
  const onMouseDown = useCallback(
    (e: React.MouseEvent) => {
      isDragging.current = true;
      startY.current = e.clientY;
      startHeight.current = height;
      e.preventDefault();
    },
    [height],
  );

  useEffect(() => {
    const onMouseMove = (e: MouseEvent) => {
      if (!isDragging.current) return;
      const delta = startY.current - e.clientY; // drag up = bigger
      const newHeight = Math.max(
        120,
        Math.min(600, startHeight.current + delta),
      );
      setHeight(newHeight);
    };

    const onMouseUp = () => {
      isDragging.current = false;
      // refit all terminals after resize
      instances.forEach((t) => t.fitAddon.fit());
    };

    window.addEventListener("mousemove", onMouseMove);
    window.addEventListener("mouseup", onMouseUp);
    return () => {
      window.removeEventListener("mousemove", onMouseMove);
      window.removeEventListener("mouseup", onMouseUp);
    };
  }, [instances]);

  return (
    <div
      className="flex-shrink-0 flex flex-col border-t border-black/[0.06] dark:border-white/[0.06]"
      style={{ height: isMinimized ? 28 : height }}
    >
      {/* Drag handle — hover to resize */}
      {!isMinimized && (
        <div
          onMouseDown={onMouseDown}
          className="h-1 flex-shrink-0 cursor-row-resize hover:bg-[#dc503c]/40 transition-colors"
          title="Drag to resize"
        />
      )}

      {/* Terminal header */}
      <div
        className="flex items-center h-7 flex-shrink-0
        bg-gray-50 dark:bg-[#252526] border-b border-black/[0.04] dark:border-white/[0.04]"
      >
        {/* Tab list */}
        <div className="flex items-center h-full overflow-x-auto flex-1">
          {instances.map((instance) => (
            <TerminalTab
              key={instance.id}
              instance={instance}
              active={instance.id === activeId}
              onClick={() => {
                setActiveId(instance.id);
                setTimeout(() => instance.term.focus(), 50);
              }}
              onClose={() => closeTerminal(instance.id)}
            />
          ))}
        </div>

        {/* Actions */}
        <div className="flex items-center gap-0.5 px-2 flex-shrink-0">
          {/* New terminal */}
          <button
            onClick={addTerminal}
            className="p-1 text-gray-400 dark:text-zinc-500
              hover:text-gray-700 dark:hover:text-zinc-300 transition-colors"
            title="New terminal"
          >
            <Plus size={12} />
          </button>

          {/* Minimize/expand */}
          <button
            onClick={() => setIsMinimized((m) => !m)}
            className="p-1 text-gray-400 dark:text-zinc-500
              hover:text-gray-700 dark:hover:text-zinc-300 transition-colors"
            title={isMinimized ? "Expand" : "Minimize"}
          >
            {isMinimized ? <ChevronUp size={12} /> : <Minus size={12} />}
          </button>

          {/* Close panel */}
          <button
            onClick={onClose}
            className="p-1 text-gray-400 dark:text-zinc-500
              hover:text-red-500 transition-colors"
            title="Close terminal"
          >
            <X size={12} />
          </button>
        </div>
      </div>

      {/* Terminal content */}
      {!isMinimized && (
        <div className="flex-1 overflow-hidden relative">
          {instances.map((instance) => (
            <div
              key={instance.id}
              className="absolute inset-0 p-1"
              style={{ display: instance.id === activeId ? "block" : "none" }}
            >
              <TerminalInstance
                projectId={projectId}
                instance={instance}
                visible={instance.id === activeId}
                api={api}
              />
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
