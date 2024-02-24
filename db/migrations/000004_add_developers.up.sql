CREATE TABLE IF NOT EXISTS "developers" (
    "id" SERIAL PRIMARY KEY,
    "name" varchar NOT NULL,
    "title" varchar NOT NULL,
    "picture" varchar NOT NULL,
    "email" varchar NOT NULL,
    "github_link" varchar NOT NULL
);