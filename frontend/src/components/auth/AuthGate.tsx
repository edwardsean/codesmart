"use client";

import { useAuth } from "@/hooks/useAuth";
import Loading from "@/components/ui/Loading";

export default function AuthGate({ children }: { children: React.ReactNode }) {
  //   const { loading } = useAuth();

  //   if (loading) return <Loading />;

  return <>{children}</>;
}
