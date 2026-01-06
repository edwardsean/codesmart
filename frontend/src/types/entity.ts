export type User = {
    id: number;
    email: string;
    username: string;
    github_id: number;
    createdAt: string;
}

export interface Repository {
  id: string;
  name: string;
  full_name: string;
  description: string;
  html_url: string
  language: string
  stargazers_count: number
  forks_count: number
  updated_at: string
  private: boolean
  default_branch: string
  owner: {
    login: string
    avatar_url: string
  }
}