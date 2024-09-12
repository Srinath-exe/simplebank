ALTER TABLE IF EXISTS "accounts" DROP CONSTRAINT IF EXISTS "accounts_owner_fkey";

ALTER TABLE IF EXISTS "entries" DROP CONSTRAINT IF EXISTS "entries_account_id_fkey";
ALTER TABLE IF EXISTS "transfers" DROP CONSTRAINT IF EXISTS "transfers_from_account_id_fkey";
ALTER TABLE IF EXISTS "transfers" DROP CONSTRAINT IF EXISTS "transfers_to_account_id_fkey";


DROP TABLE IF EXISTS "transfers" CASCADE;
DROP TABLE IF EXISTS "entires" CASCADE;
DROP TABLE IF EXISTS "accounts" CASCADE;