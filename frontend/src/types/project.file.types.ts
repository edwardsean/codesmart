export interface FileNode {
  path: string;
  name: string;
  type: "file" | "dir";
  children?: FileNode[];
}

export interface FileContent {
  path: string;
  name: string;
  content: string;
}
