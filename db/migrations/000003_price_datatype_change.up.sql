ALTER TABLE "users" ALTER COLUMN "balance" TYPE real;
ALTER TABLE "stocks" ALTER COLUMN "price" TYPE real;
ALTER TABLE "sell_orders" ALTER COLUMN "price" TYPE real;
ALTER TABLE "buy_orders" ALTER COLUMN "price" TYPE real;
ALTER TABLE "trades" ALTER COLUMN "price" TYPE real;
ALTER TABLE "portfolio" ALTER COLUMN "purchase_price" TYPE real;
ALTER TABLE "portfolio" ADD COLUMN "quantity" integer NOT NULL;
ALTER TABLE "price_history" ALTER COLUMN "price" TYPE real;