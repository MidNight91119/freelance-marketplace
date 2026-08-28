CREATE TYPE roles AS ENUM (
    'client',
    'freelancer'
);

CREATE TABLE "users" (
    "id" bigserial PRIMARY KEY,
    "name" varchar(255) NOT NULL,
    "email" varchar(255) NOT NULL UNIQUE,
    "hashed_password" varchar(255) NOT NULL,
    "role" roles NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);