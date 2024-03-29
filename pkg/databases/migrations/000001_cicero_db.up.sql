BEGIN;

--Set timezone
SET TIME ZONE 'Asia/Bangkok';

--Install uuid extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

--users_id -> U000001
--Create sequence
CREATE SEQUENCE users_id_seq START WITH 1 INCREMENT BY 1;
CREATE SEQUENCE products_id_seq START WITH 1 INCREMENT BY 1;
CREATE SEQUENCE orders_id_seq START WITH 1 INCREMENT BY 1;

--Auto update
CREATE OR REPLACE FUNCTION set_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;   
END;
$$ language 'plpgsql';

CREATE TABLE "User" (
  "id" VARCHAR(7) PRIMARY KEY DEFAULT CONCAT('U', LPAD(NEXTVAL('users_id_seq')::TEXT, 6, '0')),
  "fname" VARCHAR NOT NULL,
  "lname" VARCHAR NOT NULL,
  "role_id" INT NOT NULL,
  "email" VARCHAR UNIQUE NOT NULL,
  "password" VARCHAR NOT NULL,
  "phone" VARCHAR NOT NULL UNIQUE,
  "avatar" VARCHAR DEFAULT 'https://storage.googleapis.com/cicero-bucket-prod/user.jpg',
  "dob" VARCHAR NOT NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT now(),
  "updated_at" TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE "Role" (
  "id" SERIAL PRIMARY KEY,
  "title" VARCHAR NOT NULL UNIQUE
);

CREATE TABLE "Oauth" (
  "id" uuid NOT NULL UNIQUE PRIMARY KEY DEFAULT uuid_generate_v4(),
  "user_id" VARCHAR NOT NULL,
  "access_token" VARCHAR NOT NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT now(),
  "updated_at" TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE "Wishlist" (
  "id" uuid NOT NULL UNIQUE PRIMARY KEY DEFAULT uuid_generate_v4(),
  "user_id" VARCHAR NOT NULL,
  "product_id" VARCHAR NOT NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE "Product" (
  "id" VARCHAR(7) PRIMARY KEY DEFAULT CONCAT('P', LPAD(NEXTVAL('products_id_seq')::TEXT, 6, '0')),
  "product_title" VARCHAR NOT NULL,
  "product_price" FLOAT NOT NULL DEFAULT 0,
  "product_color" VARCHAR NOT NULL,
  "product_sex" VARCHAR NOT NULL,
  "product_desc" VARCHAR NOT NULL,
  "product_size" VARCHAR NOT NULL,
  "product_category" VARCHAR NOT NULL,
  "product_stock" INT NOT NULL DEFAULT 0,
  "created_at" TIMESTAMP NOT NULL DEFAULT now(),
  "updated_at" TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE "Image" (
  "id" uuid NOT NULL UNIQUE PRIMARY KEY DEFAULT uuid_generate_v4(),
  "product_id" VARCHAR NOT NULL,
  "url" VARCHAR NOT NULL,
  "filename" VARCHAR NOT NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE "Cart" (
  "id" uuid NOT NULL UNIQUE PRIMARY KEY DEFAULT uuid_generate_v4(),
  "size" VARCHAR NOT NULL,
  "qty" INT NOT NULL DEFAULT 1,
  "user_id" VARCHAR NOT NULL,
  "product_id" VARCHAR NOT NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT now(),
  "updated_at" TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE "Order" (
  "id" VARCHAR(7) PRIMARY KEY DEFAULT CONCAT('O', LPAD(NEXTVAL('orders_id_seq')::TEXT, 6, '0')),
  "user_id" VARCHAR NOT NULL,
  "total" FLOAT NOT NULL DEFAULT 0,
  "status" VARCHAR NOT NULL,
  "products" jsonb NOT NULL,
  "address" jsonb NOT NULL,
  "payment_detail" jsonb NOT NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT now(),
  "updated_at" TIMESTAMP NOT NULL DEFAULT now()
);



ALTER TABLE "Oauth" ADD FOREIGN KEY ("user_id") REFERENCES "User" ("id")  ON DELETE CASCADE;
ALTER TABLE "User" ADD FOREIGN KEY ("role_id") REFERENCES "Role" ("id") ON DELETE CASCADE;
ALTER TABLE "Wishlist" ADD FOREIGN KEY ("user_id") REFERENCES "User" ("id") ON DELETE CASCADE;
ALTER TABLE "Wishlist" ADD FOREIGN KEY ("product_id") REFERENCES "Product" ("id") ON DELETE CASCADE;
ALTER TABLE "Image" ADD FOREIGN KEY ("product_id") REFERENCES "Product" ("id") ON DELETE CASCADE;
ALTER TABLE "Cart" ADD FOREIGN KEY ("user_id") REFERENCES "User" ("id") ON DELETE CASCADE;
ALTER TABLE "Cart" ADD FOREIGN KEY ("product_id") REFERENCES "Product" ("id") ON DELETE CASCADE;
ALTER TABLE "Order" ADD FOREIGN KEY ("user_id") REFERENCES "User" ("id") ON DELETE CASCADE;

CREATE TRIGGER set_updated_at_timestamp_users_table BEFORE UPDATE ON "User" FOR EACH ROW EXECUTE PROCEDURE set_updated_at_column();
CREATE TRIGGER set_updated_at_timestamp_oauth_table BEFORE UPDATE ON "Oauth" FOR EACH ROW EXECUTE PROCEDURE set_updated_at_column();
CREATE TRIGGER set_updated_at_timestamp_product_table BEFORE UPDATE ON "Product" FOR EACH ROW EXECUTE PROCEDURE set_updated_at_column();
CREATE TRIGGER set_updated_at_timestamp_image_table BEFORE UPDATE ON "Image" FOR EACH ROW EXECUTE PROCEDURE set_updated_at_column();
CREATE TRIGGER set_updated_at_timestamp_cart_table BEFORE UPDATE ON "Cart" FOR EACH ROW EXECUTE PROCEDURE set_updated_at_column();
CREATE TRIGGER set_updated_at_timestamp_order_table BEFORE UPDATE ON "Order" FOR EACH ROW EXECUTE PROCEDURE set_updated_at_column();


COMMIT;