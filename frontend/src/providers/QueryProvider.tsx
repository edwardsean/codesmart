"use client";

import { QueryClientProvider } from "@tanstack/react-query";
import { queryClient } from "@/lib/queryClient";

export default function QueryProvider({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    //makes the query client accessible to every component inside
    //without the provider, useQuery and useMutation dont know which cache to use and will throw an error
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
}
