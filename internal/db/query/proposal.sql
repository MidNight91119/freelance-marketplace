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

-- name: GetProposal :one
SELECT * FROM proposals
WHERE id = $1
LIMIT 1;

-- name: AcceptProposal :one
UPDATE proposals
SET status = 'accepted', updated_at = now()
WHERE id = $1
RETURNING *;

-- name: RejectOtherProposals :exec
UPDATE proposals
SET status = 'rejected', updated_at = now()
WHERE project_id = $1 AND id <> $2;