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
export async function fetchGitHubUser(
  accessToken: string,
): Promise<{ login: string }> {
  const res = await fetch("https://api.github.com/user", {
    headers: {
      Authorization: `Bearer ${accessToken}`,
      Accept: "application/vnd.github+json",
      "User-Agent": "mahora-bot",
      "X-GitHub-Api-Version": "2022-11-28",
    },
  });

  if (!res.ok) {
    const body = await res.text();
    throw new Error(`fetchGitHubUser: ${res.status} ${body}`);
  }

  return res.json() as Promise<{ login: string }>;
}
