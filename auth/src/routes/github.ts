export default app;

app.get("/", async (c) => {
  const { tid, state } = c.req.query();
  if (!tid || !state) return c.json({ error: "missing params" }, 400);

  const parts = state.split(":");
  if (parts.length !== 3) return c.json({ error: "invalid state" }, 400);
  const [uuid, stateTid, sig] = parts;

  const valid = await verifyHMAC(
    c.env.WORKER_SECRET,
    `${uuid}:${stateTid}`,
    sig,
  );
  if (!valid) return c.json({ error: "invalid state signature" }, 403);

  const params = new URLSearchParams({
    client_id: c.env.GITHUB_CLIENT_ID,
    redirect_uri: `${c.env.WORKER_BASE_URL}/callback`,
    scope: "repo,notifications",
    state,
  });

  return c.redirect(`https://github.com/login/oauth/authorize?${params}`);
});
