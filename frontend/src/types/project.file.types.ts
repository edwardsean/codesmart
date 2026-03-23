export interface FileNode {
  id: number;
  path: string;
  name: string;
  type: "file" | "dir";
  children?: FileNode[];
}

export interface FileContent {
  id: number;
  path: string;
  name: string;
  content: string;
}
