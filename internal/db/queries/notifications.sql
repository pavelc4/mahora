INSERT INTO notifications (user_id, repo_full, event_type, event_id)
VALUES (?, ?, ?, ?)
ON CONFLICT(user_id, repo_full, event_type, event_id) DO NOTHING;

SELECT COUNT(*) FROM notifications
WHERE user_id = ? AND repo_full = ? AND event_type = ? AND event_id = ?;
