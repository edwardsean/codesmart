import { QueryClient } from "@tanstack/react-query"

export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 1000 * 60,      // 1 minute fresh
      retry: 1,                   // retry once on failure
      refetchOnWindowFocus: true, // refetch when tab regains focus
    },
  },
})