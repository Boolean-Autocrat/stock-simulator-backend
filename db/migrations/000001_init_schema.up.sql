CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "users" (
  "id" uuid DEFAULT uuid_generate_v4(),
  "full_name" varchar,
  "email" varchar UNIQUE,
  "picture" varchar,
  "balance" decimal DEFAULT 10000.00,
  PRIMARY KEY ("id")
);

CREATE TABLE "stocks" (
  "id" uuid DEFAULT uuid_generate_v4(),
  "symbol" varchar,
  "price" decimal,
  "is_crypto" bool,
  "is_stock" bool,
  PRIMARY KEY ("id")
);

CREATE TABLE "portfolio" (
  "id" uuid DEFAULT uuid_generate_v4(),
  "user_id" uuid,
  "stock_id" uuid,
  "purchase_price" decimal,
  "purchased_at" timestamp NOT NULL DEFAULT (now()),
  PRIMARY KEY ("id")
);

CREATE TABLE "price_history" (
  "id" uuid DEFAULT uuid_generate_v4(),
  "stock_id" uuid,
  "price" decimal,
  "price_at" timestamp NOT NULL DEFAULT (now()),
  PRIMARY KEY ("id")
);

CREATE TABLE "news" (
  "id" uuid DEFAULT uuid_generate_v4(),
  "title" varchar,
  "description" varchar,
  "photo" varchar,
  "likes" int,
  "dislikes" int,
  PRIMARY KEY ("id")
);

CREATE TABLE "news_sentiment" (
  "id" uuid DEFAULT uuid_generate_v4(),
  "article_id" uuid,
  "user_id" uuid,
  "like" bool,
  "dislike" bool,
  PRIMARY KEY ("id")
);

ALTER TABLE "portfolio" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "portfolio" ADD FOREIGN KEY ("stock_id") REFERENCES "stocks" ("id");

ALTER TABLE "price_history" ADD FOREIGN KEY ("stock_id") REFERENCES "stocks" ("id");

ALTER TABLE "news_sentiment" ADD FOREIGN KEY ("article_id") REFERENCES "news" ("id");

ALTER TABLE "news_sentiment" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
