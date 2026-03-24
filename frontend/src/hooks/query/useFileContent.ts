import useAxiosPrivate from "@/hooks/useAxiosPrivate"
import { fileService, UpdateFilePayload } from "@/services/file.service";
import { useMutation, useQuery } from "@tanstack/react-query";
import { useAuthStore } from "@/stores/auth.store";
import { useQueryClient } from "@tanstack/react-query";

const FILE_CONTENT_KEYS = {
    content: (projectId: string, path: string) => ["file", projectId, path] as const,
}

export function useFileContent(projectId: string, path: string) {
    const api = useAxiosPrivate();
    const { _hasHydrated, account } = useAuthStore();

    return useQuery({
        queryKey: FILE_CONTENT_KEYS.content(projectId, path),
        queryFn: () => fileService(api).getFileContent(Number(projectId), path),
        enabled: _hasHydrated && !!account?.accessToken && !!path, //only fetch when a path is selected
        staleTime: 1000 * 60 * 5,
    });
}

export function usePrefetchFileContent(projectId: string) {
    const api = useAxiosPrivate()
    const queryClient = useQueryClient()

    return (path: string) => {
        queryClient.prefetchQuery({
        queryKey: FILE_CONTENT_KEYS.content(projectId, path),
        queryFn: () => fileService(api).getFileContent(Number(projectId), path),
        staleTime: 1000 * 60 * 5,
    });
  }
}

export function useUpdateFileContent() {
    const api = useAxiosPrivate();
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: ({projectId, payload}: {projectId: string, payload: UpdateFilePayload}) => fileService(api).updateFile(Number(projectId), payload),
        onSuccess: (data, { projectId, payload}) => queryClient.invalidateQueries({queryKey: FILE_CONTENT_KEYS.content(projectId, payload.path)})
    })
}