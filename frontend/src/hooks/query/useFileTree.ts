import useAxiosPrivate from "@/hooks/useAxiosPrivate";
import { fileService, RenameFilePayload, CreateFilePayload } from "@/services/file.service";
import { useAuthStore } from "@/stores/auth.store";
import { FileNode } from "@/types/project.file.types";
import { useQuery, useMutation } from "@tanstack/react-query";
import { useQueryClient } from "@tanstack/react-query";
import { FILE } from "dns";
import { rename } from "fs";
import { use } from "react";

const FILE_TREE_KEYS = {
    tree: (projectId: string) => ["file-tree", projectId] as const,
}

export function useFileTree(projectId: string) {
    const api = useAxiosPrivate();
    const { _hasHydrated, account } = useAuthStore();

    return useQuery({
        queryKey: FILE_TREE_KEYS.tree(projectId),
        queryFn: () => fileService(api).getFileTree(Number(projectId)),
        enabled: _hasHydrated && !!account?.accessToken, //only fetch when we have a valid access token
        staleTime: 1000 * 60 * 5, //cache for 5 minutes
    });
}

export function useCreateFile(projectId: string) {
    const api = useAxiosPrivate();
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: (payload: CreateFilePayload) => 
            fileService(api).createFile(Number(projectId), payload),
        onMutate: async ({ path, is_dir }) => {
            await queryClient.cancelQueries({ queryKey: FILE_TREE_KEYS.tree(projectId) });
            const previousTree = queryClient.getQueryData(FILE_TREE_KEYS.tree(projectId));

            queryClient.setQueryData(FILE_TREE_KEYS.tree(projectId), (old: FileNode[] | undefined) => {
                if (!old) return [];

                //create the new node
                const newNode: FileNode = {
                    path, 
                    name: path.split("/").slice(-1)[0],
                    type: is_dir ? "dir" : "file",
                    children: is_dir ? [] : undefined,
                }

                return insertNode(old, newNode);
            });

            return { previousTree }; //if the mutation fails, we can roll back to this state
        },
        onError: (err, variables, context) => {
            if (context?.previousTree) {
                queryClient.setQueryData(FILE_TREE_KEYS.tree(projectId), context.previousTree);
            }
        },
        
    })
}

export function useRenameFile(projectId: string) {
    const api = useAxiosPrivate();
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: (payload: RenameFilePayload) => fileService(api).renameFile(Number(projectId), payload),
        onMutate: async ({ old_path, new_path }) => {
            await queryClient.cancelQueries({ queryKey: FILE_TREE_KEYS.tree(projectId) });
            const previousTree = queryClient.getQueryData(FILE_TREE_KEYS.tree(projectId));

            queryClient.setQueryData(FILE_TREE_KEYS.tree(projectId), (old: FileNode[] | undefined) => {
                if (!old) return [];
                return renameNode(old, old_path, new_path);
            });

            return { previousTree }
        },
        onError: (err, variables, context) => {
            if (context?.previousTree) {
                queryClient.setQueryData(FILE_TREE_KEYS.tree(projectId), context.previousTree);
            }
        },
    })
}

export function useDeleteFile(projectId: string) {
    const api = useAxiosPrivate();
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: (path: string) => fileService(api).deleteFile(Number(projectId), path),
        onMutate: async (path: string) => {
            await queryClient.cancelQueries({ queryKey: FILE_TREE_KEYS.tree(projectId) });
            const previousTree = queryClient.getQueryData(FILE_TREE_KEYS.tree(projectId));
        
            queryClient.setQueryData(FILE_TREE_KEYS.tree(projectId), (old: FileNode[] | undefined) => {
                if (!old) return [];
                return deleteNode(old, path);
            });

            return { previousTree };
        },
        onError: (err, variables, context) => {
            if (context?.previousTree) {
                queryClient.setQueryData(FILE_TREE_KEYS.tree(projectId), context.previousTree);
            }
        },
        
    })
}

function deleteNode(tree: FileNode[], targetPath: string): FileNode[] {
    return tree.filter(node => node.path !== targetPath)
        .map(node => {
            if (node.children) {
                return {
                    ...node,
                    children: deleteNode(node.children, targetPath), //recurse into children to find the node to delete
                }
            }

            return node; //no match, return as is
        })
}
              
function renameNode(tree: FileNode[], oldPath: string, newPath: string): FileNode[] {
    return tree.map(node => {
        if (node.path === oldPath) {
            return {
                ...node,
                path: newPath,
                name: newPath.split("/").pop()!, 
            };
        }
        if (node.children) {
            //if node is parent
            //update children paths that start with oldPath to newPath
            const prefix = oldPath + "/";
            if(node.path.startsWith(prefix)) {
                const updatedPath = node.path.replace(prefix, newPath + "/");
                return {
                    ...node,
                    path: updatedPath,
                    children: node.children,
                }
            }
            return {
                ...node,
                children: renameNode(node.children, oldPath, newPath), //recurse into children to find the node to rename
            }   
        }
        return node; //no match, return as is
    })

}
function insertNode(tree: FileNode[], newNode: FileNode): FileNode[] {
    const parts = newNode.path.split("/");
    const parentPath = parts.slice(0, -1).join("/"); //creates a new array containing all elements except the last one

    //base case, root level
    if (parentPath === "") {
        return [...tree, newNode];
    }

    return tree.map(node => {
        if (node.path === parentPath && node.children) { //found the parent, insert the new node here
            return {
                ...node,
                children: [...node.children, newNode],
            };
        }
        if(node.children) { //recurse into children to find the parent
            return {
                ...node,
                children: insertNode(node.children, newNode),
            }
        }

        return node; //no match, return as is
    });

    
}