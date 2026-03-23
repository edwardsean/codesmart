import { useMutation, useQueryClient } from "@tanstack/react-query";
import useAxiosPrivate from "@/hooks/useAxiosPrivate";
import { projectService } from "@/services/project.service";
import { PROJECT_KEYS } from "@/hooks/query/useProjects";
import { CreateProjectRequest } from "@/services/project.service";

export function useCreateProjects() {
    const api = useAxiosPrivate();
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: (payload: CreateProjectRequest) => 
            projectService(api).createProject(payload)
        ,
        onSuccess: () => {
            //invalidates projects list, refetches, new project appears
            queryClient.invalidateQueries({queryKey: PROJECT_KEYS.all});
        }
    })
}
