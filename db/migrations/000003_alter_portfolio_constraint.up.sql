ALTER TABLE portfolio DROP CONSTRAINT portfolio_stock_key;
ALTER TABLE portfolio ADD CONSTRAINT portfolio_stock_user_key UNIQUE ("stock", "user");