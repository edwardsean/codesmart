// services/fileService.ts
import { AxiosInstance } from "axios";
import { FileNode, FileContent } from "@/types/entity";


interface CreateFilePayload {
  path: string;
  content?: string;
  language?: string;
  is_dir?: boolean;
}

interface UpdateFilePayload {
  path?: string;    // rename
  content?: string; // save
}

// converts flat list from API into nested tree structure
export function buildTree(flat: FileNode[]): FileNode[] {
  const root: FileNode[] = [];
  const map = new Map<string, FileNode>();

  // sort so dirs come before files, then alphabetically
  const sorted = [...flat].sort((a, b) => {
    if (a.type !== b.type) return a.type === "dir" ? -1 : 1;
    return a.path.localeCompare(b.path);
  });

  for (const node of sorted) {
    const parts = node.path.split("/");
    if (parts.length === 1) {
      root.push(node);
      map.set(node.path, node);
    } else {
      const parentPath = parts.slice(0, -1).join("/");
      const parent = map.get(parentPath);
      if (parent) {
        if (!parent.children) parent.children = [];
        parent.children.push(node);
      }
      map.set(node.path, node);
    }
  }

  return root;
}

export const fileService = (api: AxiosInstance) => ({
  getFileTree: async (projectId: number): Promise<FileNode[]> => {
    const res = await api.get(`/api/projects/${projectId}/files`);
    const flat: FileNode[] = res.data.files ?? [];
    return buildTree(flat);
  },

  getFileContent: async (projectId: number, path: string): Promise<FileContent> => {
    const res = await api.get(`/api/projects/${projectId}/files/content`, {
      params: { path },
    });
    return res.data.file;
  },

  createFile: async (projectId: number, payload: CreateFilePayload): Promise<FileContent> => {
    const res = await api.post(`/api/projects/${projectId}/files`, payload);
    return res.data.file;
  },

  updateFile: async (projectId: number, fileId: number, payload: UpdateFilePayload): Promise<void> => {
    await api.put(`/api/projects/${projectId}/files/${fileId}`, payload);
  },

  deleteFile: async (projectId: number, fileId: number): Promise<void> => {
    await api.delete(`/api/projects/${projectId}/files/${fileId}`);
  },
});