CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS "users" (
  "id" uuid DEFAULT uuid_generate_v4(),
  "full_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "picture" varchar NOT NULL,
  "balance" decimal DEFAULT 10000.00 NOT NULL,
  PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "access_tokens" (
  "id" uuid DEFAULT uuid_generate_v4(),
  "user_id" uuid NOT NULL UNIQUE,
  "token" varchar NOT NULL,
  "expires_at" timestamptz NOT NULL,
  PRIMARY KEY ("id"),
  FOREIGN KEY ("user_id") REFERENCES "users" ("id")
);

CREATE TABLE IF NOT EXISTS "refresh_tokens" (
  "id" uuid DEFAULT uuid_generate_v4(),
  "user_id" uuid NOT NULL UNIQUE,
  "token" varchar NOT NULL,
  "expires_at" timestamptz NOT NULL,
  PRIMARY KEY ("id"),
  FOREIGN KEY ("user_id") REFERENCES "users" ("id")
);

CREATE TABLE IF NOT EXISTS "stocks" (
  "id" uuid DEFAULT uuid_generate_v4(),
  "name" varchar NOT NULL,
  "symbol" varchar NOT NULL,
  "price" decimal NOT NULL,
  "is_crypto" bool NOT NULL,
  "is_stock" bool NOT NULL,
  "quantity" int NOT NULL,
  PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "portfolio" (
  "id" uuid DEFAULT uuid_generate_v4(),
  "user_id" uuid NOT NULL,
  "stock_id" uuid NOT NULL,
  "purchase_price" decimal NOT NULL,
  "purchased_at" timestamptz NOT NULL DEFAULT (now()),
  PRIMARY KEY ("id"),
  FOREIGN KEY ("user_id") REFERENCES "users" ("id"),
  FOREIGN KEY ("stock_id") REFERENCES "stocks" ("id")
);

CREATE TABLE IF NOT EXISTS "price_history" (
  "id" uuid DEFAULT uuid_generate_v4(),
  "stock_id" uuid NOT NULL,
  "price" decimal NOT NULL,
  "price_at" timestamptz NOT NULL DEFAULT (now()),
  PRIMARY KEY ("id"),
  FOREIGN KEY ("stock_id") REFERENCES "stocks" ("id")
);

CREATE TABLE IF NOT EXISTS "news" (
  "id" uuid DEFAULT uuid_generate_v4(),
  "title" varchar NOT NULL,
  "author" varchar NOT NULL,
  "content" TEXT NOT NULL,
  "tag" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "news_sentiment" (
  "id" uuid DEFAULT uuid_generate_v4(),
  "article_id" uuid NOT NULL,
  "user_id" uuid NOT NULL UNIQUE,
  "like" bool NOT NULL,
  "dislike" bool NOT NULL,
  PRIMARY KEY ("id"),
  FOREIGN KEY ("article_id") REFERENCES "news" ("id"),
  FOREIGN KEY ("user_id") REFERENCES "users" ("id")
);
