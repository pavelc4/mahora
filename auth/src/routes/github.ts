import { Hono } from "hono";
import { Bindings } from "../types";
import { generateHMAC } from "../lib/hmac";
import { v4 as uuidv4 } from "uuid";

const app = new Hono<{ Bindings: Bindings }>();

app.get("/", async (c) => {
  const tid = c.req.query("tid");
  if (!tid) return c.json({ error: "missing telegram_id" }, 400);

  const uuid = uuidv4();
  const sig = await generateHMAC(c.env.WORKER_SECRET, `${uuid}:${tid}`);
  const state = `${uuid}:${tid}:${sig}`;

  const params = new URLSearchParams({
    client_id: c.env.GITHUB_CLIENT_ID,
    redirect_uri: `${c.env.WORKER_BASE_URL}/callback`,
    scope: "repo,notifications",
    state,
  });

  return c.redirect(`https://github.com/login/oauth/authorize?${params}`);
});

export default app;
