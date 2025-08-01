-- Table Creations
-- @block
CREATE TABLE users (
  userID INT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(100) NOT NULL,
  email VARCHAR(100) UNIQUE NOT NULL,
  password TEXT NOT NULL,
  totalLoaned INT NOT NULL,
  totalPaid INT NOT NULL,
  activeEMI INT NOT NULL,
  completedEMI INT NOT NULL
);

-- @block
CREATE TABLE emiRecords (
  recordID INT PRIMARY KEY AUTO_INCREMENT,
  ownerID INT NOT NULL,
  title VARCHAR(200) NOT NULL,
  totalAmount INT NOT NULL,
  paidAmount INT NOT NULL,
  installmentAmount INT NOT NULL,
  startDate DATE NOT NULL,
  endDate DATE NOT NULL,
  deductDay INT NOT NULL,
  FOREIGN KEY (ownerID) REFERENCES users(userID) ON DELETE CASCADE
);

-- Drop Tables
-- @block
DROP TABLE users;

-- @block
DROP TABLE emiRecords;