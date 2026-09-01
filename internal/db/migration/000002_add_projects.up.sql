CREATE TYPE project_status AS ENUM (
  'open',
  'in_progress',
  'completed',
  'cancelled'
);

CREATE TABLE "projects" (
  "id" bigserial PRIMARY KEY,
  "client_id" bigint NOT NULL REFERENCES "users" ("id"),
  "title" varchar(255) NOT NULL,
  "description" text NOT NULL,
  "category" varchar(255) NOT NULL,
  "budget_min" bigint NOT NULL,
  "budget_max" bigint NOT NULL,
  "status" project_status NOT NULL DEFAULT 'open',
  "deadline" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),

  CONSTRAINT budget_positive CHECK (budget_min > 0),
  CONSTRAINT budget_range CHECK (budget_max >= budget_min),
  CONSTRAINT deadline_after_creation CHECK (deadline > created_at)
);