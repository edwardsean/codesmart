// services/fileService.ts
import { AxiosInstance } from "axios";
import { FileNode, FileContent } from "@/types/project.file.types";


export interface CreateFilePayload {
  path: string;
  content?: string;
  language?: string;
  is_dir?: boolean;
}

export interface UpdateFilePayload {
  path: string;    // rename
  content?: string; // save
}

export interface RenameFilePayload {
  old_path: string;
  new_path: string;
}

// converts flat list from API into nested tree structure
export function buildTree(flat: FileNode[]): FileNode[]{
  const root: FileNode[] = [];
  const map = new Map<string, FileNode>();

  for (const node of flat) { //for each node, make an object filenode, with children initialized to []
    map.set(node.path, {
      ...node,
      children: []
    })
  }

  //build parent-child relationships
  for (const node of flat) {
    const currentNode = map.get(node.path);
    const lastSlashIndex = node.path.lastIndexOf("/");

    if(currentNode) {
      if (lastSlashIndex === -1) {
        //no slash, this is a root node
        root.push(currentNode);
      } else  {
        //has parent, get the parent path
        const parentPath = node.path.substring(0, lastSlashIndex);
        const parent = map.get(parentPath);

        if(parent) {
          parent.children!.push(currentNode);
        } else {
          root.push(currentNode); //shouldnt happen
        }
      }
    }
    
  }

//sort directories first, then files, alphabetically
  const sortNodes = (nodes: FileNode[]) => {
    nodes.sort((a, b) => {
      if (a.type !== b.type) {
        return a.type === "dir" ? -1 : 1;
      }

      return a.name.localeCompare(b.name);
    });

    nodes.forEach(node => {
      if (node.children && node.children.length > 0) {
        sortNodes(node.children);
      }
    });
  };

  sortNodes(root);
  return root;
}


export const fileService = (api: AxiosInstance) => ({
  getFileTree: async (projectId: number): Promise<FileNode[]> => {
    const res = await api.get(`/api/projects/${projectId}/files`);
    console.log("API RESPONSE: ", res.data);
    const flat: FileNode[] = res.data.files ?? [];
    return buildTree(flat);
  },

  getFileContent: async (projectId: number, path: string): Promise<FileContent> => {
    const res = await api.get(`/api/projects/${projectId}/files/content`, {
      params: { path },
    });
    return res.data.file as FileContent;
  },

  createFile: async (projectId: number, payload: CreateFilePayload): Promise<FileContent> => {
    const res = await api.post(`/api/projects/${projectId}/files`, payload);
    return res.data.file as FileContent;
  },

  updateFile: async (projectId: number, payload: UpdateFilePayload): Promise<void> => {
    await api.put(`/api/projects/${projectId}/files`, payload);
  },

  renameFile: async (projectId: number, payload: RenameFilePayload): Promise<void> => {
    await api.put(`/api/projects/${projectId}/files/rename`, payload);
  },

  deleteFile: async (projectId: number, path: string): Promise<void> => {
    await api.delete(`/api/projects/${projectId}/files/${path}`);
  },
});