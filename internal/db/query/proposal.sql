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

-- name: ListProposalsByProject :many
SELECT pr.*, u.name AS freelancer_name
FROM proposals pr
JOIN users u ON u.id = pr.freelancer_id
WHERE pr.project_id = $1
ORDER BY pr.created_at DESC;