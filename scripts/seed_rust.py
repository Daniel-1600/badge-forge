#!/usr/bin/env python3
"""Wipe tahrir_go and seed it from the Rust project.

- persons + github_users: one per person in rust-lang/team (people/*.toml)
- badges: one per team, working group and project group (teams/*.toml)
- assertions: one per lead, member and alumnus of each of those teams
- github_*: recently closed PRs of rust-lang/rust, from the GitHub API

Needs: git, gh (logged in), the tahrir-pg container running, and migrations
applied. Real emails from rust-lang/team are not copied: each person gets their
GitHub noreply address instead. Keep the data local.

Usage:
    python3 -I scripts/seed_rust.py
    PAGES=10 python3 -I scripts/seed_rust.py
"""
import json
import os
import pathlib
import secrets
import subprocess
import sys
import tempfile
import tomllib

TEAM_REPO = os.environ.get("TEAM_REPO", "https://github.com/rust-lang/team")
PR_REPO = os.environ.get("PR_REPO", "rust-lang/rust")
PAGES = int(os.environ.get("PAGES", "3"))  # 100 closed PRs per page
CONTAINER = os.environ.get("CONTAINER", "tahrir-pg")
DB = os.environ.get("DB_NAME", "tahrir_go").split()[0]
PG_USER = os.environ.get("DB_USER", "tahrir")

BADGE_IMAGE = "https://www.rust-lang.org/static/images/rust-logo-blk.svg"
# marker-teams (e.g. "all", "leads") only group other teams' people.
BADGE_KINDS = {"team", "working-group", "project-group"}
ROLE_TEXT = {"lead": "Lead of", "member": "Member of", "alumni": "Former member of"}


def noreply(github_id, login):
    return f"{github_id}+{login}@users.noreply.github.com"


def load_team(root):
    people = {}
    for f in sorted((root / "people").glob("*.toml")):
        d = tomllib.loads(f.read_text())
        people[d["github"].lower()] = d

    badges, memberships = [], []
    for f in sorted((root / "teams").glob("*.toml")):
        d = tomllib.loads(f.read_text())
        kind = d.get("kind", "team")
        if kind not in BADGE_KINDS:
            continue
        name, site = d["name"], d.get("website", {})
        badge_id = f"rust-team-{name}"
        label = site.get("name") or f"Rust {name} {kind.replace('-', ' ')}"
        badges.append({
            "id": badge_id,
            "name": label,
            "description": site.get("description") or f"Part of the Rust {name} {kind.replace('-', ' ')}.",
            "criteria": f"Listed as a lead, member or alumnus in rust-lang/team teams/{name}.toml.",
            "tags": f"rust,{kind}",
        })
        # A person keeps their highest role: lead > member > alumni.
        roles = {}
        for role, key in (("alumni", "alumni"), ("member", "members"), ("lead", "leads")):
            for m in d.get("people", {}).get(key, []):
                login = (m["github"] if isinstance(m, dict) else m).lower()
                if login in people:
                    roles[login] = role
        for login, role in roles.items():
            p = people[login]
            memberships.append({
                "id": f"{badge_id}-{p['github']}",
                "badge_id": badge_id,
                "github_id": p["github-id"],
                "salt": secrets.token_hex(8),
                "recipient": noreply(p["github-id"], p["github"]),
                "issued_for": f"{ROLE_TEXT[role]} {label} (rust-lang/team)",
            })

    persons = [{
        "github_id": p["github-id"],
        "login": p["github"],
        "email": noreply(p["github-id"], p["github"]),
    } for p in people.values()]
    return persons, badges, memberships


def fetch_prs():
    prs = []
    for page in range(1, PAGES + 1):
        out = subprocess.run(
            ["gh", "api", f"repos/{PR_REPO}/pulls?state=closed&sort=updated&direction=desc&per_page=100&page={page}"],
            check=True, capture_output=True, text=True).stdout
        batch = json.loads(out)
        prs.extend(batch)
        if len(batch) < 100:
            break
    return prs


# Everything is passed to Postgres as one JSON document and unpacked there, so
# no value is ever spliced into SQL text.
SQL = r"""
BEGIN;
CREATE TEMP TABLE seed (d jsonb);
COPY seed FROM STDIN;
%s
\.

TRUNCATE github_pull_request_labels, github_labels, github_pull_requests,
         github_users, github_repositories, assertions, badges, persons
         RESTART IDENTITY;

INSERT INTO persons (email, nickname, website, avatar, created_on, optout)
SELECT p->>'email', p->>'login', 'https://github.com/' || (p->>'login'),
       'https://avatars.githubusercontent.com/u/' || (p->>'github_id'), now(), false
FROM seed, jsonb_array_elements(d->'persons') p;

INSERT INTO github_users (id, login, type, avatar_url, html_url, person_id)
SELECT (p->>'github_id')::bigint, p->>'login', 'User',
       'https://avatars.githubusercontent.com/u/' || (p->>'github_id'),
       'https://github.com/' || (p->>'login'), pe.id
FROM seed, jsonb_array_elements(d->'persons') p
JOIN persons pe ON pe.nickname = p->>'login';

INSERT INTO badges (id, name, image, description, criteria, created_on, tags)
SELECT b->>'id', b->>'name', %s, b->>'description', b->>'criteria', now(), b->>'tags'
FROM seed, jsonb_array_elements(d->'badges') b;

INSERT INTO assertions (id, badge_id, person_id, salt, issued_on, recipient, issued_for)
SELECT m->>'id', m->>'badge_id', gu.person_id, m->>'salt', now(), m->>'recipient', m->>'issued_for'
FROM seed, jsonb_array_elements(d->'memberships') m
JOIN github_users gu ON gu.id = (m->>'github_id')::bigint;

-- Pull requests. Authors already loaded from rust-lang/team keep their person link.
INSERT INTO github_repositories (id, full_name, html_url)
SELECT DISTINCT (pr->'base'->'repo'->>'id')::bigint, pr->'base'->'repo'->>'full_name', pr->'base'->'repo'->>'html_url'
FROM seed, jsonb_array_elements(d->'prs') pr
ON CONFLICT (id) DO NOTHING;

INSERT INTO github_users (id, login, type, avatar_url, html_url)
SELECT DISTINCT ON ((pr->'user'->>'id')::bigint)
       (pr->'user'->>'id')::bigint, pr->'user'->>'login', pr->'user'->>'type',
       pr->'user'->>'avatar_url', pr->'user'->>'html_url'
FROM seed, jsonb_array_elements(d->'prs') pr
ON CONFLICT (id) DO NOTHING;

INSERT INTO github_pull_requests
SELECT (pr->>'id')::bigint, (pr->'base'->'repo'->>'id')::bigint, (pr->>'number')::int,
       (pr->'user'->>'id')::bigint, pr->>'title', pr->>'body', pr->>'state',
       COALESCE((pr->>'draft')::bool, false), pr->>'author_association',
       pr->'base'->>'ref', pr->'head'->>'ref', pr->'head'->>'sha', pr->>'merge_commit_sha',
       pr->'milestone'->>'title', pr->>'html_url',
       (pr->>'created_at')::timestamptz, (pr->>'updated_at')::timestamptz,
       (pr->>'closed_at')::timestamptz, (pr->>'merged_at')::timestamptz
FROM seed, jsonb_array_elements(d->'prs') pr;

INSERT INTO github_labels (id, repository_id, name, color, description)
SELECT DISTINCT ON ((l->>'id')::bigint)
       (l->>'id')::bigint, (pr->'base'->'repo'->>'id')::bigint, l->>'name', l->>'color', l->>'description'
FROM seed, jsonb_array_elements(d->'prs') pr, jsonb_array_elements(pr->'labels') l;

INSERT INTO github_pull_request_labels (pull_request_id, label_id)
SELECT (pr->>'id')::bigint, (l->>'id')::bigint
FROM seed, jsonb_array_elements(d->'prs') pr, jsonb_array_elements(pr->'labels') l;

COMMIT;

SELECT (SELECT count(*) FROM persons) AS persons,
       (SELECT count(*) FROM badges) AS badges,
       (SELECT count(*) FROM assertions) AS assertions,
       (SELECT count(*) FROM github_pull_requests) AS pull_requests,
       (SELECT count(*) FROM github_pull_requests WHERE merged_at IS NOT NULL) AS merged,
       (SELECT count(*) FROM github_users WHERE person_id IS NOT NULL) AS linked_users,
       (SELECT count(*) FROM github_users) AS github_users;
"""


def copy_text(s):
    # COPY text format: escape backslashes; json.dumps leaves no raw tabs/newlines.
    return s.replace("\\", "\\\\")


def main():
    with tempfile.TemporaryDirectory() as tmp:
        root = pathlib.Path(tmp) / "team"
        print(f"cloning {TEAM_REPO}", file=sys.stderr)
        subprocess.run(["git", "clone", "-q", "--depth", "1", TEAM_REPO, str(root)], check=True)
        persons, badges, memberships = load_team(root)

    print(f"fetching up to {PAGES * 100} closed PRs from {PR_REPO}", file=sys.stderr)
    prs = fetch_prs()

    doc = json.dumps({"persons": persons, "badges": badges, "memberships": memberships, "prs": prs})
    sql = SQL % (copy_text(doc), "'" + BADGE_IMAGE + "'")
    subprocess.run(
        ["podman", "exec", "-i", CONTAINER, "psql", "-U", PG_USER, "-d", DB, "-v", "ON_ERROR_STOP=1", "-q"],
        input=sql, text=True, check=True)


if __name__ == "__main__":
    main()
