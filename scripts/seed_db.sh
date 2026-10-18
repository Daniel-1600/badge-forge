#!/usr/bin/env bash
# Wipe persons, badges and assertions in tahrir_go and reload them from the
# tahrir reference database (a Fedora Badges copy) in the same container.
#
# Needs: the tahrir-pg container running with both databases, and migrations
# applied to tahrir_go. The data is real people's, so keep it local.
#
# Usage:
#   ./scripts/seed_db.sh
set -euo pipefail

CONTAINER="${CONTAINER:-tahrir-pg}"
SRC_DB="${SRC_DB:-tahrir}"
DST_DB="${DST_DB:-tahrir_go}"
PG_USER="${PG_USER:-tahrir}"

psql_src() { podman exec -i "$CONTAINER" psql -U "$PG_USER" -d "$SRC_DB" -v ON_ERROR_STOP=1 -q "$@"; }
psql_dst() { podman exec -i "$CONTAINER" psql -U "$PG_USER" -d "$DST_DB" -v ON_ERROR_STOP=1 -q "$@"; }

# Dump the three tables as one COPY stream, then load it in a single
# transaction so a failure leaves tahrir_go untouched.
#  - badges.issuer_id is NULL: in tahrir it points at the issuers table, but
#    tahrir_go references persons, and no person matches.
#  - test-badge and its assertions are leftovers from manual testing.
{
  echo 'BEGIN;'
  echo 'TRUNCATE assertions, badges, persons;'

  echo 'COPY persons (id, email, nickname, website, bio, created_on, optout, rank, last_login, avatar) FROM STDIN;'
  psql_src -c "\copy (select id, email, nickname, website, bio, created_on, opt_out, rank, last_login, _avatar from persons) to stdout"
  echo '\.'

  echo 'COPY badges (id, name, image, description, criteria, issuer_id, created_on, tags, stl) FROM STDIN;'
  psql_src -c "\copy (select id, name, image, description, criteria, null, created_on, tags, stl from badges where id <> 'test-badge') to stdout"
  echo '\.'

  echo 'COPY assertions (id, badge_id, person_id, salt, issued_on, recipient, issued_for) FROM STDIN;'
  psql_src -c "\copy (select id, badge_id, person_id, salt, issued_on, recipient, issued_for from assertions where badge_id <> 'test-badge') to stdout"
  echo '\.'

  # New persons from POST /persons must not collide with the copied IDs.
  echo "SELECT setval('persons_id_seq', (SELECT max(id) FROM persons));"
  echo 'COMMIT;'
} | psql_dst

psql_dst -c "select (select count(*) from persons) as persons, (select count(*) from badges) as badges, (select count(*) from assertions) as assertions;"
