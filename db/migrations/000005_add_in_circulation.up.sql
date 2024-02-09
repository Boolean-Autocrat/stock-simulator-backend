ALTER TABLE "stocks" ADD COLUMN "in_circulation" int NOT NULL DEFAULT 0;

CREATE TABLE IF NOT EXISTS "portfolio_purchase_history" (
    "portfolio_id" uuid,
    "purchase_price" decimal NOT NULL,
    "purchased_at" timestamptz NOT NULL DEFAULT (now()),
    PRIMARY KEY ("portfolio_id"),
    FOREIGN KEY ("portfolio_id") REFERENCES "portfolio" ("id")
);

ALTER TABLE "portfolio" DROP COLUMN "purchase_price";
ALTER TABLE "portfolio" DROP COLUMN "purchased_at";