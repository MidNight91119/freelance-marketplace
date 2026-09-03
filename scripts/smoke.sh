#!/usr/bin/env bash
# Manual end-to-end smoke test. Assumes `make server` is already running.
#   ./scripts/smoke.sh
#
# Signs up a client + two freelancers with random emails (so it is re-runnable),
# then walks the whole flow and prints the status code for each step.

set -uo pipefail
BASE="${BASE:-http://localhost:8080}"
PASS="secret123"

pass=0; fail=0

# check <expected-status> <label> <curl args...>
check() {
  local want=$1 label=$2; shift 2
  local out code body
  out=$(curl -s -w $'\n%{http_code}' "$@")
  code=${out##*$'\n'}
  body=${out%$'\n'*}
  if [[ "$code" == "$want" ]]; then
    printf '  \033[32m✓\033[0m %-3s %s\n' "$code" "$label"
    pass=$((pass+1))
  else
    printf '  \033[31m✗\033[0m %-3s %s \033[2m(wanted %s)\033[0m\n      %s\n' "$code" "$label" "$want" "$body"
    fail=$((fail+1))
  fi
}

# body-only helper, for when we need to read an id or a token out of the response
get() { curl -s "$@"; }

signup() { # name email role
  get -X POST "$BASE/api/auth/signup" -H 'Content-Type: application/json' \
    -d "{\"name\":\"$1\",\"email\":\"$2\",\"password\":\"$PASS\",\"role\":\"$3\"}" >/dev/null
}

login() { # email  -> prints token
  get -X POST "$BASE/api/auth/login" -H 'Content-Type: application/json' \
    -d "{\"email\":\"$1\",\"password\":\"$PASS\"}" |
    sed -n 's/.*"accessToken":"\([^"]*\)".*/\1/p'
}

json_id() { sed -n 's/.*"id":\([0-9]*\).*/\1/p' | head -1; }

R=$RANDOM
CA="clienta$R@example.com"
CB="clientb$R@example.com"
F1="free1$R@example.com"
F2="free2$R@example.com"

echo "── auth ─────────────────────────────────────────────"
check 201 "signup client A"            -X POST "$BASE/api/auth/signup" -H 'Content-Type: application/json' -d "{\"name\":\"Client A\",\"email\":\"$CA\",\"password\":\"$PASS\",\"role\":\"client\"}"
check 409 "signup same email again"    -X POST "$BASE/api/auth/signup" -H 'Content-Type: application/json' -d "{\"name\":\"Client A\",\"email\":\"$CA\",\"password\":\"$PASS\",\"role\":\"client\"}"
check 400 "signup with bad role"       -X POST "$BASE/api/auth/signup" -H 'Content-Type: application/json' -d "{\"name\":\"X\",\"email\":\"x$R@example.com\",\"password\":\"$PASS\",\"role\":\"admin\"}"

signup "Client B" "$CB" client
signup "Free One" "$F1" freelancer
signup "Free Two" "$F2" freelancer

check 200 "login with correct password" -X POST "$BASE/api/auth/login" -H 'Content-Type: application/json' -d "{\"email\":\"$CA\",\"password\":\"$PASS\"}"
check 401 "login with wrong password"   -X POST "$BASE/api/auth/login" -H 'Content-Type: application/json' -d "{\"email\":\"$CA\",\"password\":\"wrongpassword\"}"
check 401 "login with unknown email"    -X POST "$BASE/api/auth/login" -H 'Content-Type: application/json' -d "{\"email\":\"nobody$R@example.com\",\"password\":\"$PASS\"}"

TA=$(login "$CA"); TB=$(login "$CB"); T1=$(login "$F1"); T2=$(login "$F2")
AUTH_A=(-H "Authorization: Bearer $TA")
AUTH_B=(-H "Authorization: Bearer $TB")
AUTH_1=(-H "Authorization: Bearer $T1")
AUTH_2=(-H "Authorization: Bearer $T2")
JSON=(-H 'Content-Type: application/json')

echo
echo "── auth middleware ──────────────────────────────────"
check 200 "GET /me with valid token"   "$BASE/api/auth/me" "${AUTH_A[@]}"
check 401 "GET /me with no header"     "$BASE/api/auth/me"
check 401 "GET /me with garbage token" "$BASE/api/auth/me" -H "Authorization: Bearer not.a.token"
check 401 "GET /me with wrong scheme"  "$BASE/api/auth/me" -H "Authorization: Basic $TA"

echo
echo "── projects ─────────────────────────────────────────"
PROJ="{\"title\":\"Smoke $R\",\"description\":\"desc\",\"category\":\"Smoke$R\",\"budgetMin\":50000,\"budgetMax\":100000,\"deadline\":\"2026-12-15\"}"
check 403 "freelancer creates project"    -X POST "$BASE/api/projects" "${AUTH_1[@]}" "${JSON[@]}" -d "$PROJ"
check 401 "anonymous creates project"     -X POST "$BASE/api/projects" "${JSON[@]}" -d "$PROJ"
check 400 "budgetMax below budgetMin"     -X POST "$BASE/api/projects" "${AUTH_A[@]}" "${JSON[@]}" -d "{\"title\":\"x\",\"description\":\"d\",\"category\":\"c\",\"budgetMin\":90000,\"budgetMax\":1000,\"deadline\":\"2026-12-15\"}"
check 400 "deadline in the past"          -X POST "$BASE/api/projects" "${AUTH_A[@]}" "${JSON[@]}" -d "{\"title\":\"x\",\"description\":\"d\",\"category\":\"c\",\"budgetMin\":1000,\"budgetMax\":9000,\"deadline\":\"2020-01-01\"}"

# one request: check the status code AND capture the id from the same response
OUT=$(curl -s -w $'\n%{http_code}' -X POST "$BASE/api/projects" "${AUTH_A[@]}" "${JSON[@]}" -d "$PROJ")
CODE=${OUT##*$'\n'}
BODY=${OUT%$'\n'*}
PID=$(echo "$BODY" | json_id)
if [[ "$CODE" == "201" ]]; then
  printf '  \033[32m✓\033[0m %-3s %s\n' "$CODE" "client creates project"
  pass=$((pass+1))
else
  printf '  \033[31m✗\033[0m %-3s %s \033[2m(wanted 201)\033[0m\n      %s\n' "$CODE" "client creates project" "$BODY"
  fail=$((fail+1))
fi
echo "  (working project id = $PID)"

check 200 "list projects"                 "$BASE/api/projects" "${AUTH_1[@]}"
check 200 "list filtered by category"     "$BASE/api/projects?category=Smoke$R" "${AUTH_1[@]}"
check 400 "list with non-numeric budget"  "$BASE/api/projects?minBudget=abc" "${AUTH_1[@]}"
check 401 "list without a token"          "$BASE/api/projects"

echo
echo "── proposals ────────────────────────────────────────"
PROP='{"coverLetter":"I have 3 years of experience","proposedPrice":75000,"estimatedDuration":20}'
check 201 "freelancer 1 submits"             -X POST "$BASE/api/projects/$PID/proposals" "${AUTH_1[@]}" "${JSON[@]}" -d "$PROP"
check 409 "freelancer 1 submits again"       -X POST "$BASE/api/projects/$PID/proposals" "${AUTH_1[@]}" "${JSON[@]}" -d "$PROP"
check 201 "freelancer 2 submits"             -X POST "$BASE/api/projects/$PID/proposals" "${AUTH_2[@]}" "${JSON[@]}" -d "$PROP"
check 403 "client submits a proposal"        -X POST "$BASE/api/projects/$PID/proposals" "${AUTH_A[@]}" "${JSON[@]}" -d "$PROP"
check 404 "propose to missing project"       -X POST "$BASE/api/projects/999999/proposals" "${AUTH_1[@]}" "${JSON[@]}" -d "$PROP"
check 400 "propose with bad path id"         -X POST "$BASE/api/projects/abc/proposals" "${AUTH_1[@]}" "${JSON[@]}" -d "$PROP"
check 400 "propose with missing coverLetter" -X POST "$BASE/api/projects/$PID/proposals" "${AUTH_1[@]}" "${JSON[@]}" -d '{"proposedPrice":5000,"estimatedDuration":20}'

echo
echo "── viewing proposals (ownership) ────────────────────"
check 200 "owner views proposals"          "$BASE/api/projects/$PID/proposals" "${AUTH_A[@]}"
check 403 "OTHER CLIENT views proposals"   "$BASE/api/projects/$PID/proposals" "${AUTH_B[@]}"
check 403 "freelancer views proposals"     "$BASE/api/projects/$PID/proposals" "${AUTH_1[@]}"
check 401 "anonymous views proposals"      "$BASE/api/projects/$PID/proposals"
check 404 "proposals of missing project"   "$BASE/api/projects/999999/proposals" "${AUTH_A[@]}"

echo
echo "── accepting a proposal (the transaction) ───────────"
# grab the two proposal ids on the working project
PROPS=$(get "$BASE/api/projects/$PID/proposals" "${AUTH_A[@]}")
P1=$(echo "$PROPS" | tr '}' '\n' | sed -n 's/.*"proposalId":\([0-9]*\).*/\1/p' | head -1)
P2=$(echo "$PROPS" | tr '}' '\n' | sed -n 's/.*"proposalId":\([0-9]*\).*/\1/p' | tail -1)
echo "  (proposals $P1 and $P2)"

check 403 "freelancer accepts"             -X PUT "$BASE/api/proposals/$P1/accept" "${AUTH_1[@]}"
check 403 "OTHER CLIENT accepts"           -X PUT "$BASE/api/proposals/$P1/accept" "${AUTH_B[@]}"
check 401 "anonymous accepts"              -X PUT "$BASE/api/proposals/$P1/accept"
check 404 "accept missing proposal"        -X PUT "$BASE/api/proposals/999999/accept" "${AUTH_A[@]}"
check 400 "accept with bad path id"        -X PUT "$BASE/api/proposals/abc/accept" "${AUTH_A[@]}"
check 200 "OWNER accepts proposal 1"       -X PUT "$BASE/api/proposals/$P1/accept" "${AUTH_A[@]}"
check 409 "accept the same one again"      -X PUT "$BASE/api/proposals/$P1/accept" "${AUTH_A[@]}"
check 409 "accept the OTHER proposal"      -X PUT "$BASE/api/proposals/$P2/accept" "${AUTH_A[@]}"

echo "  after accepting:"
check 200 "project dropped from open listing" "$BASE/api/projects?category=Smoke$R" "${AUTH_A[@]}"
echo "    open listing for this category → $(get "$BASE/api/projects?category=Smoke$R" "${AUTH_A[@]}")  (should be [], project is now in_progress)"
echo "    proposals→ $(get "$BASE/api/projects/$PID/proposals" "${AUTH_A[@]}" | grep -o '"status":"[a-z]*"' | tr '\n' ' ')"
check 409 "new proposal on in_progress project" -X POST "$BASE/api/projects/$PID/proposals" "${AUTH_1[@]}" "${JSON[@]}" -d "$PROP"

echo
printf '\033[1m%d passed, %d failed\033[0m\n' "$pass" "$fail"
[[ $fail -eq 0 ]]
