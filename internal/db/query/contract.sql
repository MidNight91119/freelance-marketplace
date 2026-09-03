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