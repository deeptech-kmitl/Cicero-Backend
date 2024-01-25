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

COMMIT;
