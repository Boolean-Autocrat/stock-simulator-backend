CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS "users" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "full_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "picture" varchar NOT NULL,
  "balance" real NOT NULL
);

CREATE TABLE IF NOT EXISTS "watchlist" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "user" uuid NOT NULL,
  "stock" uuid NOT NULL,
  "added_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE IF NOT EXISTS "access_tokens" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "user" uuid NOT NULL,
  "token" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE IF NOT EXISTS "stocks" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "name" varchar NOT NULL,
  "symbol" varchar NOT NULL,
  "price" real NOT NULL,
  "ipo_quantity" int NOT NULL,
  "in_circulation" int NOT NULL DEFAULT 0,
  "is_stock" bool NOT NULL,
  "is_crypto" bool NOT NULL,
  "trend" varchar NOT NULL DEFAULT 'unchanged',
  "percentage_change" real NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS "price_history" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "stock" uuid NOT NULL,
  "price" real NOT NULL,
  "price_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE IF NOT EXISTS "portfolio" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "user" uuid NOT NULL,
  "stock" uuid UNIQUE NOT NULL,
  "quantity" int NOT NULL
);

CREATE TABLE IF NOT EXISTS "orders" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "user" uuid NOT NULL,
  "stock" uuid NOT NULL,
  "quantity" int NOT NULL,
  "fulfilled_quantity" int NOT NULL DEFAULT 0,
  "price" real NOT NULL,
  "is_buy" bool NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE IF NOT EXISTS "trade_history" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "stock" uuid NOT NULL,
  "quantity" int NOT NULL,
  "price" real NOT NULL,
  "buyer" uuid NOT NULL,
  "seller" uuid NOT NULL,
  "traded_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE IF NOT EXISTS "ipo_history" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "user" uuid NOT NULL,
  "stock" uuid NOT NULL,
  "quantity" int NOT NULL,
  "price" real NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE IF NOT EXISTS "news" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "title" varchar NOT NULL,
  "author" varchar NOT NULL,
  "content" text NOT NULL,
  "tag" varchar NOT NULL,
  "image" varchar NOT NULL DEFAULT 'https://placehold.co/500x500',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE IF NOT EXISTS "news_sentiment" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "article" uuid NOT NULL,
  "user" uuid NOT NULL,
  "like" bool NOT NULL,
  "dislike" bool NOT NULL
);

ALTER TABLE "watchlist" ADD FOREIGN KEY ("user") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "watchlist" ADD FOREIGN KEY ("stock") REFERENCES "stocks" ("id") ON DELETE CASCADE;

ALTER TABLE "access_tokens" ADD FOREIGN KEY ("user") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "price_history" ADD FOREIGN KEY ("stock") REFERENCES "stocks" ("id") ON DELETE CASCADE;

ALTER TABLE "portfolio" ADD FOREIGN KEY ("user") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "portfolio" ADD FOREIGN KEY ("stock") REFERENCES "stocks" ("id") ON DELETE CASCADE;

ALTER TABLE "orders" ADD FOREIGN KEY ("user") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "orders" ADD FOREIGN KEY ("stock") REFERENCES "stocks" ("id") ON DELETE CASCADE;

ALTER TABLE "trade_history" ADD FOREIGN KEY ("stock") REFERENCES "stocks" ("id") ON DELETE CASCADE;

ALTER TABLE "trade_history" ADD FOREIGN KEY ("buyer") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "trade_history" ADD FOREIGN KEY ("seller") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "ipo_history" ADD FOREIGN KEY ("user") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "ipo_history" ADD FOREIGN KEY ("stock") REFERENCES "stocks" ("id") ON DELETE CASCADE;

ALTER TABLE "news_sentiment" ADD FOREIGN KEY ("article") REFERENCES "news" ("id") ON DELETE CASCADE;

ALTER TABLE "news_sentiment" ADD FOREIGN KEY ("user") REFERENCES "users" ("id") ON DELETE CASCADE;