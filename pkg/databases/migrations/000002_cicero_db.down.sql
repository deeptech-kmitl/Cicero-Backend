BEGIN;

TRUNCATE TABLE "User" CASCADE;
TRUNCATE TABLE "Role" CASCADE;
TRUNCATE TABLE "Oauth" CASCADE;
TRUNCATE TABLE "Product" CASCADE;
TRUNCATE TABLE "Wishlist" CASCADE;
TRUNCATE TABLE "Image" CASCADE;
TRUNCATE TABLE "Cart" CASCADE;
TRUNCATE TABLE "Order" CASCADE;

SELECT SETVAL ((SELECT PG_GET_SERIAL_SEQUENCE('"Role"', 'id')), 1, FALSE);

COMMIT;