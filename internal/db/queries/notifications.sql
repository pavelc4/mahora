-- name: InsertNotification :exec
INSERT INTO notifications (user_id, repo_full, event_type, event_id)
VALUES (?, ?, ?, ?)
ON CONFLICT(user_id, repo_full, event_type, event_id) DO NOTHING;

-- name: HasNotification :one
SELECT COUNT(*) FROM notifications
WHERE user_id = ? AND repo_full = ? AND event_type = ? AND event_id = ?;

-- name: DeleteOldNotifications :exec
DELETE FROM notifications
WHERE sent_at < strftime('%Y-%m-%dT%H:%M:%fZ', 'now', '-90 days');
