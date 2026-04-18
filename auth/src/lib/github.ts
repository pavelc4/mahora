export async function exchangeCode(
  code: string,
  clientId: string,
  clientSecret: string,
): Promise<string> {
  const resp = await fetch("https://github.com/login/oauth/access_token", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Accept: "application/json",
    },
    body: JSON.stringify({
      client_id: clientId,
      client_secret: clientSecret,
      code,
    }),
  });

  const data = (await resp.json()) as { access_token?: string; error?: string };

  if (data.error || !data.access_token) {
    throw new Error(data.error ?? "failed to exchange code");
  }

  return data.access_token;
}
