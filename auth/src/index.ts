import { Hono } from "hono";
import { Bindings } from "./types";
import github from "./routes/github";
import callback from "./routes/callback";
import token from "./routes/token";

const app = new Hono<{ Bindings: Bindings }>();

app.route("/auth/github", github);
app.route("/callback", callback);
app.route("/auth/token", token);

app.get("/health", (c) => c.json({ status: "ok" }));

app.onError((err, c) => {
  console.error(err);
  return c.json({ error: "internal server error" }, 500);
});

export default app;
