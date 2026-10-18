#!/usr/bin/env bash
# Replay recently merged rust-lang/rust PRs into a running badge-forge so you can
# watch real contributors earn badges and the milestone rule fire.
#
# Needs: gh (logged in), jq, curl, openssl, and badge-forge running on $API with
# migrations applied. Results are for local testing only: the contributors have
# not opted in, so don't publish them.
#
# Usage:
#   ./scripts/rust_replay.sh
#   API=http://localhost:9090 PAGES=5 ./scripts/rust_replay.sh
set -euo pipefail

API="${API:-http://localhost:8080}"
PAGES="${PAGES:-3}"                       # 100 closed PRs per page
REPO="${REPO:-rust-lang/rust}"
BADGE_ID="rust-merged-pr"
# Must match badge_id of the milestone rule in config/rules.yaml.
MILESTONE_BADGE_ID="${MILESTONE_BADGE_ID:-badge-first-steps}"
OUT="$(mktemp)"
trap 'rm -f "$OUT"' EXIT

for cmd in gh jq curl openssl; do
  command -v "$cmd" >/dev/null || { echo "missing dependency: $cmd" >&2; exit 1; }
done

curl -sf "$API/healthz" >/dev/null || { echo "badge-forge is not reachable at $API" >&2; exit 1; }

# find_or_create_person NICKNAME -> prints the person's ID
find_or_create_person() {
  curl -sf -X POST "$API/persons" -H 'Content-Type: application/json' \
    -d "$(jq -n --arg n "$1" '{Nickname: $n}')" | jq -r '.ID'
}

# create_badge ID NAME DESCRIPTION CRITERIA (a 500 just means it already exists)
create_badge() {
  curl -s -o /dev/null -X POST "$API/badges" -H 'Content-Type: application/json' \
    -d "$(jq -n --arg id "$1" --arg name "$2" --arg desc "$3" --arg crit "$4" --argjson issuer "$ISSUER_ID" '{
      ID: $id, Name: $name, Description: $desc, Criteria: $crit, IssuerID: $issuer,
      Image: "https://www.rust-lang.org/static/images/rust-logo-blk.svg",
      Tags: "rust,contribution"}')" || true
}

# 1. The issuer is a person row (badges.issuer_id references persons.id).
ISSUER_ID="$(find_or_create_person "rust-lang")"

# 2. The per-PR badge, plus the milestone badge the rule engine awards.
create_badge "$BADGE_ID" "Merged PR to $REPO" \
  "Had a pull request merged into $REPO." "One merged PR in $REPO."
create_badge "$MILESTONE_BADGE_ID" "First Steps" \
  "Awarded automatically for your first 3 contributions." "Three contributions."

# 3. Fetch merged PRs, skipping bots and rollup PRs (rollups would credit the
#    person who assembled them instead of the original authors). Sort by merge
#    time so milestone badges are awarded in the right order.
for page in $(seq 1 "$PAGES"); do
  gh api "repos/$REPO/pulls?state=closed&sort=updated&direction=desc&per_page=100&page=$page" \
    --jq '.[]
          | select(.merged_at != null)
          | select(.user.type != "Bot" and (.user.login | endswith("[bot]") | not))
          | select(.title | startswith("Rollup of") | not)
          | [.merged_at, .user.login, (.number | tostring)] | @tsv'
done | sort > "$OUT"

echo "Replaying $(wc -l < "$OUT") merged PRs from $REPO into $API"

# 4. For each PR: find-or-create the author, then award the badge. The
#    assertion ID comes from the PR number, so re-runs don't double-award.
while IFS=$'\t' read -r merged_at login number; do
  person_id="$(find_or_create_person "$login")"
  curl -s -o /dev/null -w "%{http_code}  $login  $REPO#$number\n" \
    -X POST "$API/assertions" -H 'Content-Type: application/json' \
    -d "$(jq -n --arg id "rust-pr-$number" --arg badge "$BADGE_ID" --argjson person "$person_id" \
          --arg salt "$(openssl rand -hex 8)" --arg for "Merged $REPO#$number ($merged_at)" \
          '{ID: $id, BadgeID: $badge, PersonID: $person, Salt: $salt, IssuedFor: $for}')"
done < "$OUT"

echo
echo "Top contributors in this window (3+ merges earn \"$MILESTONE_BADGE_ID\"):"
cut -f2 "$OUT" | sort | uniq -c | sort -rn | head -5
echo
echo "Check one with: curl -s $API/persons/nickname/<login>/badges | jq '.[].Badge.Name'"
