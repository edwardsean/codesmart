"use client";

import { useRef, useState, useCallback, useEffect } from "react";

interface ResizablePanelProps {
  children: React.ReactNode;
  defaultWidth?: number;
  minWidth?: number;
  maxWidth?: number;
  collapsed?: boolean;
  onCollapse?: (collapsed: boolean) => void;
  side?: "left" | "right";
}

export default function ResizablePanel({
  children,
  defaultWidth = 220,
  minWidth = 160,
  maxWidth = 480,
  collapsed = false,
  onCollapse,
  side = "left",
}: ResizablePanelProps) {
  const [width, setWidth] = useState(defaultWidth);
  const [isDragging, setIsDragging] = useState(false);
  const panelRef = useRef<HTMLDivElement>(null);
  const startXRef = useRef(0);
  const startWidthRef = useRef(0);

  const onMouseDown = useCallback(
    (e: React.MouseEvent) => {
      e.preventDefault();
      setIsDragging(true);
      startXRef.current = e.clientX;
      startWidthRef.current = width;
    },
    [width],
  );

  useEffect(() => {
    if (!isDragging) return;

    const onMouseMove = (e: MouseEvent) => {
      const delta =
        side === "left"
          ? e.clientX - startXRef.current
          : startXRef.current - e.clientX;
      const newWidth = Math.max(
        minWidth,
        Math.min(maxWidth, startWidthRef.current + delta),
      );
      setWidth(newWidth);
    };

    const onMouseUp = () => setIsDragging(false);

    window.addEventListener("mousemove", onMouseMove);
    window.addEventListener("mouseup", onMouseUp);
    return () => {
      window.removeEventListener("mousemove", onMouseMove);
      window.removeEventListener("mouseup", onMouseUp);
    };
  }, [isDragging, minWidth, maxWidth, side]);

  if (collapsed) return null;

  return (
    <div
      ref={panelRef}
      className="relative flex-shrink-0 h-full overflow-hidden"
      style={{ width }}
    >
      {/* Panel content */}
      <div className="h-full overflow-hidden">{children}</div>

      {/* Drag handle */}
      <div
        onMouseDown={onMouseDown}
        onDoubleClick={() => onCollapse?.(true)}
        title="Drag to resize · Double-click to collapse"
        className={`absolute top-0 ${side === "left" ? "right-0" : "left-0"} w-1 h-full cursor-col-resize
          hover:bg-[#dc503c]/40 transition-colors z-10
          ${isDragging ? "bg-[#dc503c]/60" : "bg-transparent"}`}
      />
    </div>
  );
}
