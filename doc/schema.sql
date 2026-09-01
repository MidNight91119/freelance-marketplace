-- SQL dump generated using DBML (dbml.dbdiagram.io)
-- Database: PostgreSQL
-- Generated at: 2026-09-01T11:41:16.981Z

CREATE TYPE "roles" AS ENUM (
  'client',
  'freelancer'
);

CREATE TYPE "project_status" AS ENUM (
  'open',
  'in_progress',
  'completed',
  'cancelled'
);

CREATE TYPE "proposal_status" AS ENUM (
  'pending',
  'accepted',
  'rejected'
);

CREATE TYPE "contract_status" AS ENUM (
  'active',
  'completed',
  'cancelled'
);

CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "name" varchar(255) NOT NULL,
  "email" varchar(255) UNIQUE NOT NULL,
  "hashed_password" varchar(255) NOT NULL,
  "role" roles NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "projects" (
  "id" bigserial PRIMARY KEY,
  "client_id" bigint NOT NULL,
  "title" varchar(255) NOT NULL,
  "description" text NOT NULL,
  "category" varchar(255) NOT NULL,
  "budget_min" bigint NOT NULL,
  "budget_max" bigint NOT NULL,
  "status" project_status NOT NULL DEFAULT 'open',
  "deadline" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "proposals" (
  "id" bigserial PRIMARY KEY,
  "project_id" bigint NOT NULL,
  "freelancer_id" bigint NOT NULL,
  "cover_letter" text NOT NULL,
  "proposed_price" bigint NOT NULL,
  "estimated_duration_days" bigint NOT NULL,
  "status" proposal_status NOT NULL DEFAULT 'pending',
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "contracts" (
  "id" bigserial PRIMARY KEY,
  "project_id" bigint UNIQUE NOT NULL,
  "proposal_id" bigint UNIQUE NOT NULL,
  "client_id" bigint NOT NULL,
  "freelancer_id" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "status" contract_status NOT NULL DEFAULT 'active',
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE UNIQUE INDEX ON "proposals" ("project_id", "freelancer_id");

COMMENT ON COLUMN "projects"."budget_max" IS 'must be >= budget_min';

COMMENT ON COLUMN "projects"."deadline" IS 'must be > created_at';

ALTER TABLE "projects" ADD FOREIGN KEY ("client_id") REFERENCES "users" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "proposals" ADD FOREIGN KEY ("project_id") REFERENCES "projects" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "proposals" ADD FOREIGN KEY ("freelancer_id") REFERENCES "users" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "contracts" ADD FOREIGN KEY ("project_id") REFERENCES "projects" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "contracts" ADD FOREIGN KEY ("proposal_id") REFERENCES "proposals" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "contracts" ADD FOREIGN KEY ("client_id") REFERENCES "users" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "contracts" ADD FOREIGN KEY ("freelancer_id") REFERENCES "users" ("id") DEFERRABLE INITIALLY IMMEDIATE;
