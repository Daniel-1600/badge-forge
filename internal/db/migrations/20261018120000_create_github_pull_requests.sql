-- +goose Up
-- GitHub data imported from the pulls API. IDs are GitHub's own, so
-- re-importing a PR updates it instead of duplicating it.
CREATE TABLE github_repositories (
    id          BIGINT PRIMARY KEY,
    full_name   TEXT NOT NULL UNIQUE,
    html_url    TEXT NOT NULL
);

CREATE TABLE github_users (
    id          BIGINT PRIMARY KEY,
    login       TEXT NOT NULL UNIQUE,
    type        TEXT NOT NULL,
    avatar_url  TEXT,
    html_url    TEXT,
    person_id   INTEGER UNIQUE,
    CONSTRAINT fk_github_users_person FOREIGN KEY (person_id) REFERENCES persons(id)
);

CREATE TABLE github_pull_requests (
    id                  BIGINT PRIMARY KEY,
    repository_id       BIGINT NOT NULL,
    number              INTEGER NOT NULL,
    author_id           BIGINT NOT NULL,
    title               TEXT NOT NULL,
    body                TEXT,
    state               TEXT NOT NULL,
    draft               BOOLEAN NOT NULL DEFAULT FALSE,
    author_association  TEXT,
    base_ref            TEXT NOT NULL,
    head_ref            TEXT NOT NULL,
    head_sha            TEXT NOT NULL,
    merge_commit_sha    TEXT,
    milestone           TEXT,
    html_url            TEXT NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL,
    updated_at          TIMESTAMPTZ NOT NULL,
    closed_at           TIMESTAMPTZ,
    merged_at           TIMESTAMPTZ,
    CONSTRAINT uq_github_pull_requests_repo_number UNIQUE (repository_id, number),
    CONSTRAINT fk_github_pull_requests_repo FOREIGN KEY (repository_id) REFERENCES github_repositories(id),
    CONSTRAINT fk_github_pull_requests_author FOREIGN KEY (author_id) REFERENCES github_users(id)
);

CREATE INDEX idx_github_pull_requests_author_merged ON github_pull_requests (author_id, merged_at);

CREATE TABLE github_labels (
    id              BIGINT PRIMARY KEY,
    repository_id   BIGINT NOT NULL,
    name            TEXT NOT NULL,
    color           TEXT,
    description     TEXT,
    CONSTRAINT uq_github_labels_repo_name UNIQUE (repository_id, name),
    CONSTRAINT fk_github_labels_repo FOREIGN KEY (repository_id) REFERENCES github_repositories(id)
);

CREATE TABLE github_pull_request_labels (
    pull_request_id BIGINT NOT NULL,
    label_id        BIGINT NOT NULL,
    PRIMARY KEY (pull_request_id, label_id),
    CONSTRAINT fk_github_pr_labels_pr FOREIGN KEY (pull_request_id) REFERENCES github_pull_requests(id) ON DELETE CASCADE,
    CONSTRAINT fk_github_pr_labels_label FOREIGN KEY (label_id) REFERENCES github_labels(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE github_pull_request_labels;
DROP TABLE github_labels;
DROP TABLE github_pull_requests;
DROP TABLE github_users;
DROP TABLE github_repositories;
