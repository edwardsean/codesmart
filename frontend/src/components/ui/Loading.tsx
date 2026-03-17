import React from "react";

interface LoadingScreenProps {
  message?: string;
}

const DOTS = 8;

export default function LoadingScreen({
  message = "Loading...",
}: LoadingScreenProps) {
  return (
    <div className="fixed inset-0 flex flex-col items-center justify-center bg-white">
      <div className="relative w-16 h-16">
        {Array.from({ length: DOTS }).map((_, i) => {
          const angle = (i / DOTS) * 360;
          const rad = (angle * Math.PI) / 180;
          const x = 50 + 38 * Math.cos(rad) - 6;
          const y = 50 + 38 * Math.sin(rad) - 6;
          const delay = (i / DOTS) * 0.8;

          return (
            <span
              key={i}
              style={{
                position: "absolute",
                left: `${x}%`,
                top: `${y}%`,
                width: 7,
                height: 7,
                borderRadius: "50%",
                backgroundColor: "#ef4444",
                animation: `lakoo-pulse 0.8s ease-in-out ${delay}s infinite`,
              }}
            />
          );
        })}
      </div>

      {message && (
        <p
          style={{ fontFamily: "Georgia, serif" }}
          className="mt-6 text-sm text-gray-400 tracking-widest uppercase"
        >
          {message}
        </p>
      )}

      <style>{`
        @keyframes lakoo-pulse {
          0%, 100% { opacity: 0.15; transform: scale(0.8); }
          50%       { opacity: 1;    transform: scale(1.2); }
        }
      `}</style>
    </div>
  );
}
