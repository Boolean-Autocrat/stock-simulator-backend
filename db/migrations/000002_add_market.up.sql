CREATE TABLE IF NOT EXISTS "sell_orders" (
    "id" uuid,
    "user" uuid NOT NULL,
    "stock" uuid NOT NULL,
    "price" int NOT NULL,
    "quantity" int NOT NULL,
    "fulfilled" int NOT NULL DEFAULT 0,
    PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "buy_orders" (
    "id" uuid,
    "user" uuid NOT NULL,
    "stock" uuid NOT NULL,
    "price" int NOT NULL,
    "quantity" int NOT NULL,
    "fulfilled" int NOT NULL DEFAULT 0,
    PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "trades" (
    "id" uuid,
    "buy_order" uuid NOT NULL,
    "sell_order" uuid NOT NULL,
    "price" int NOT NULL,
    "quantity" int NOT NULL,
    PRIMARY KEY ("id")
);