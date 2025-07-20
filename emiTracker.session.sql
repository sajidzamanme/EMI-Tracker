-- Table Creations
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
CREATE TABLE subscriptions (
  subID INT PRIMARY KEY AUTO_INCREMENT,
  ownerID INT NOT NULL,
  FOREIGN KEY (ownerID) REFERENCES users(userID),
  subName VARCHAR(200),
  totalAmount INT,
  paidAmount INT,
  paymentAmount INT,
  startDate DATE,
  endDate DATE,
  deductDay INT
);
-- Table Creations

-- Table Insertions
-- @block
INSERT INTO users(name, email, pass, totalLoaned, totalPaid) VALUES(
  "test",
  "test2@mail.com",
  "testpass",
  20000,
  15000
);

-- @block
INSERT INTO subscriptions(ownerID, subName, totalAmount, paidAmount, paymentAmount, startDate, endDate, deductDay) VALUES(
  9,
  "test sub 3",
  20000,
  5000,
  2000,
  '2025-07-20',
  '2026-12-20',
  5
);
-- Table Insertions

-- Table Updates
-- @block
UPDATE users
SET name = "up",
  email = "up@mail.com",
  pass = "upupup"
WHERE userID = 2;

-- @block
UPDATE subscriptions
SET paidAmount = 0 WHERE subID = 1;
-- Table Updates

-- Table Deletions
-- @block
DELETE FROM users WHERE userID = 1;

-- @block
DELETE FROM subscriptions WHERE subID = 2;
-- Table Deletions


-- Table Selections
-- @block
SELECT * FROM users;

-- @block
SELECT * FROM users WHERE userID = 3;

-- @block
SELECT * FROM subscriptions;
-- Table Selections


-- @block
SELECT * FROM users
INNER JOIN subscriptions
ON subscriptions.ownerID = users.userID;




-- @block
DROP TABLE users;

-- @block
DROP TABLE subscriptions;


-- @block
SELECT name, subID, subName, totalAmount, paidAmount, paymentAmount
FROM subscriptions INNER JOIN users ON userID = ownerID
WHERE ownerID = 9;

-- @block
SELECT * FROM subscriptions WHERE subID = 3

