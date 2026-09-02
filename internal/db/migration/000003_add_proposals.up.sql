CREATE TYPE "proposal_status" AS ENUM (
  'pending',
  'accepted',
  'rejected'
);

CREATE TABLE "proposals" (
  "id" bigserial PRIMARY KEY,
  "project_id" bigint NOT NULL REFERENCES "projects" ("id"),
  "freelancer_id" bigint NOT NULL REFERENCES "users" ("id"),
  "cover_letter" text NOT NULL,
  "proposed_price" bigint NOT NULL,
  "estimated_duration_days" bigint NOT NULL,
  "status" proposal_status NOT NULL DEFAULT 'pending',
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),

  CONSTRAINT one_proposal_per_freelancer UNIQUE (project_id, freelancer_id),
  CONSTRAINT price_positive CHECK (proposed_price > 0),
  CONSTRAINT duration_positive CHECK (estimated_duration_days > 0)
);