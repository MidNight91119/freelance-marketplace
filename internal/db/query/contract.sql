-- name: CreateContract :one
INSERT INTO "contracts" (
  project_id,
  proposal_id,
  client_id,
  freelancer_id,
  amount
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;

-- name: ListContractsByUser :many
SELECT
  c.*,
  p.title AS project_title,
  cl.name AS client_name,
  f.name AS freelancer_name
FROM "contracts" c
JOIN projects p ON p.id = c.project_id
JOIN users cl ON cl.id = c.client_id
JOIN users f ON f.id = c.freelancer_id
WHERE c.client_id = $1 OR c.freelancer_id = $1
ORDER BY c.created_at DESC;