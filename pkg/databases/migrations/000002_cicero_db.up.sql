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
    ('admmin', 'JAAA', 2, 'admin@gmail.com', '28139b84-3205-41c5-a437-9419bc1f7dfc', '1234567890', 'avatar_url'),
    ('customer', 'JUUU', 1, 'customer@gmail.com', 'e67079e5-8869-4b6c-ba97-6eb00c127a16', '9876543210', 'avatar_url');

COMMIT;
