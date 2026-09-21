# Freelance Marketplace

A full-stack marketplace where clients post projects, freelancers submit proposals, and accepting a proposal
creates a contract inside a single database transaction.

**Live:** [frontend](https://freelance-marketplace-web.onrender.com) · [backend](https://freelance-marketplace-8yz8.onrender.com)

Both run on Render's free tier and sleep after 15 minutes idle — the first request after a pause takes ~30s to wake. Not a bug.

Built as Project 2 of the Super30 grind. Backend is hand-written Go on the standard library; the frontend is
React + TypeScript consuming the real API — no mock data anywhere.

---

## What it does


| Role           | Can                                                                                                         |
| -------------- | ----------------------------------------------------------------------------------------------------------- |
| **Client**     | create projects · view own projects · view proposals on own projects · accept one proposal · view contracts |
| **Freelancer** | browse and filter open projects · submit one proposal per project · view own proposals · view contracts     |


Accepting a proposal is the centrepiece: in one transaction it marks that proposal `accepted`, every other
proposal on the project `rejected`, the project `in_progress`, and inserts a `contracts` row snapshotting the
agreed price. Any failure rolls the whole thing back — this is tested against a real Postgres, including the
rollback path.

## Stack


| Layer      | Choice                                                                          | Why                                                                                                                                                                      |
| ---------- | ------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| Language   | **Go 1.26**                                                                     | *static binary, fast compile, stdlib HTTP that doesn't need a framework*                                                                                                 |
| HTTP       | `net/http` **stdlib** — method+pattern routing, hand-written middleware chain   | Go 1.22+ made frameworks optional. Auth and role middleware are plain `func(http.Handler) http.Handler`, so the router reads as a map of the whole authorization surface |
| Database   | **Postgres + pgx/v5**                                                           | Contracts and payments need ACID and real foreign keys                                                                                                                   |
| Queries    | **sqlc**                                                                        | Hand-written SQL, generated type-safe Go. Errors move from runtime to compile time                                                                                       |
| Migrations | **golang-migrate**                                                              | Numbered, append-only, reversible                                                                                                                                        |
| Auth       | **JWT (**`golang-jwt/v5`**) + bcrypt**                                          | Spec-mandated; the keyfunc asserts the signing method to close algorithm-confusion                                                                                       |
| Validation | `go-playground/validator` on request structs, business rules as explicit checks | Server-side only enforces anything; the client is UX                                                                                                                     |
| Config     | **viper**, `app.env` locally, env vars in production                            |                                                                                                                                                                          |
| Frontend   | **React 19 + TypeScript + Vite + React Router**                                 | Market standard; one hand-written stylesheet, no UI framework                                                                                                            |
| Deploy     | **Render**                                                                      | Proven pipeline                                                                                                                                                          |


**Deliberately not used:** Redis (no measured read hotspot), a job queue (nothing is slow), gRPC (only consumer
is a browser), Kubernetes (one binary, one database).

## API

All routes except signup/login require `Authorization: Bearer <jwt>`.


| Method | Path                           | Who               | Notes                                                                                                               |
| ------ | ------------------------------ | ----------------- | ------------------------------------------------------------------------------------------------------------------- |
| POST   | `/api/auth/signup`             | public            | bcrypt, unique email, role ∈ {client, freelancer}                                                                   |
| POST   | `/api/auth/login`              | public            | returns a JWT carrying `userId`, `role`, `email`                                                                    |
| POST   | `/api/projects`                | client            | `budgetMin > 0`, `budgetMax ≥ budgetMin`, deadline in the future                                                    |
| GET    | `/api/projects`                | any               | open projects; filters `category`, `minBudget`, `maxBudget` (range overlap); returns `clientName` + `proposalCount` |
| GET    | `/api/projects/mine`           | client            | caller's projects at any status                                                                                     |
| POST   | `/api/projects/{id}/proposals` | freelancer        | one per freelancer per project — enforced by a UNIQUE constraint, not an application check                          |
| GET    | `/api/projects/{id}/proposals` | client, **owner** | ownership checked in the handler against the token's user id                                                        |
| GET    | `/api/proposals/mine`          | freelancer        | caller's proposals with the project title joined in                                                                 |
| PUT    | `/api/proposals/{id}/accept`   | client, **owner** | **the transaction**                                                                                                 |
| GET    | `/api/contracts`               | any               | scoped to the caller as client or freelancer — the `WHERE` clause is the authorization                              |


Errors are `{ "code": "...", "message": "..." }` with a fixed vocabulary
(`INVALID_REQUEST`, `UNAUTHORIZED`, `FORBIDDEN`, `INVALID_CREDENTIALS`, `EMAIL_ALREADY_EXISTS`,
`PROJECT_NOT_FOUND`, `PROJECT_NOT_OPEN`, `PROPOSAL_ALREADY_EXISTS`, `PROPOSAL_NOT_FOUND`,
`PROPOSAL_ALREADY_PROCESSED`).

## Design decisions**

- **Role checks are middleware; ownership checks are in handlers.** Role is knowable from the token alone, so
it's route-level. Ownership needs a database read, so it can only live where the row is loaded.
- **Invariants live in the database.** Duplicate proposals, one contract per project, budget ordering — all
constraints. Application checks lose races; constraints don't. The Go validation is for good error messages.
- `contracts.amount` **is a snapshot, not a cache.** If the freelancer edits their proposal later, the signed
contract must not move.
- **CHECK constraints never reference** `now()`**.** They're re-evaluated on every `UPDATE`; "deadline in the
future" is a rule about the moment of creation and belongs in Go.
- **No** `?mine=true` **flag.** "My projects" is a separate route because it's role-gated (a flag can't be
middleware-gated) and because it changes what "the list" means (any status vs open only).
- **Frontend auth state is React state.** `localStorage` is persistence only; an `AuthProvider` owns `role`,
so login/logout re-render every consumer without a navigation as a side effect.



## Schema

Four tables: `users` → `projects` → `proposals` → `contracts`. See `[doc/db.dbml](doc/db.dbml)`.

Key constraints: `UNIQUE (project_id, freelancer_id)` on proposals; `UNIQUE (project_id)` and
`UNIQUE (proposal_id)` on contracts; `CHECK (budget_max >= budget_min)`; enums for role and every status.

## Tests

```
make test    # db tests against a real Postgres + handler tests against a gomock Store
make smoke   # 51-check end-to-end run against a live server
```

- `internal/db/sqlc` — every query against a real database, including the accept-proposal transaction
and its rollback (the failure is arranged on the *last* step so earlier writes are proven undone).
- `internal/api` — handlers against a mock `Store`: signup, login, and the role-gated `mine` routes
(200 / 403 / 401, with `Times(0)` proving middleware stopped the request before the handler).
- `scripts/smoke.sh` — signs up a client and two freelancers, walks the full flow, and asserts every
status code including the rejections. Re-runnable.

Not yet covered by handler tests: routes 3–8. They're covered end to end by smoke.

## Running locally

Requires Go 1.26+, Node 20+, Docker.

```bash
make postgres && make createdb && make migrateup   # Postgres in Docker, schema applied
cp app.env.example app.env                          # or write DB_SOURCE, TOKEN_SYMMETRIC_KEY, etc.
make server                                         # :8080

cd web && npm install && npm run dev                # :5173, talks to :8080 by default
```

Set `VITE_API_URL` to point the frontend elsewhere. The backend reads its config from `app.env` locally and
from environment variables in production.

## Layout

```
cmd/api/main.go            entrypoint
internal/api/              handlers, router, middleware, response helpers
internal/db/migration/     golang-migrate, append-only
internal/db/query/         hand-written SQL   ← source of truth
internal/db/sqlc/          generated Go       ← never hand-edited
internal/db/mock/          gomock Store
internal/token/            JWT maker
internal/util/             config, password hashing
scripts/smoke.sh           end-to-end checks
web/                       React + Vite + TypeScript
```

