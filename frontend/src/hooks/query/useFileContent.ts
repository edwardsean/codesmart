import useAxiosPrivate from "@/hooks/useAxiosPrivate"
import { fileService, UpdateFilePayload } from "@/services/file.service";
import { useMutation, useQuery } from "@tanstack/react-query";
import { useAuthStore } from "@/stores/auth.store";
import { useQueryClient } from "@tanstack/react-query";

const FILE_CONTENT_KEYS = {
    content: (projectId: string, path: string | undefined, fileId: number | undefined) => ["file", projectId, fileId, path] as const,
}

export function useFileContent(projectId: string, path: string | undefined, fileId: number | undefined) {
    const api = useAxiosPrivate();
    const { _hasHydrated, account } = useAuthStore();

    return useQuery({
        queryKey: FILE_CONTENT_KEYS.content(projectId, path, fileId),
        queryFn: () => fileService(api).getFileContent(Number(projectId), path!),
        enabled: _hasHydrated && !!account?.accessToken && !!path, //only fetch when a path is selected
        staleTime: 1000 * 60 * 5,
    });
}

export function usePrefetchFileContent(projectId: string) {
    const api = useAxiosPrivate()
    const queryClient = useQueryClient()

    return (path: string, fileId: number) => {
        queryClient.prefetchQuery({
        queryKey: FILE_CONTENT_KEYS.content(projectId, path, fileId),
        queryFn: () => fileService(api).getFileContent(Number(projectId), path),
        staleTime: 1000 * 60 * 5,
    });
  }
}

export function useUpdateFileContent() {
    const api = useAxiosPrivate();
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: ({projectId, fileId, payload}: {projectId: string, fileId: number, payload: UpdateFilePayload}) => fileService(api).updateFile(Number(projectId), fileId, payload),
        onSuccess: (data, { projectId, fileId, payload}) => queryClient.invalidateQueries({queryKey: FILE_CONTENT_KEYS.content(projectId, payload.path, fileId)})
    })
}