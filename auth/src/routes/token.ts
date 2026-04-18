import { Hono } from "hono";
import { Bindings } from "../types";
import { getToken, deleteToken } from "../lib/kv";

const app = new Hono<{ Bindings: Bindings }>();

app.get("/", async (c) => {
  const { state, tid } = c.req.query();
  if (!state || !tid) return c.json({ error: "missing params" }, 400);

  const entry = await getToken(c.env.KV, state);
  if (!entry) return c.json({ error: "not found" }, 404);

  if (entry.telegramId !== tid) return c.json({ error: "unauthorized" }, 403);

  await deleteToken(c.env.KV, state);

  return c.json({ token: entry.token });
});

export default app;
