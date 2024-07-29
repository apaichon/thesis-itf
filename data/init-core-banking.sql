
CREATE TABLE BankAccount (
    AccountID INTEGER PRIMARY KEY AUTOINCREMENT,
    UserID INTEGER REFERENCES UserAccount(UserID),
    AccountNumber TEXT NOT NULL UNIQUE,
    Balance REAL DEFAULT 0,
    AccountTypeID SMALLINT,
     CreatedAt DATETIME,
    CreatedBy INTEGER,
    UpdatedAt DATETIME,
    UpdatedBy INTEGER,
    Remarks TEXT,
    StatusID TINYINT,
    FOREIGN KEY (UserID) REFERENCES UserAccount(UserID)
)


CREATE TABLE UserProfile (
    UserID INTEGER PRIMARY KEY AUTOINCREMENT,
    FirstName TEXT,
    LastName TEXT,
    DateOfBirth DATE,
    Address TEXT,
    CreatedAt DATETIME,
    CreatedBy INTEGER,
    UpdatedAt DATETIME,
    UpdatedBy INTEGER,
    Remarks TEXT,
    StatusID TINYINT,
    FOREIGN KEY (UserID) REFERENCES UserAccount(UserID)
)

CREATE TABLE UserAccount (
    UserID INTEGER PRIMARY KEY AUTOINCREMENT,
    Username TEXT NOT NULL UNIQUE,
    PasswordHash TEXT NOT NULL,
    Salt TEXT NOT NULL,
    Email TEXT,
    PhoneNumber TEXT,
    CreatedAt DATETIME,
    CreatedBy INTEGER,
    UpdatedAt DATETIME,
    UpdatedBy INTEGER,
    Remarks TEXT,
    StatusID TINYINT
)

CREATE TABLE TransactionHistory (
    TransactionID INTEGER PRIMARY KEY AUTOINCREMENT,
    RefID INTEGER,
    TransactionType TEXT,
    AccountID INTEGER REFERENCES BankAccount(AccountID),
    Amount REAL,
    TransactionDate DATETIME DEFAULT CURRENT_TIMESTAMP,
    CreatedAt DATETIME,
    CreatedBy INTEGER,
    Remarks TEXT,
    StatusID TINYINT,
    FOREIGN KEY (AccountID) REFERENCES BankAccount(AccountID)
)


CREATE TABLE Withdrawal (
    WithdrawalID INTEGER PRIMARY KEY AUTOINCREMENT,
    AccountID INTEGER,
    Amount REAL,
    WithdrawalDate DATETIME DEFAULT CURRENT_TIMESTAMP,
    CreatedBy INTEGER,
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
)


CREATE TABLE Deposit (
    DepositID INTEGER PRIMARY KEY AUTOINCREMENT,
    AccountID INTEGER,
    Amount REAL,
    DepositDate DATETIME DEFAULT CURRENT_TIMESTAMP,
    CreatedBy INTEGER,
    CreatedAt  DATETIME DEFAULT CURRENT_TIMESTAMP
)


-- Insert sample data into UserAccount
INSERT INTO UserAccount (Username, PasswordHash, Salt, Email, PhoneNumber, CreatedAt, CreatedBy, UpdatedAt, UpdatedBy, Remarks, StatusID)
VALUES 
('johndoe', 'hashedpassword123', 'salt123', 'john.doe@example.com', '1234567890', '2024-07-29 10:00:00', 1, '2024-07-29 10:00:00', 1, 'Initial user', 1),
('janedoe', 'hashedpassword456', 'salt456', 'jane.doe@example.com', '0987654321', '2024-07-29 10:00:00', 1, '2024-07-29 10:00:00', 1, 'Initial user', 1);

-- Insert sample data into UserProfile
INSERT INTO UserProfile (UserID, FirstName, LastName, DateOfBirth, Address, CreatedAt, CreatedBy, UpdatedAt, UpdatedBy, Remarks, StatusID)
VALUES 
(1, 'John', 'Doe', '1990-01-01', '123 Main St', '2024-07-29 10:00:00', 1, '2024-07-29 10:00:00', 1, 'Profile for John Doe', 1),
(2, 'Jane', 'Doe', '1992-02-02', '456 Main St', '2024-07-29 10:00:00', 1, '2024-07-29 10:00:00', 1, 'Profile for Jane Doe', 1);

-- Insert sample data into BankAccount
INSERT INTO BankAccount (UserID, AccountNumber, Balance, AccountTypeID, CreatedAt, CreatedBy, UpdatedAt, UpdatedBy, Remarks, StatusID)
VALUES 
(1, 'ACC1234567890', 1000.00, 1, '2024-07-29 10:00:00', 1, '2024-07-29 10:00:00', 1, 'John Doe account', 1),
(2, 'ACC0987654321', 1500.00, 1, '2024-07-29 10:00:00', 1, '2024-07-29 10:00:00', 1, 'Jane Doe account', 1);

-- Insert sample data into TransactionHistory
INSERT INTO TransactionHistory (RefID, TransactionType, AccountID, Amount, TransactionDate, CreatedAt, CreatedBy, Remarks, StatusID)
VALUES 
(1, 'Deposit', 1, 500.00, '2024-07-29 10:05:00', '2024-07-29 10:05:00', 1, 'Initial deposit', 1),
(2, 'Withdrawal', 2, 200.00, '2024-07-29 10:10:00', '2024-07-29 10:10:00', 1, 'Initial withdrawal', 1);

-- Insert sample data into Withdrawal
INSERT INTO Withdrawal (AccountID, Amount, WithdrawalDate, CreatedBy, CreatedAt)
VALUES 
(1, 100.00, '2024-07-29 10:15:00', 1, '2024-07-29 10:15:00'),
(2, 150.00, '2024-07-29 10:20:00', 1, '2024-07-29 10:20:00');

-- Insert sample data into Deposit
INSERT INTO Deposit (AccountID, Amount, DepositDate, CreatedBy, CreatedAt)
VALUES 
(1, 200.00, '2024-07-29 10:25:00', 1, '2024-07-29 10:25:00'),
(2, 300.00, '2024-07-29 10:30:00', 1, '2024-07-29 10:30:00');

