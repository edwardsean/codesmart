"use client";

import React from "react";
import { useInitAuth } from "@/hooks/useInitAuth";

export default function InitAuth({ children }: { children: React.ReactNode }) {
  useInitAuth();
  return <>{children}</>;
}
