import { KVTokenEntry } from "../types";

const TOKEN_TTL = 300;

export async function saveToken(
  kv: KVNamespace,
  state: string,
  entry: KVTokenEntry,
): Promise<void> {
  await kv.put(state, JSON.stringify(entry), { expirationTtl: TOKEN_TTL });
}

export async function getToken(
  kv: KVNamespace,
  state: string,
): Promise<KVTokenEntry | null> {
  const raw = await kv.get(state);
  if (!raw) return null;
  return JSON.parse(raw) as KVTokenEntry;
}

export async function deleteToken(
  kv: KVNamespace,
  state: string,
): Promise<void> {
  await kv.delete(state);
}
