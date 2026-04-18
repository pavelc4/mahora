import { Hono } from "hono";
import { Bindings } from "../types";
import { verifyHMAC } from "../lib/hmac";
import { exchangeCode, fetchGitHubUser } from "../lib/github";
import { saveToken } from "../lib/kv";

const app = new Hono<{ Bindings: Bindings }>();

app.get("/", async (c) => {
  const { code, state } = c.req.query();
  if (!code || !state) return c.json({ error: "invalid request" }, 400);

  const parts = state.split(":");
  if (parts.length !== 3) return c.json({ error: "invalid state" }, 400);

  const [uuid, tid, sig] = parts;

  const valid = await verifyHMAC(c.env.WORKER_SECRET, `${uuid}:${tid}`, sig);
  if (!valid) return c.json({ error: "invalid state signature" }, 403);

  let token: string;
  try {
    token = await exchangeCode(
      code,
      c.env.GITHUB_CLIENT_ID,
      c.env.GITHUB_CLIENT_SECRET,
    );
  } catch (e) {
    return c.json({ error: "failed to exchange code" }, 500);
  }

  let githubLogin: string;
  try {
    const user = await fetchGitHubUser(token);
    githubLogin = user.login;
  } catch (e) {
    return c.json({ error: "failed to fetch github user" }, 500);
  }

  await saveToken(c.env.KV, uuid, { token, telegramId: tid, githubLogin });

  // UI yang sudah dipercantik dengan Tailwind CSS & Material 3 vibes
  return c.html(`
    <!DOCTYPE html>
    <html lang="id">
      <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>Autentikasi Berhasil</title>
        <script src="https://cdn.tailwindcss.com"></script>
        <style>
          @import url('https://fonts.googleapis.com/css2?family=Outfit:wght@400;600;700&display=swap');
          body { font-family: 'Outfit', sans-serif; }
        </style>
      </head>
      <body class="bg-slate-50 flex items-center justify-center min-h-screen p-4">
        <div class="bg-white max-w-md w-full rounded-[32px] p-8 text-center shadow-[0_8px_30px_rgb(0,0,0,0.04)] border border-slate-100">

          <div class="w-20 h-20 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-6">
            <svg class="w-10 h-10 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M5 13l4 4L19 7"></path>
            </svg>
          </div>

          <h2 class="text-3xl font-bold text-slate-900 mb-3">GitHub Connected!</h2>
          <p class="text-slate-500 mb-8 text-lg leading-relaxed">
           Autentication successful. You can close this tab and return to Telegram.
          </p>

          </div>
      </body>
    </html>
  `);
});

export default app;
