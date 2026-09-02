-- name: CreateProject :one
INSERT INTO "projects" (
  client_id,
  title,
  description,
  category,
  budget_min,
  budget_max,
  deadline
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: ListProjects :many
SELECT p.*, u.name AS client_name
FROM projects p
JOIN users u ON u.id = p.client_id
WHERE p.status = 'open'
  AND (sqlc.narg(category)::varchar IS NULL OR p.category = sqlc.narg(category))
  AND (sqlc.narg(min_budget)::bigint IS NULL OR p.budget_max >= sqlc.narg(min_budget))
  AND (sqlc.narg(max_budget)::bigint IS NULL OR p.budget_min <= sqlc.narg(max_budget))
ORDER BY p.created_at DESC;