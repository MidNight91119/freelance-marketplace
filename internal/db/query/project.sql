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