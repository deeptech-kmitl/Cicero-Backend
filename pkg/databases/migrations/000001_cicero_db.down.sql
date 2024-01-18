BEGIN;

DROP TRIGGER IF EXISTS set_updated_at_timestamp_users_table ON "User";

DROP FUNCTION IF EXISTS set_updated_at_column();

DROP TABLE IF EXISTS "User" CASCADE;
DROP TABLE IF EXISTS "Role" CASCADE;

DROP SEQUENCE IF EXISTS users_id_seq;


COMMIT;