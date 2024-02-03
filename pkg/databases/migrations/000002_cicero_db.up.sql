BEGIN;

INSERT INTO "Role" (
    "title"
)
VALUES
    ('customer'),
    ('admin');

INSERT INTO "User" (
    "fname", 
    "lname", 
    "role_id", 
    "email", 
    "password", 
    "phone", 
    "avatar"
) 
VALUES
    ('admmin', 'JAAA', 2, 'admin@gmail.com', '$2a$10$qS6yyjwgkcsgUsa9o1ZiTeHcAXDbJC7QVpE8c/kpc8dHIRl6iBO/m', '1234567890', 'avatar_url'),
    ('customer', 'JUUU', 1, 'customer@gmail.com', '$2a$10$qS6yyjwgkcsgUsa9o1ZiTeHcAXDbJC7QVpE8c/kpc8dHIRl6iBO/m', '9876543210', 'avatar_url');


INSERT INTO "Product" ("product_title", "product_price", "product_color", "product_sex", "product_desc", "product_size", "product_category")
VALUES
  ('Running Shoes', '49.99', 'Blue', 'Male', 'Comfortable running shoes with advanced cushioning', 'XL', 'Footwear'),
  ('Smartphone', '799.99', 'Black', 'Unisex', 'High-performance smartphone with dual cameras', 'N/A', 'Electronics'),
  ('Backpack', '29.99', 'Red', 'Female', 'Durable backpack with multiple compartments', 'N/A', 'Accessories'),
  ('T-shirt', '19.99', 'Green', 'Male', 'Casual cotton t-shirt for everyday wear', 'M', 'Apparel');

INSERT INTO "Image" ("product_id", "url", "filename")
VALUES
  ('P000001', 'https://example.com/running-shoes.jpg', 'running-shoes.jpg'),
  ('P000002', 'https://example.com/smartphone.jpg', 'smartphone.jpg'),
  ('P000003', 'https://example.com/backpack.jpg', 'backpack.jpg'),
  ('P000004', 'https://example.com/t-shirt.jpg', 't-shirt.jpg');

INSERT INTO "Wishlist" ("user_id", "product_id")
VALUES
  ('U000002', 'P000001'),
  ('U000002', 'P000002');

INSERT INTO "Cart" ("size", "qty", "product_id", "user_id")
VALUES
  ('XL', 1, 'P000001', 'U000002'),
  ('M', 2, 'P000002', 'U000002');

COMMIT;
