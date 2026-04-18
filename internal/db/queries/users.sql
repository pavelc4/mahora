SELECT * FROM users WHERE telegram_id = ? LIMIT 1;

INSERT INTO users (telegram_id, github_login, github_token, updated_at)
VALUES (?, ?, ?, strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
ON CONFLICT(telegram_id) DO UPDATE SET
    github_login = excluded.github_login,
    github_token = excluded.github_token,
    updated_at   = excluded.updated_at
RETURNING *;

UPDATE users SET github_token = NULL, github_login = NULL,
    updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now')
WHERE telegram_id = ?;
