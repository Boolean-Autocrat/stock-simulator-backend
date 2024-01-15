CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "users" (
  "id" uuid DEFAULT uuid_generate_v4(),
  "full_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "picture" varchar NOT NULL,
  "balance" decimal DEFAULT 10000.00,
  PRIMARY KEY ("id")
);

CREATE TABLE "access_tokens" (
  "id" uuid DEFAULT uuid_generate_v4(),
  "user_id" uuid NOT NULL,
  "token" varchar NOT NULL,
  "expires_at" timestamp NOT NULL,
  PRIMARY KEY ("id"),
  FOREIGN KEY ("user_id") REFERENCES "users" ("id")
);

CREATE TABLE "refresh_tokens" (
  "id" uuid DEFAULT uuid_generate_v4(),
  "user_id" uuid NOT NULL,
  "token" varchar NOT NULL,
  "expires_at" timestamp NOT NULL,
  PRIMARY KEY ("id"),
  FOREIGN KEY ("user_id") REFERENCES "users" ("id")
);


CREATE TABLE "stocks" (
  "id" uuid DEFAULT uuid_generate_v4(),
  "name" varchar NOT NULL,
  "symbol" varchar NOT NULL,
  "price" decimal NOT NULL,
  "is_crypto" bool NOT NULL,
  "is_stock" bool NOT NULL,
  "quantity" int NOT NULL,
  PRIMARY KEY ("id")
);

CREATE TABLE "portfolio" (
  "id" uuid DEFAULT uuid_generate_v4(),
  "user_id" uuid NOT NULL,
  "stock_id" uuid NOT NULL,
  "purchase_price" decimal NOT NULL,
  "purchased_at" timestamp NOT NULL DEFAULT (now()),
  PRIMARY KEY ("id")
);

CREATE TABLE "price_history" (
  "id" uuid DEFAULT uuid_generate_v4(),
  "stock_id" uuid NOT NULL,
  "price" decimal NOT NULL,
  "price_at" timestamp NOT NULL DEFAULT (now()),
  PRIMARY KEY ("id")
);

CREATE TABLE "news" (
  "id" uuid DEFAULT uuid_generate_v4(),
  "title" varchar NOT NULL,
  "description" varchar NOT NULL,
  "photo" varchar,
  "likes" int DEFAULT 0 NOT NULL,
  "dislikes" int DEFAULT 0 NOT NULL,
  PRIMARY KEY ("id")
);

CREATE TABLE "news_sentiment" (
  "id" uuid DEFAULT uuid_generate_v4(),
  "article_id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  "like" bool NOT NULL,
  "dislike" bool NOT NULL,
  PRIMARY KEY ("id")
);

ALTER TABLE "portfolio" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "portfolio" ADD FOREIGN KEY ("stock_id") REFERENCES "stocks" ("id");

ALTER TABLE "price_history" ADD FOREIGN KEY ("stock_id") REFERENCES "stocks" ("id");

ALTER TABLE "news_sentiment" ADD FOREIGN KEY ("article_id") REFERENCES "news" ("id");

ALTER TABLE "news_sentiment" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
