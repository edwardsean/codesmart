import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import useAxiosPrivate from "@/hooks/useAxiosPrivate";
import { projectService } from "@/services/project.service";
import { useAuthStore } from "@/stores/auth.store";
import { ProjectListItem } from "@/types/project.types";


export const PROJECT_KEYS = {
    all: ["projects"] as const,
    detail: (id: number) => ["projects", id] as const
}

export function useProjects() {
    const api = useAxiosPrivate();
    const { _hasHydrated, account } = useAuthStore();

    return useQuery({ //use query = reads data from the cache or fetches it if not cached
        queryKey: PROJECT_KEYS.all,
        queryFn: () => projectService(api).getProjects(), //how to fetch if not in cache
        enabled: _hasHydrated && !!account?.accessToken,
    });
}

export function useProject(id: number) {
    const api = useAxiosPrivate();
    const { _hasHydrated, account } = useAuthStore();
    const queryClient = useQueryClient();

    return useQuery({
        queryKey: PROJECT_KEYS.detail(id),
        queryFn: async () => {
            const project = await projectService(api).getProjectByID(id);

            //sync detail data to the list cache
            queryClient.setQueryData(
                PROJECT_KEYS.all,
                (old: ProjectListItem[] | undefined) => {
                    if(!old) return old;
                    return old.map(p => 
                        p.id === project.id
                        ? {
                            ...p,
                            total_levels: project.total_levels,
                            completed_levels: project.completed_levels,
                            progress_percentage: project.progress_percentage,
                            status: project.status,
                        } : p
                    )
                }
            )
            
            return project;
        },
        enabled: _hasHydrated && !!account?.accessToken,
        refetchInterval: (query) => //poll every 5s while AI is generating levels
            query.state.data?.total_levels === 0 ? 5000 : false
    })
}

export function useDeleteProject() {
  const api = useAxiosPrivate();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: number) => projectService(api).deleteProject(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: PROJECT_KEYS.all });
    },
  });
}