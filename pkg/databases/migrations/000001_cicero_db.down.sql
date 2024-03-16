BEGIN;

DROP TRIGGER IF EXISTS set_updated_at_timestamp_users_table ON "User";
DROP TRIGGER IF EXISTS set_updated_at_timestamp_oauth_table ON "Oauth";
DROP TRIGGER IF EXISTS set_updated_at_timestamp_product_table ON "Product";
DROP TRIGGER IF EXISTS set_updated_at_timestamp_image_table ON "Image";
DROP TRIGGER IF EXISTS set_updated_at_timestamp_cart_table ON "Cart";
DROP TRIGGER IF EXISTS set_updated_at_timestamp_order_table ON "Order";

DROP FUNCTION IF EXISTS set_updated_at_column();

DROP TABLE IF EXISTS "User" CASCADE;
DROP TABLE IF EXISTS "Role" CASCADE;
DROP TABLE IF EXISTS "Oauth" CASCADE;
DROP TABLE IF EXISTS "Wishlist" CASCADE;
DROP TABLE IF EXISTS "Product" CASCADE;
DROP TABLE IF EXISTS "Image" CASCADE;
DROP TABLE IF EXISTS "Cart" CASCADE;
DROP TABLE IF EXISTS "Order" CASCADE;

DROP SEQUENCE IF EXISTS users_id_seq;
DROP SEQUENCE IF EXISTS products_id_seq;
DROP SEQUENCE IF EXISTS orders_id_seq;


COMMIT;