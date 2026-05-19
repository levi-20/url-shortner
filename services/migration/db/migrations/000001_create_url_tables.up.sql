CREATE TABLE redirection (
    code       CHAR(7)    PRIMARY KEY,          -- indexed (PK)
    url        TEXT       NOT NULL,
    url_hash   CHAR(64)   NOT NULL UNIQUE,      -- indexed (UNIQUE)
    expire_at  TIMESTAMP  NULL
);

CREATE TABLE metadata (
    code        CHAR(7)    PRIMARY KEY REFERENCES redirection(code),  -- indexed (PK)
    created_at  TIMESTAMP  NOT NULL DEFAULT now(),
    clicks     BIGINT     NOT NULL DEFAULT 0
);