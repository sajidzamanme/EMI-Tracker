
-- @block
CREATE TABLE users (
  userID INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(100),
  email VARCHAR(100) UNIQUE,
  pass VARCHAR(30),
  totalLoaned INT,
  totalPaid INT
);

-- @block
INSERT INTO users(name, email, pass, totalLoaned, totalPaid) VALUES(
  "test",
  "test2@mail.com",
  "testpass",
  20000,
  15000
);

-- @block
UPDATE users
SET name = "up",
  email = "up@mail.com",
  pass = "upupup"
WHERE userID = 2;

-- @block
DELETE FROM users WHERE userID = 1;

-- @block
SELECT * FROM users;

-- @block
SELECT * FROM users WHERE userID = 3

-- @block
DROP TABLE users;