ALTER TABLE accounts DROP FOREIGN KEY accounts_ibfk_1;
ALTER TABLE accounts DROP INDEX owner_currency_key;
DROP TABLE IF EXISTS users;