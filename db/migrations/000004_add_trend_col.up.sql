ALTER TABLE "stocks" ADD COLUMN "trend" varchar NOT NULL DEFAULT 'static';
ALTER TABLE "stocks" ADD COLUMN "percent_change" real NOT NULL DEFAULT 0.0;