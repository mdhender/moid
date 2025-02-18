--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

CREATE TABLE articles
(
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    title          TEXT     NOT NULL,
    slug           TEXT     NOT NULL UNIQUE,
    published      INTEGER  NOT NULL DEFAULT 0,
    date_published DATETIME,
    date_updated   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT articles_slug_format CHECK (slug GLOB '[a-z0-9-]*')
);

CREATE INDEX idx_articles_published ON articles (published, date_published);
CREATE INDEX idx_articles_slug ON articles (slug);
