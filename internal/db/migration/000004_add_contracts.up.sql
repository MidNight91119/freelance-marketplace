CREATE TYPE "contract_status" AS ENUM (
  'active',
  'completed',
  'cancelled'
);

CREATE TABLE "contracts" (
  "id" bigserial PRIMARY KEY,
  "project_id" bigint NOT NULL REFERENCES "projects" ("id"),
  "proposal_id" bigint NOT NULL REFERENCES "proposals" ("id"),
  "client_id" bigint NOT NULL REFERENCES "users" ("id"),
  "freelancer_id" bigint NOT NULL REFERENCES "users" ("id"),
  "amount" bigint NOT NULL,
  "status" contract_status NOT NULL DEFAULT 'active',
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),

  CONSTRAINT one_contract_per_project UNIQUE (project_id),
  CONSTRAINT one_contract_per_proposal UNIQUE (proposal_id),
  CONSTRAINT amount_positive CHECK (amount > 0)
);