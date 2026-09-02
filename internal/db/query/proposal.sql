-- name: CreateProposal :one
INSERT INTO "proposals" (
  project_id,
  freelancer_id,
  cover_letter,
  proposed_price,
  estimated_duration_days
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;