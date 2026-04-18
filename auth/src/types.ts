export interface Bindings {
  KV: KVNamespace;
  GITHUB_CLIENT_ID: string;
  GITHUB_CLIENT_SECRET: string;
  WORKER_SECRET: string;
  WORKER_BASE_URL: string;
}

export interface KVTokenEntry {
  token: string;
  telegramId: string;
  githubLogin: string;
}
